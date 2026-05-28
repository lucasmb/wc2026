package hooks

import (
	"fmt"
	"sort"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

func RegisterTournamentHooks(app core.App) {
	// Recalculate standings when group stage matches conclude
	app.OnRecordAfterUpdateSuccess("matches_id").BindFunc(func(e *core.RecordEvent) error {
		phase := e.Record.GetString("phase")
		status := e.Record.GetString("status")

		if phase != "group" || status != "finished" {
			return e.Next()
		}

		groupCode := e.Record.GetString("group_code")
		if groupCode == "" {
			return e.Next()
		}

		if err := updateStandingsForGroup(e.App, groupCode); err != nil {
			return fmt.Errorf("failed to recalculate tournament standings for Group %s: %w", groupCode, err)
		}

		return e.Next()
	})
}

func updateStandingsForGroup(txApp core.App, groupCode string) error {
	// Retrieve teams in the group
	teams, err := txApp.FindRecordsByFilter(
		"teams_id",
		"group_code = {:groupCode}",
		"",
		0,
		0,
		dbx.Params{"groupCode": groupCode},
	)
	if err != nil {
		return err
	}

	// Retrieve group matches that are finished
	matches, err := txApp.FindRecordsByFilter(
		"matches_id",
		"phase = 'group' && group_code = {:groupCode} && status = 'finished'",
		"",
		0,
		0,
		dbx.Params{"groupCode": groupCode},
	)
	if err != nil {
		return err
	}

	type TeamStats struct {
		Record *core.Record
		Pts    int
		GF     int
		GA     int
		GD     int
	}
	statsMap := make(map[string]*TeamStats)
	for _, t := range teams {
		statsMap[t.Id] = &TeamStats{Record: t}
	}

	for _, m := range matches {
		homeID := m.GetString("home_team")
		awayID := m.GetString("away_team")
		homeScore := m.GetInt("score_home")
		awayScore := m.GetInt("score_away")

		hStats, okH := statsMap[homeID]
		aStats, okA := statsMap[awayID]

		if !okH || !okA {
			continue
		}

		hStats.GF += homeScore
		hStats.GA += awayScore
		hStats.GD = hStats.GF - hStats.GA

		aStats.GF += awayScore
		aStats.GA += homeScore
		aStats.GD = aStats.GF - aStats.GA

		if homeScore > awayScore {
			hStats.Pts += 3
		} else if awayScore > homeScore {
			aStats.Pts += 3
		} else {
			hStats.Pts += 1
			aStats.Pts += 1
		}
	}

	sortedStats := make([]*TeamStats, 0, len(statsMap))
	for _, v := range statsMap {
		sortedStats = append(sortedStats, v)
	}

	// Clean integer indexing sort implementation
	sort.Slice(sortedStats, func(i, j int) bool {
		sI := sortedStats[i]
		sJ := sortedStats[j]
		if sI.Pts != sJ.Pts {
			return sI.Pts > sJ.Pts
		}
		if sI.GD != sJ.GD {
			return sI.GD > sJ.GD
		}
		return sI.GF > sJ.GF
	})

	for index, s := range sortedStats {
		rank := index + 1
		s.Record.Set("group_points", s.Pts)
		s.Record.Set("goals_for", s.GF)
		s.Record.Set("goals_against", s.GA)
		s.Record.Set("goal_difference", s.GD)
		s.Record.Set("group_rank", rank)

		if err := txApp.Save(s.Record); err != nil {
			return err
		}
	}

	return nil
}
