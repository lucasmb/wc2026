package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/migrations"
)

func init() {
	migrations.Register(func(app core.App) error {
		// 1. Fetch Teams Collection Reference
		teamsColl, err := app.FindCollectionByNameOrId("teams_id")
		if err != nil {
			return err
		}

		// 2. Define the 48 Qualified Teams (Groups A to L)
		teamsData := []map[string]any{
			// Group A
			{"code": "MEX", "name": "Mexico", "group_code": "A", "flag_url": "https://flagcdn.com/mx.svg"},
			{"code": "RSA", "name": "South Africa", "group_code": "A", "flag_url": "https://flagcdn.com/za.svg"},
			{"code": "KOR", "name": "Korea Republic", "group_code": "A", "flag_url": "https://flagcdn.com/kr.svg"},
			{"code": "CZE", "name": "Czechia", "group_code": "A", "flag_url": "https://flagcdn.com/cz.svg"},

			// Group B
			{"code": "CAN", "name": "Canada", "group_code": "B", "flag_url": "https://flagcdn.com/ca.svg"},
			{"code": "BIH", "name": "Bosnia and Herzegovina", "group_code": "B", "flag_url": "https://flagcdn.com/ba.svg"},
			{"code": "QAT", "name": "Qatar", "group_code": "B", "flag_url": "https://flagcdn.com/qa.svg"},
			{"code": "SUI", "name": "Switzerland", "group_code": "B", "flag_url": "https://flagcdn.com/ch.svg"},

			// Group C
			{"code": "BRA", "name": "Brazil", "group_code": "C", "flag_url": "https://flagcdn.com/br.svg"},
			{"code": "MAR", "name": "Morocco", "group_code": "C", "flag_url": "https://flagcdn.com/ma.svg"},
			{"code": "HAI", "name": "Haiti", "group_code": "C", "flag_url": "https://flagcdn.com/ht.svg"},
			{"code": "SCO", "name": "Scotland", "group_code": "C", "flag_url": "https://flagcdn.com/gb-sct.svg"},

			// Group D
			{"code": "USA", "name": "United States", "group_code": "D", "flag_url": "https://flagcdn.com/us.svg"},
			{"code": "PAR", "name": "Paraguay", "group_code": "D", "flag_url": "https://flagcdn.com/py.svg"},
			{"code": "AUS", "name": "Australia", "group_code": "D", "flag_url": "https://flagcdn.com/au.svg"},
			{"code": "TUR", "name": "Türkiye", "group_code": "D", "flag_url": "https://flagcdn.com/tr.svg"},

			// Group E
			{"code": "GER", "name": "Germany", "group_code": "E", "flag_url": "https://flagcdn.com/de.svg"},
			{"code": "CUW", "name": "Curaçao", "group_code": "E", "flag_url": "https://flagcdn.com/cw.svg"},
			{"code": "CIV", "name": "Côte d'Ivoire", "group_code": "E", "flag_url": "https://flagcdn.com/ci.svg"},
			{"code": "ECU", "name": "Ecuador", "group_code": "E", "flag_url": "https://flagcdn.com/ec.svg"},

			// Group F
			{"code": "NED", "name": "Netherlands", "group_code": "F", "flag_url": "https://flagcdn.com/nl.svg"},
			{"code": "JPN", "name": "Japan", "group_code": "F", "flag_url": "https://flagcdn.com/jp.svg"},
			{"code": "ALB", "name": "Albania", "group_code": "F", "flag_url": "https://flagcdn.com/al.svg"},
			{"code": "TUN", "name": "Tunisia", "group_code": "F", "flag_url": "https://flagcdn.com/tn.svg"},

			// Group G
			{"code": "BEL", "name": "Belgium", "group_code": "G", "flag_url": "https://flagcdn.com/be.svg"},
			{"code": "EGY", "name": "Egypt", "group_code": "G", "flag_url": "https://flagcdn.com/eg.svg"},
			{"code": "IRN", "name": "IR Iran", "group_code": "G", "flag_url": "https://flagcdn.com/ir.svg"},
			{"code": "NZL", "name": "New Zealand", "group_code": "G", "flag_url": "https://flagcdn.com/nz.svg"},

			// Group H
			{"code": "ESP", "name": "Spain", "group_code": "H", "flag_url": "https://flagcdn.com/es.svg"},
			{"code": "CPV", "name": "Cabo Verde", "group_code": "H", "flag_url": "https://flagcdn.com/cv.svg"},
			{"code": "KSA", "name": "Saudi Arabia", "group_code": "H", "flag_url": "https://flagcdn.com/sa.svg"},
			{"code": "URU", "name": "Uruguay", "group_code": "H", "flag_url": "https://flagcdn.com/uy.svg"},

			// Group I
			{"code": "FRA", "name": "France", "group_code": "I", "flag_url": "https://flagcdn.com/fr.svg"},
			{"code": "SEN", "name": "Senegal", "group_code": "I", "flag_url": "https://flagcdn.com/sn.svg"},
			{"code": "BOL", "name": "Bolivia", "group_code": "I", "flag_url": "https://flagcdn.com/bo.svg"},
			{"code": "NOR", "name": "Norway", "group_code": "I", "flag_url": "https://flagcdn.com/no.svg"},

			// Group J
			{"code": "ARG", "name": "Argentina", "group_code": "J", "flag_url": "https://flagcdn.com/ar.svg"},
			{"code": "ALG", "name": "Algeria", "group_code": "J", "flag_url": "https://flagcdn.com/dz.svg"},
			{"code": "AUT", "name": "Austria", "group_code": "J", "flag_url": "https://flagcdn.com/at.svg"},
			{"code": "JOR", "name": "Jordan", "group_code": "J", "flag_url": "https://flagcdn.com/jo.svg"},

			// Group K
			{"code": "POR", "name": "Portugal", "group_code": "K", "flag_url": "https://flagcdn.com/pt.svg"},
			{"code": "COD", "name": "Congo DR", "group_code": "K", "flag_url": "https://flagcdn.com/cd.svg"},
			{"code": "UZB", "name": "Uzbekistan", "group_code": "K", "flag_url": "https://flagcdn.com/uz.svg"},
			{"code": "COL", "name": "Colombia", "group_code": "K", "flag_url": "https://flagcdn.com/co.svg"},

			// Group L
			{"code": "ENG", "name": "England", "group_code": "L", "flag_url": "https://flagcdn.com/gb-eng.svg"},
			{"code": "CRO", "name": "Croatia", "group_code": "L", "flag_url": "https://flagcdn.com/hr.svg"},
			{"code": "GHA", "name": "Ghana", "group_code": "L", "flag_url": "https://flagcdn.com/gh.svg"},
			{"code": "PAN", "name": "Panama", "group_code": "L", "flag_url": "https://flagcdn.com/pa.svg"},
		}

		// Insert each team record in a single transaction-safe batch
		return app.RunInTransaction(func(txApp core.App) error {
			for _, team := range teamsData {
				record := core.NewRecord(teamsColl)
				record.Set("name", team["name"])
				record.Set("code", team["code"])
				record.Set("flag_url", team["flag_url"])
				record.Set("group_code", team["group_code"])
				record.Set("group_points", 0)
				record.Set("goals_for", 0)
				record.Set("goals_against", 0)
				record.Set("goal_difference", 0)
				record.Set("group_rank", 1)

				if err := txApp.Save(record); err != nil {
					return err
				}
			}

			// 3. Create Default Tournament settings record
			settingsColl, err := txApp.FindCollectionByNameOrId("settings_id")
			if err != nil {
				return err
			}

			settingsRecord := core.NewRecord(settingsColl)
			settingsRecord.Id = "tournsettings26" // Safe 15-char static ID
			settingsRecord.Set("current_phase", "group")
			settingsRecord.Set("predictions_open", true)

			return txApp.Save(settingsRecord)
		})
	}, func(app core.App) error {
		// Rollback logic
		return app.RunInTransaction(func(txApp core.App) error {
			// Delete the settings record
			settingsRec, err := txApp.FindRecordById("settings_id", "tournsettings26")
			if err == nil {
				if err := txApp.Delete(settingsRec); err != nil {
					return err
				}
			}

			// Empty all teams
			teams, err := txApp.FindAllRecords("teams_id")
			if err == nil {
				for _, team := range teams {
					_ = txApp.Delete(team)
				}
			}
			return nil
		})
	})
}
