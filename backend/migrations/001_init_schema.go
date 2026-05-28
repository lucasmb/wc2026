package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	migrations.Register(func(app core.App) error {
		// 1. Retrieve the existing system "users" collection
		usersColl, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		// Ensure users has avatar_url
		hasAvatarUrl := false
		for _, f := range usersColl.Fields {
			if f.GetName() == "avatar_url" {
				hasAvatarUrl = true
				break
			}
		}
		if !hasAvatarUrl {
			usersColl.Fields.Add(&core.URLField{
				Name:        "avatar_url",
				Presentable: true,
			})
			if err := app.Save(usersColl); err != nil {
				return err
			}
		}

		// 2. Collection: teams (Static ID: "teams_id")
		teams := core.NewBaseCollection("teams")
		teams.Id = "teams_id"
		teams.ListRule = types.Pointer("") // Public read
		teams.ViewRule = types.Pointer("")
		teams.Fields.Add(
			&core.TextField{Name: "name", Required: true},
			&core.TextField{Name: "code", Required: true, Min: 3, Max: 3},
			&core.URLField{Name: "flag_url"},
			&core.TextField{Name: "group_code", Max: 1},
			&core.NumberField{Name: "group_points"},
			&core.NumberField{Name: "goals_for"},
			&core.NumberField{Name: "goals_against"},
			&core.NumberField{Name: "goal_difference"},
			&core.NumberField{Name: "group_rank"},
		)
		teams.Indexes = types.JSONArray[string]{
			"CREATE UNIQUE INDEX idx_team_code ON teams (code)",
		}
		if err := app.Save(teams); err != nil {
			return err
		}

		// 3. Collection: matches (Static ID: "matches_id")
		matches := core.NewBaseCollection("matches")
		matches.Id = "matches_id"
		matches.ListRule = types.Pointer("")
		matches.ViewRule = types.Pointer("")
		matches.Fields.Add(
			&core.SelectField{
				Name:      "phase",
				Required:  true,
				Values:    []string{"group", "r32", "r16", "qf", "sf", "final", "third"},
				MaxSelect: 1,
			},
			&core.TextField{Name: "group_code", Max: 1},
			&core.NumberField{Name: "match_number", Required: true},
			&core.RelationField{
				Name:         "home_team",
				CollectionId: "teams_id", // Clean static reference
				MaxSelect:    1,
			},
			&core.RelationField{
				Name:         "away_team",
				CollectionId: "teams_id",
				MaxSelect:    1,
			},
			&core.DateField{Name: "kickoff", Required: true},
			&core.TextField{Name: "venue"},
			&core.TextField{Name: "city"},
			&core.NumberField{Name: "score_home"},
			&core.NumberField{Name: "score_away"},
			&core.RelationField{
				Name:         "winner",
				CollectionId: "teams_id",
				MaxSelect:    1,
			},
			&core.SelectField{
				Name:      "status",
				Required:  true,
				Values:    []string{"upcoming", "live", "finished"},
				MaxSelect: 1,
			},
			&core.NumberField{Name: "match_day"},
		)
		if err := app.Save(matches); err != nil {
			return err
		}

		// 4. Collection: prediction_groups (Static ID: "prediction_groups_id")
		predGroups := core.NewBaseCollection("prediction_groups")
		predGroups.Id = "prediction_groups_id"
		predGroups.ListRule = types.Pointer("@request.auth.id != ''")
		predGroups.ViewRule = types.Pointer("@request.auth.id != ''")
		predGroups.CreateRule = types.Pointer("@request.auth.id != ''")
		predGroups.UpdateRule = types.Pointer("@request.auth.id = owner")
		predGroups.DeleteRule = types.Pointer("@request.auth.id = owner")
		predGroups.Fields.Add(
			&core.TextField{Name: "name", Required: true},
			&core.RelationField{
				Name:         "owner",
				CollectionId: usersColl.Id,
				Required:     true,
				MaxSelect:    1,
			},
			&core.TextField{Name: "invite_code"},
			&core.DateField{Name: "invite_expires_at"},
			&core.BoolField{Name: "is_public"},
		)
		predGroups.Indexes = types.JSONArray[string]{
			"CREATE UNIQUE INDEX idx_group_invite_code ON prediction_groups (invite_code)",
		}
		if err := app.Save(predGroups); err != nil {
			return err
		}

		// 5. Collection: group_members (Static ID: "group_members_id")
		members := core.NewBaseCollection("group_members")
		members.Id = "group_members_id"
		members.ListRule = types.Pointer("@request.auth.id != ''")
		members.ViewRule = types.Pointer("@request.auth.id != ''")
		members.CreateRule = types.Pointer("@request.auth.id != ''")
		members.DeleteRule = types.Pointer("@request.auth.id = user")
		members.Fields.Add(
			&core.RelationField{
				Name:          "prediction_group",
				CollectionId:  "prediction_groups_id",
				Required:      true,
				MaxSelect:     1,
				CascadeDelete: true,
			},
			&core.RelationField{
				Name:          "user",
				CollectionId:  usersColl.Id,
				Required:      true,
				MaxSelect:     1,
				CascadeDelete: true,
			},
			&core.NumberField{Name: "total_points"},
			&core.NumberField{Name: "rank"},
		)
		members.Indexes = types.JSONArray[string]{
			"CREATE UNIQUE INDEX idx_group_user ON group_members (prediction_group, user)",
		}
		if err := app.Save(members); err != nil {
			return err
		}

		// 6. Collection: predictions (Static ID: "predictions_id")
		predictions := core.NewBaseCollection("predictions")
		predictions.Id = "predictions_id"
		predictions.ListRule = types.Pointer("@request.auth.id = user")
		predictions.ViewRule = types.Pointer("@request.auth.id = user")
		predictions.CreateRule = types.Pointer("@request.auth.id != ''")
		predictions.UpdateRule = types.Pointer("@request.auth.id = user")
		predictions.Fields.Add(
			&core.RelationField{
				Name:          "user",
				CollectionId:  usersColl.Id,
				Required:      true,
				MaxSelect:     1,
				CascadeDelete: true,
			},
			&core.RelationField{
				Name:          "match",
				CollectionId:  "matches_id",
				Required:      true,
				MaxSelect:     1,
				CascadeDelete: true,
			},
			&core.RelationField{
				Name:          "prediction_group",
				CollectionId:  "prediction_groups_id",
				Required:      true,
				MaxSelect:     1,
				CascadeDelete: true,
			},
			&core.NumberField{Name: "predicted_home"},
			&core.NumberField{Name: "predicted_away"},
			&core.NumberField{Name: "points_awarded"},
			&core.BoolField{Name: "is_locked"},
			&core.DateField{Name: "submitted_at"},
		)
		predictions.Indexes = types.JSONArray[string]{
			"CREATE UNIQUE INDEX idx_user_match_group ON predictions (user, match, prediction_group)",
		}
		if err := app.Save(predictions); err != nil {
			return err
		}

		// 7. Collection: bonus_predictions (Static ID: "bonus_predictions_id")
		bonus := core.NewBaseCollection("bonus_predictions")
		bonus.Id = "bonus_predictions_id"
		bonus.ListRule = types.Pointer("@request.auth.id = user")
		bonus.ViewRule = types.Pointer("@request.auth.id = user")
		bonus.CreateRule = types.Pointer("@request.auth.id != ''")
		bonus.UpdateRule = types.Pointer("@request.auth.id = user")
		bonus.Fields.Add(
			&core.RelationField{
				Name:          "user",
				CollectionId:  usersColl.Id,
				Required:      true,
				MaxSelect:     1,
				CascadeDelete: true,
			},
			&core.RelationField{
				Name:          "prediction_group",
				CollectionId:  "prediction_groups_id",
				Required:      true,
				MaxSelect:     1,
				CascadeDelete: true,
			},
			&core.RelationField{
				Name:         "champion",
				CollectionId: "teams_id",
				MaxSelect:    1,
			},
			&core.TextField{Name: "top_scorer_name"},
			&core.NumberField{Name: "points_awarded"},
			&core.BoolField{Name: "is_locked"},
		)
		bonus.Indexes = types.JSONArray[string]{
			"CREATE UNIQUE INDEX idx_bonus_user_group ON bonus_predictions (user, prediction_group)",
		}
		if err := app.Save(bonus); err != nil {
			return err
		}

		// 8. Collection: settings (Static ID: "settings_id")
		settings := core.NewBaseCollection("settings")
		settings.Id = "settings_id"
		settings.ListRule = types.Pointer("")
		settings.ViewRule = types.Pointer("")
		settings.Fields.Add(
			&core.SelectField{
				Name:      "current_phase",
				Required:  true,
				Values:    []string{"group", "r32", "r16", "qf", "sf", "final", "third"},
				MaxSelect: 1,
			},
			&core.BoolField{Name: "predictions_open"},
			&core.DateField{Name: "last_sync_at"},
		)
		if err := app.Save(settings); err != nil {
			return err
		}

		return nil
	}, func(app core.App) error {
		// Rollback order safely
		colls := []string{"settings_id", "bonus_predictions_id", "predictions_id", "group_members_id", "prediction_groups_id", "matches_id", "teams_id"}
		for _, id := range colls {
			c, err := app.FindCollectionByNameOrId(id)
			if err == nil {
				if err := app.Delete(c); err != nil {
					return err
				}
			}
		}
		return nil
	})
}
