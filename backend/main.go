package main

import (
	"fmt"
	"log"
	"os"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/spf13/cobra"

	"wc2026/hooks"
	_ "wc2026/migrations"
	"wc2026/routes"
)

func main() {
	app := pocketbase.New()

	// Auto-run migrations on start
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: true,
	})

	// 1. Register custom Cobra simulation CLI command (Exits cleanly when finished)
	app.RootCmd.AddCommand(&cobra.Command{
		Use:   "simulate",
		Short: "Runs a full tournament simulation from start to finish",
		Run: func(cmd *cobra.Command, args []string) {
			// Bootstrap database hooks to run during simulate operations
			hooks.RegisterTournamentHooks(app)
			hooks.RegisterScoringHooks(app)
			hooks.RegisterPredictionHooks(app)

			// Execute simulation (blocking call, exits when done)
			runSimulation(app)
		},
	})

	// 2. Register custom Cobra database records reset CLI command (Preserves OAuth/System config)
	app.RootCmd.AddCommand(&cobra.Command{
		Use:   "reset",
		Short: "Clears user prediction records and resets match scores while preserving system/OAuth settings",
		Run: func(cmd *cobra.Command, args []string) {
			err := app.RunInTransaction(func(txApp core.App) error {
				fmt.Println("🧹 Resetting prediction database records...")

				// Wipe out user predictions, group memberships, and groups
				_, _ = txApp.DB().Delete("predictions", nil).Execute()
				_, _ = txApp.DB().Delete("group_members", nil).Execute()
				_, _ = txApp.DB().Delete("prediction_groups", nil).Execute()

				// Delete any non-group stage matches (knockouts)
				_, _ = txApp.DB().Delete("matches", dbx.NewExp("phase != 'group'")).Execute()

				// Reset all group matches back to upcoming using the safe Record API
				matchesToReset, err := txApp.FindAllRecords("matches_id")
				if err == nil {
					for _, m := range matchesToReset {
						m.Set("score_home", nil)
						m.Set("score_away", nil)
						m.Set("status", "upcoming")
						m.Set("winner", nil)
						if err := txApp.Save(m); err != nil {
							return fmt.Errorf("failed resetting match ID %s: %w", m.Id, err)
						}
					}
				}

				// Reset team statistics
				_, err = txApp.DB().Update("teams", dbx.Params{
					"group_points":    0,
					"goals_for":       0,
					"goals_against":   0,
					"goal_difference": 0,
					"group_rank":      1,
				}, nil).Execute()
				if err != nil {
					return fmt.Errorf("failed resetting team stats: %w", err)
				}

				fmt.Println("✅ Database records cleanly reset! Match fixtures and team points restored to pristine states.")
				return nil
			})

			if err != nil {
				log.Fatalf("❌ Database reset failed: %v", err)
			}
		},
	})

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))

		routes.RegisterLeaderboard(app, se)
		routes.RegisterSync(app, se)
		routes.RegisterExternalAPI(app, se)

		return se.Next()
	})

	hooks.RegisterUserHooks(app)
	hooks.RegisterTournamentHooks(app)
	hooks.RegisterScoringHooks(app)
	hooks.RegisterPredictionHooks(app)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
