package routes

import (
	"net/http"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type MatchPredictionRow struct {
	ID             string `json:"id"`
	Match          string `json:"match"`
	PredictedHome  int    `json:"predicted_home"`
	PredictedAway  int    `json:"predicted_away"`
	PointsAwarded  int    `json:"points_awarded"`
	UserName       string `json:"userName"`
	UserAvatarUrl  string `json:"userAvatarUrl"`
	UserAvatar     string `json:"userAvatar"`
}

func RegisterMatchPredictions(app core.App, se *core.ServeEvent) {
	se.Router.GET("/api/wc/match-predictions/{groupId}", func(e *core.RequestEvent) error {
		if e.Auth == nil {
			return apis.NewUnauthorizedError("Authentication required", nil)
		}

		groupID := e.Request.PathValue("groupId")

		member, err := app.FindFirstRecordByFilter(
			"group_members_id",
			"prediction_group = {:groupId} && user = {:userId}",
			dbx.Params{"groupId": groupID, "userId": e.Auth.Id},
		)
		if err != nil || member == nil {
			return apis.NewForbiddenError("You are not a member of this prediction group", nil)
		}

		predictions, err := app.FindRecordsByFilter(
			"predictions_id",
			"prediction_group = {:groupId}",
			"",
			0,
			0,
			dbx.Params{"groupId": groupID},
		)
		if err != nil {
			return apis.NewBadRequestError("Could not retrieve predictions", err)
		}

		userIDs := make(map[string]bool)
		for _, p := range predictions {
			userIDs[p.GetString("user")] = true
		}

		userMap := make(map[string]*core.Record)
		for userID := range userIDs {
			userRec, err := app.FindRecordById("users", userID)
			if err == nil {
				userMap[userID] = userRec
			}
		}

		rows := make([]MatchPredictionRow, 0, len(predictions))
		for _, p := range predictions {
			userID := p.GetString("user")
			userRec := userMap[userID]

			userName := "Unknown"
			userAvatarUrl := ""
			userAvatar := ""

			if userRec != nil {
				userName = userRec.GetString("username")
				userAvatarUrl = userRec.GetString("avatar_url")
				userAvatar = userRec.GetString("avatar")
			}

			rows = append(rows, MatchPredictionRow{
				ID:            p.Id,
				Match:         p.GetString("match"),
				PredictedHome: p.GetInt("predicted_home"),
				PredictedAway: p.GetInt("predicted_away"),
				PointsAwarded: p.GetInt("points_awarded"),
				UserName:      userName,
				UserAvatarUrl: userAvatarUrl,
				UserAvatar:    userAvatar,
			})
		}

		return e.JSON(http.StatusOK, rows)
	})
}
