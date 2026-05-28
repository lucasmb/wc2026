package migrations

import (
	"fmt"
	"time"

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

		// Retrieve all seeded teams to resolve relations
		teams, err := app.FindAllRecords("teams_id")
		if err != nil {
			return err
		}

		// Organize teams by their assigned Group Code
		groupsMap := make(map[string][]*core.Record)
		for _, team := range teams {
			gCode := team.GetString("group_code")
			groupsMap[gCode] = append(groupsMap[gCode], team)
		}

		groupCodes := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L"}
		matchCounter := 1

		// Base Tournament Start Date: June 11, 2026 (12:00 UTC)
		baseStartDate := time.Date(2026, 6, 11, 12, 0, 0, 0, time.UTC)

		venues := []string{
			"Estadio Azteca (Mexico City)", "MetLife Stadium (East Rutherford)",
			"SoFi Stadium (Los Angeles)", "BC Place (Vancouver)",
			"Hard Rock Stadium (Miami)", "AT&T Stadium (Dallas)",
			"Arrowhead Stadium (Kansas City)", "Mercedes-Benz Stadium (Atlanta)",
			"Gillette Stadium (Boston)", "Lincoln Financial Field (Philadelphia)",
			"NRG Stadium (Houston)", "Lumen Field (Seattle)",
		}

		return app.RunInTransaction(func(txApp core.App) error {
			for gIdx, code := range groupCodes {
				groupTeams := groupsMap[code]
				if len(groupTeams) < 4 {
					return fmt.Errorf("group %s does not contain 4 seeded teams", code)
				}

				T1, T2, T3, T4 := groupTeams[0], groupTeams[1], groupTeams[2], groupTeams[3]

				// Programmatic pairings to generate three matchdays for the group
				// Matchday 1
				m1Home, m1Away := T1, T2
				m2Home, m2Away := T3, T4
				// Matchday 2
				m3Home, m3Away := T1, T3
				m4Home, m4Away := T2, T4
				// Matchday 3 (Simultaneous)
				m5Home, m5Away := T1, T4
				m6Home, m6Away := T2, T3

				// Calculate realistic staggered kickoff times
				staggerOffset := time.Duration(gIdx) * 12 * time.Hour

				// Matchday 1 dates
				m1Date := baseStartDate.Add(staggerOffset)
				m2Date := m1Date.Add(4 * time.Hour)

				// Matchday 2 dates (approx 7 days later)
				m3Date := m1Date.Add(7 * 24 * time.Hour)
				m4Date := m3Date.Add(4 * time.Hour)

				// Matchday 3 dates (approx 14 days later - strictly simultaneous)
				m5Date := m1Date.Add(14 * 24 * time.Hour)
				m6Date := m5Date // STRICT SIMULTANEOUS KICKOFF

				fixtures := []struct {
					home     *core.Record
					away     *core.Record
					kickoff  time.Time
					matchDay int
				}{
					{m1Home, m1Away, m1Date, 1},
					{m2Home, m2Away, m2Date, 1},
					{m3Home, m3Away, m3Date, 2},
					{m4Home, m4Away, m4Date, 2},
					{m5Home, m5Away, m5Date, 3}, // Simultaneous pair
					{m6Home, m6Away, m6Date, 3},
				}

				for _, f := range fixtures {
					record := core.NewRecord(matchesColl)
					record.Set("phase", "group")
					record.Set("group_code", code)
					record.Set("match_number", matchCounter)
					record.Set("home_team", f.home.Id)
					record.Set("away_team", f.away.Id)

					// Parse the time.Time cleanly to types.DateTime
					dt, err := types.ParseDateTime(f.kickoff)
					if err != nil {
						return fmt.Errorf("failed to parse kickoff time: %w", err)
					}
					record.Set("kickoff", dt)

					record.Set("venue", venues[(matchCounter-1)%len(venues)])
					record.Set("city", "Host City")
					record.Set("status", "upcoming")
					record.Set("match_day", f.matchDay)

					if err := txApp.Save(record); err != nil {
						return err
					}
					matchCounter++
				}
			}
			return nil
		})
	}, func(app core.App) error {
		// Rollback logic
		return app.RunInTransaction(func(txApp core.App) error {
			matches, err := txApp.FindAllRecords("matches_id")
			if err == nil {
				for _, match := range matches {
					_ = txApp.Delete(match)
				}
			}
			return nil
		})
	})
}
