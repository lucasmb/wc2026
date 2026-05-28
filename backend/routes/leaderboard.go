package routes

import (
	"net/http"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

type LeaderboardRow struct {
	UserID      string `db:"id" json:"userId"`
	Username    string `db:"username" json:"username"`
	AvatarUrl   string `db:"avatar_url" json:"avatarUrl"`
	TotalPoints int    `db:"total_points" json:"totalPoints"`
	Rank        int    `db:"rank" json:"rank"`
}

func RegisterLeaderboard(app core.App, se *core.ServeEvent) {
	se.Router.GET("/api/wc/leaderboard/{groupId}", func(e *core.RequestEvent) error {
		// Enforce standard Auth
		if e.Auth == nil {
			return apis.NewUnauthorizedError("Authentication required", nil)
		}

		groupID := e.Request.PathValue("groupId")

		// 1. Verify group membership prior to disclosure
		var count int
		err := app.DB().
			Select("COUNT(*)").
			From("group_members").
			Where(dbx.HashExp{"prediction_group": groupID, "user": e.Auth.Id}).
			Row(&count)

		if err != nil || count == 0 {
			return apis.NewForbiddenError("You are not a member of this prediction group", nil)
		}

		// 2. Query rankings and join core user data fields cleanly
		var rows []LeaderboardRow
		err = app.DB().
			Select("users.id", "users.username", "users.avatar_url", "group_members.total_points", "group_members.rank").
			From("group_members").
			LeftJoin("users", dbx.NewExp("group_members.user = users.id")).
			Where(dbx.HashExp{"group_members.prediction_group": groupID}).
			OrderBy("group_members.rank ASC, users.username ASC").
			All(&rows)

		if err != nil {
			return apis.NewBadRequestError("Could not calculate leaderboard standings", err)
		}

		return e.JSON(http.StatusOK, rows)
	})
}
