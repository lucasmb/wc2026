package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type SyncRequestPayload struct {
	MatchNumber int    `json:"match_number"`
	ScoreHome   int    `json:"score_home"`
	ScoreAway   int    `json:"score_away"`
	Status      string `json:"status"`
}

func RegisterSync(app core.App, se *core.ServeEvent) {
	se.Router.POST("/api/wc/sync", func(e *core.RequestEvent) error {
		if !e.HasSuperuserAuth() {
			return apis.NewForbiddenError("Superuser privileges are required to perform synchronization", nil)
		}

		var payloads []SyncRequestPayload
		if err := json.NewDecoder(e.Request.Body).Decode(&payloads); err != nil {
			return apis.NewBadRequestError("Invalid request body payload format", err)
		}

		err := app.RunInTransaction(func(txApp core.App) error {
			for _, payload := range payloads {
				match, err := txApp.FindFirstRecordByFilter("matches_id", "match_number = {:matchNum}", map[string]any{"matchNum": payload.MatchNumber})
				if err != nil {
					continue
				}

				match.Set("score_home", payload.ScoreHome)
				match.Set("score_away", payload.ScoreAway)
				match.Set("status", payload.Status)

				if err := txApp.Save(match); err != nil {
					return fmt.Errorf("failed updating match #%d: %w", payload.MatchNumber, err)
				}
			}

			// Retrieve settings record and update last_sync_at using the types.NowDateTime() helper
			settings, err := txApp.FindRecordById("settings_id", "tournsettings26")
			if err == nil {
				settings.Set("last_sync_at", types.NowDateTime())
				_ = txApp.Save(settings)
			}

			return nil
		})

		if err != nil {
			return apis.NewBadRequestError("Sync process failed midway and rolled back", err)
		}

		return e.JSON(http.StatusOK, map[string]any{
			"success": true,
			"message": "Chronological match stats successfully synced",
		})
	})
}
