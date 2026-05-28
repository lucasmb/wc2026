package main

import (
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	"wc2026/hooks"
	_ "wc2026/migrations"
	"wc2026/routes"
)

func main() {
	app := pocketbase.New()

	// Register migration commands and enable automigration during dev
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: true,
	})

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// Serve Quasar SPA build from pb_public directory
		se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))

		// Register custom API routes
		routes.RegisterLeaderboard(app, se)
		routes.RegisterSync(app, se)

		return se.Next()
	})

	// Register lifecycle hook interceptors
	hooks.RegisterTournamentHooks(app)
	hooks.RegisterScoringHooks(app)
	hooks.RegisterPredictionHooks(app)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
