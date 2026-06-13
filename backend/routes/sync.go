package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

type SyncRequestPayload struct {
	MatchNumber int    `json:"match_number"`
	ScoreHome   int    `json:"score_home"`
	ScoreAway   int    `json:"score_away"`
	Status      string `json:"status"`
	Swapped     bool   `json:"swapped"`
}

func RegisterSync(app core.App, se *core.ServeEvent) {
	se.Router.POST("/api/wc/sync", func(e *core.RequestEvent) error {
		// 1. Dual-Auth Check: Allow PocketBase superuser OR secure API Key Header
		apiKeyHeader := e.Request.Header.Get("X-WC-API-Key")
		configuredApiKey := os.Getenv("WC_API_KEY")

		hasValidApiKey := (configuredApiKey != "" && apiKeyHeader == configuredApiKey)

		if !e.HasSuperuserAuth() && !hasValidApiKey {
			return apis.NewForbiddenError("Authentication required (Admin session or valid X-WC-API-Key header)", nil)
		}

		// 2. Decode the incoming JSON array of match results
		var payloads []SyncRequestPayload
		if err := json.NewDecoder(e.Request.Body).Decode(&payloads); err != nil {
			return apis.NewBadRequestError("Invalid request body payload format", err)
		}

		// 3. Process updates transactionally
		err := app.RunInTransaction(func(txApp core.App) error {
			for _, payload := range payloads {
				// Locate target match by official match number
				match, err := txApp.FindFirstRecordByFilter(
					"matches_id",
					"match_number = {:matchNum}",
					map[string]any{"matchNum": payload.MatchNumber},
				)
				if err != nil {
					// Skip unmapped matches gracefully
					continue
				}

				// Apply scores and status
				match.Set("score_home", payload.ScoreHome)
				match.Set("score_away", payload.ScoreAway)
				match.Set("status", payload.Status)

				// Saving here automatically fires the scoring hooks (hooks/scoring.go)
				if err := txApp.Save(match); err != nil {
					return fmt.Errorf("failed updating match #%d: %w", payload.MatchNumber, err)
				}
			}

			// Record sync timestamp in global settings
			settings, err := txApp.FindRecordById("settings_id", "tournsettings26")
			if err == nil {
				settings.Set("last_sync_at", types.NowDateTime())
				_ = txApp.Save(settings)
			}

			return nil
		})

		if err != nil {
			return apis.NewBadRequestError("Sync process failed and was rolled back", err)
		}

		return e.JSON(http.StatusOK, map[string]any{
			"success": true,
			"message": "Chronological match stats successfully synchronized",
		})
	})
}
