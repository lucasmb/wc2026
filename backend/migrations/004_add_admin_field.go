package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/migrations"
)

func init() {
	migrations.Register(func(app core.App) error {
		usersColl, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		hasIsAdmin := false
		for _, f := range usersColl.Fields {
			if f.GetName() == "is_admin" {
				hasIsAdmin = true
				break
			}
		}
		if !hasIsAdmin {
			usersColl.Fields.Add(&core.BoolField{
				Name: "is_admin",
			})
			if err := app.Save(usersColl); err != nil {
				return err
			}
		}

		return nil
	}, nil)
}
