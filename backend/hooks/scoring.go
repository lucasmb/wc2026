package hooks

import (
	"database/sql"
	"fmt"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

func RegisterScoringHooks(app core.App) {
	// Trigger scoring calculations after a match successfully completes
	app.OnRecordAfterUpdateSuccess("matches_id").BindFunc(func(e *core.RecordEvent) error {
		status := e.Record.GetString("status")

		// Only run scoring if the match status was updated to "finished"
		if status != "finished" {
			return e.Next()
		}

		scoreHome := e.Record.GetInt("score_home")
		scoreAway := e.Record.GetInt("score_away")
		phase := e.Record.GetString("phase")

		// Retrieve all predictions for this match using the standard FindRecordsByFilter pattern
		predictions, err := e.App.FindRecordsByFilter(
			"predictions_id",
			"match = {:matchID}",
			"",
			0,
			0,
			dbx.Params{"matchID": e.Record.Id},
		)
		if err != nil {
			return err
		}

		// Process points calculation for each submitted prediction in a transaction
		err = e.App.RunInTransaction(func(txApp core.App) error {
			affectedGroupIDs := make(map[string]bool)

			for _, pred := range predictions {
				predHome := pred.GetInt("predicted_home")
				predAway := pred.GetInt("predicted_away")
				groupID := pred.GetString("prediction_group")

				points := calculateMatchPoints(phase, scoreHome, scoreAway, predHome, predAway)

				pred.Set("points_awarded", points)
				pred.Set("is_locked", true)
				if err := txApp.Save(pred); err != nil {
					return err
				}

				if groupID != "" {
					affectedGroupIDs[groupID] = true
				}
			}

			// Recalculate total points and ranks for members of affected prediction groups
			for groupID := range affectedGroupIDs {
				if err := recalculateGroupLeaderboard(txApp, groupID); err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("failed to process match scoring: %w", err)
		}

		return e.Next()
	})
}

func calculateMatchPoints(phase string, actualHome, actualAway, predHome, predAway int) int {
	actualDiff := actualHome - actualAway
	predDiff := predHome - predAway

	correctResult := false
	if (actualDiff > 0 && predDiff > 0) || (actualDiff < 0 && predDiff < 0) || (actualDiff == 0 && predDiff == 0) {
		correctResult = true
	}

	exactScore := (actualHome == predHome && actualAway == predAway)

	var basePoints, exactBonus int
	switch phase {
	case "group":
		basePoints, exactBonus = 2, 1
	case "r32", "r16":
		basePoints, exactBonus = 3, 2
	case "qf", "sf":
		basePoints, exactBonus = 5, 3
	case "final", "third":
		basePoints, exactBonus = 8, 5
	default:
		basePoints, exactBonus = 2, 1
	}

	totalPoints := 0
	if correctResult {
		totalPoints += basePoints
	}
	if exactScore {
		totalPoints += exactBonus
	}

	return totalPoints
}

func recalculateGroupLeaderboard(txApp core.App, groupID string) error {
	// Retrieve group members using the standard Filter pattern
	members, err := txApp.FindRecordsByFilter(
		"group_members_id",
		"prediction_group = {:groupID}",
		"",
		0,
		0,
		dbx.Params{"groupID": groupID},
	)
	if err != nil {
		return err
	}

	type MemberScore struct {
		Record *core.Record
		Points int
	}
	var scores []MemberScore

	for _, member := range members {
		userID := member.GetString("user")

		// Run direct aggregate projection safely
		var total sql.NullInt64
		err := txApp.DB().
			Select("SUM(points_awarded)").
			From("predictions").
			Where(dbx.HashExp{"user": userID, "prediction_group": groupID}).
			Row(&total)

		points := 0
		if err == nil && total.Valid {
			points = int(total.Int64)
		}

		scores = append(scores, MemberScore{
			Record: member,
			Points: points,
		})
	}

	// Sort high to low
	for i := 0; i < len(scores); i++ {
		for j := i + 1; j < len(scores); j++ {
			if scores[i].Points < scores[j].Points {
				scores[i], scores[j] = scores[j], scores[i]
			}
		}
	}

	currentRank := 1
	for idx, s := range scores {
		if idx > 0 && s.Points < scores[idx-1].Points {
			currentRank++
		}
		s.Record.Set("total_points", s.Points)
		s.Record.Set("rank", currentRank)
		if err := txApp.Save(s.Record); err != nil {
			return err
		}
	}

	return nil
}
