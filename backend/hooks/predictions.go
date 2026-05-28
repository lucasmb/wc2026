package hooks

import (
	"time"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func RegisterPredictionHooks(app core.App) {
	// Intercept Prediction Creation Requests
	app.OnRecordCreateRequest("predictions_id").BindFunc(func(e *core.RecordRequestEvent) error {
		if e.HasSuperuserAuth() {
			return e.Next()
		}

		// 1. Force the current logged-in user as the owner
		if e.Auth == nil {
			return apis.NewUnauthorizedError("Authentication is required to place a prediction", nil)
		}
		e.Record.Set("user", e.Auth.Id)
		e.Record.Set("submitted_at", time.Now().UTC())

		// 2. Fetch the target Match to inspect the Kickoff time
		matchID := e.Record.GetString("match")
		match, err := e.App.FindRecordById("matches_id", matchID)
		if err != nil {
			return apis.NewBadRequestError("Invalid match reference", err)
		}

		kickoffVal := match.GetDateTime("kickoff")
		if kickoffVal.Time().Before(time.Now().UTC()) {
			return apis.NewBadRequestError("Predictions are locked. This match has already kicked off.", nil)
		}

		return e.Next()
	})

	// Intercept Prediction Update Requests
	app.OnRecordUpdateRequest("predictions_id").BindFunc(func(e *core.RecordRequestEvent) error {
		if e.HasSuperuserAuth() {
			return e.Next()
		}

		if e.Auth == nil {
			return apis.NewUnauthorizedError("Authentication is required to modify a prediction", nil)
		}

		// Ensure users cannot tamper with prediction ownership
		e.Record.Set("user", e.Auth.Id)
		e.Record.Set("submitted_at", time.Now().UTC())

		// Fetch the match to check kickoff status
		matchID := e.Record.GetString("match")
		match, err := e.App.FindRecordById("matches_id", matchID)
		if err != nil {
			return apis.NewBadRequestError("Invalid match reference", err)
		}

		kickoffVal := match.GetDateTime("kickoff")
		if kickoffVal.Time().Before(time.Now().UTC()) {
			return apis.NewBadRequestError("Predictions are locked. This match has already kicked off.", nil)
		}

		return e.Next()
	})
}
