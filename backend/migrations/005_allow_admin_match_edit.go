package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	migrations.Register(func(app core.App) error {
		matchesColl, err := app.FindCollectionByNameOrId("matches_id")
		if err != nil {
			return err
		}

		matchesColl.UpdateRule = types.Pointer("@request.auth.is_admin = true")

		if err := app.Save(matchesColl); err != nil {
			return err
		}

		return nil
	}, nil)
}
