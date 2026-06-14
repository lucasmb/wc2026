package routes

import (
	"fmt"
	"net/http"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type LeaderboardRow struct {
	UserID      string `json:"userId"`
	Username    string `json:"username"`
	AvatarUrl   string `json:"avatarUrl"`
	Avatar      string `json:"avatar"`
	TotalPoints int    `json:"totalPoints"`
	Rank        int    `json:"rank"`
}

func RegisterLeaderboard(app core.App, se *core.ServeEvent) {
	se.Router.GET("/api/wc/leaderboard/{groupId}", func(e *core.RequestEvent) error {
		// 1. Enforce Standard Auth
		if e.Auth == nil {
			return apis.NewUnauthorizedError("Authentication required", nil)
		}

		groupID := e.Request.PathValue("groupId")

		// 2. Safe check group membership via high-level ORM
		member, err := app.FindFirstRecordByFilter(
			"group_members_id",
			"prediction_group = {:groupId} && user = {:userId}",
			dbx.Params{"groupId": groupID, "userId": e.Auth.Id},
		)
		if err != nil || member == nil {
			return apis.NewForbiddenError("You are not a member of this prediction group", nil)
		}

		// 3. Fetch all group members ordered by Rank ascending
		members, err := app.FindRecordsByFilter(
			"group_members_id",
			"prediction_group = {:groupId}",
			"rank", // Sort by rank ASC
			0,
			0,
			dbx.Params{"groupId": groupID},
		)
		if err != nil {
			return apis.NewBadRequestError("Could not retrieve group members", err)
		}

		// 4. Expand the "user" relation to retrieve user profiles
		errs := app.ExpandRecords(members, []string{"user"}, nil)
		if len(errs) > 0 {
			// Log but do not fail, fall back to "Unknown" username if expansion partially slips
			fmt.Printf("⚠️ Expansion warnings on leaderboard users: %v\n", errs)
		}

		// 5. Structure payload to match frontend JSON expectations
		rows := make([]LeaderboardRow, 0, len(members))
		for _, m := range members {
			userRec := m.ExpandedOne("user")
			username := "Unknown"
			avatarUrl := ""
			avatar := ""

			if userRec != nil {
				username = userRec.GetString("username")
				avatarUrl = userRec.GetString("avatar_url")
				avatar = userRec.GetString("avatar")
			}

			rows = append(rows, LeaderboardRow{
				UserID:      m.GetString("user"),
				Username:    username,
				AvatarUrl:   avatarUrl,
				Avatar:      avatar,
				TotalPoints: m.GetInt("total_points"),
				Rank:        m.GetInt("rank"),
			})
		}

		return e.JSON(http.StatusOK, rows)
	})
}
