package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/types"
)

// Represents qualified teams during bracket calculation
type QualifiedTeam struct {
	ID        string
	Name      string
	GroupCode string
	Points    int
	GD        int
	GF        int
}

func runSimulation(app core.App) {
	fmt.Println("==================================================")
	fmt.Println("⚽ STARTING FIFA WORLD CUP 2026 FULL SIMULATION ⚽")
	fmt.Println("==================================================")

	err := app.RunInTransaction(func(txApp core.App) error {
		// 1. Database Cleanup
		fmt.Println("🧹 Cleaning out database simulation states...")
		_, _ = txApp.DB().Delete("predictions", nil).Execute()
		_, _ = txApp.DB().Delete("group_members", nil).Execute()
		_, _ = txApp.DB().Delete("prediction_groups", nil).Execute()

		// Delete any previously generated knockout matches to start fresh
		_, _ = txApp.DB().Delete("matches", dbx.NewExp("phase != 'group'")).Execute()

		// Reset all group matches to upcoming/nil scores using the safe Record API
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

		// 2. Setup Test Group and Users
		usersColl, err := txApp.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		simUsers := []string{"Messi", "Maradona", "Scaloni"}
		userRecords := make(map[string]*core.Record)

		for _, name := range simUsers {
			email := fmt.Sprintf("%s@sim.com", name)
			user, err := txApp.FindFirstRecordByFilter("users", "email = {:email}", dbx.Params{"email": email})
			if err != nil {
				user = core.NewRecord(usersColl)
				user.Set("email", email)
				user.Set("username", name)
				user.SetVerified(true)
				user.SetPassword("simpass123")
				if err := txApp.Save(user); err != nil {
					return err
				}
			}
			userRecords[name] = user
		}

		// Create Prediction Group
		groupColl, err := txApp.FindCollectionByNameOrId("prediction_groups_id")
		if err != nil {
			return err
		}
		group := core.NewRecord(groupColl)
		group.Id = "simpredictiongp"
		group.Set("name", "Simulation Prode")
		group.Set("owner", userRecords["Messi"].Id)
		group.Set("invite_code", "SIMCODE1")
		if err := txApp.Save(group); err != nil {
			return err
		}

		// Create Memberships
		memberColl, err := txApp.FindCollectionByNameOrId("group_members_id")
		if err != nil {
			return err
		}
		for _, u := range userRecords {
			m := core.NewRecord(memberColl)
			m.Set("prediction_group", group.Id)
			m.Set("user", u.Id)
			m.Set("total_points", 0)
			m.Set("rank", 1)
			if err := txApp.Save(m); err != nil {
				return err
			}
		}

		// 3. Generate Predictions and Simulate Group Stage
		matches, err := txApp.FindRecordsByFilter("matches_id", "phase = 'group'", "match_number", 0, 0, nil)
		if err != nil {
			return err
		}

		predColl, err := txApp.FindCollectionByNameOrId("predictions_id")
		if err != nil {
			return err
		}

		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		fmt.Println("🎲 Generating unique predictions and playing 72 Group Matches...")
		for _, m := range matches {
			for _, u := range userRecords {
				p := core.NewRecord(predColl)
				p.Set("user", u.Id)
				p.Set("match", m.Id)
				p.Set("prediction_group", group.Id)
				p.Set("predicted_home", r.Intn(4))
				p.Set("predicted_away", r.Intn(4))
				p.Set("points_awarded", 0)
				p.Set("is_locked", true)
				if err := txApp.Save(p); err != nil {
					return err
				}
			}

			// Play the match
			scoreH := r.Intn(4)
			scoreA := r.Intn(4)

			m.Set("score_home", scoreH)
			m.Set("score_away", scoreA)
			m.Set("status", "finished")
			if err := txApp.Save(m); err != nil {
				return err
			}
		}

		// 4. ADVANCE PHASE ENGINE: Calculate Round of 32 Qualified Teams
		fmt.Println("⚙️ Calculating qualified teams for the Round of 32 (R32)...")
		teams, err := txApp.FindRecordsByFilter("teams_id", "", "group_code, group_rank", 0, 0, nil)
		if err != nil {
			return err
		}

		qualifiedMap := make([]QualifiedTeam, 0, 32)
		thirdPlaceTeams := make([]QualifiedTeam, 0, 12)

		for _, t := range teams {
			qTeam := QualifiedTeam{
				ID:        t.Id,
				Name:      t.GetString("name"),
				GroupCode: t.GetString("group_code"),
				Points:    t.GetInt("group_points"),
				GD:        t.GetInt("goal_difference"),
				GF:        t.GetInt("goals_for"),
			}

			rank := t.GetInt("group_rank")
			if rank <= 2 {
				// Top 2 of each group qualify automatically (24 teams)
				qualifiedMap = append(qualifiedMap, qTeam)
			} else if rank == 3 {
				// 3rd placed teams go to wildcard evaluation list
				thirdPlaceTeams = append(thirdPlaceTeams, qTeam)
			}
		}

		// Sort 3rd placed teams: Points -> Goal Difference -> Goals For
		sort.Slice(thirdPlaceTeams, func(i, j int) bool {
			tI := thirdPlaceTeams[i]
			tJ := thirdPlaceTeams[j]
			if tI.Points != tJ.Points {
				return tI.Points > tJ.Points
			}
			if tI.GD != tJ.GD {
				return tI.GD > tJ.GD
			}
			return tI.GF > tJ.GF
		})

		// Take the best 8 third-place teams (8 wildcard teams)
		for i := 0; i < 8 && i < len(thirdPlaceTeams); i++ {
			qualifiedMap = append(qualifiedMap, thirdPlaceTeams[i])
		}

		fmt.Printf("✅ Selected 32 Knockout Teams (24 top-2 teams + 8 wildcards). Example: %s, %s, %s\n",
			qualifiedMap[0].Name, qualifiedMap[1].Name, qualifiedMap[2].Name)

		// 5. Create 16 Round of 32 matches in 'upcoming' status
		fmt.Println("📅 Scheduling 16 Round of 32 matches (Matches 73 to 88)...")
		matchesColl, err := txApp.FindCollectionByNameOrId("matches_id")
		if err != nil {
			return err
		}

		for i := 0; i < 16; i++ {
			homeTeam := qualifiedMap[i*2]
			awayTeam := qualifiedMap[(i*2)+1]

			m := core.NewRecord(matchesColl)
			m.Set("phase", "r32")
			m.Set("match_number", 73+i)
			m.Set("home_team", homeTeam.ID)
			m.Set("away_team", awayTeam.ID)
			m.Set("venue", "Host Knockout Stadium")
			m.Set("city", "Host City")
			m.Set("status", "upcoming") // Leaves them unplayed so users can predict on frontend!
			m.Set("match_day", 4)       // Matchday 4 represents knockouts

			// Set kickoff 15 days in the future
			ko, _ := types.ParseDateTime(time.Now().Add(15 * 24 * time.Hour))
			m.Set("kickoff", ko)

			if err := txApp.Save(m); err != nil {
				return err
			}
		}

		// 6. Display Standings
		fmt.Println("\n🏁 Group Stage Completed! Standings updated.")
		var standings []struct {
			Username string `db:"username"`
			Points   int    `db:"total_points"`
			Rank     int    `db:"rank"`
		}
		err = txApp.DB().
			Select("users.username", "group_members.total_points", "group_members.rank").
			From("group_members").
			LeftJoin("users", dbx.NewExp("group_members.user = users.id")).
			Where(dbx.HashExp{"group_members.prediction_group": "simpredictiongp"}).
			OrderBy("group_members.rank ASC").
			All(&standings)

		if err != nil {
			return err
		}

		for _, s := range standings {
			fmt.Printf("🏆 Rank #%d: %-12s - %3d Pts\n", s.Rank, s.Username, s.Points)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("❌ Simulation failed: %v\n", err)
		return
	}

	fmt.Println("\n==================================================")
	fmt.Println("🎉 STEP Completed successfully! Server remains live.")
	fmt.Println("==================================================")
}
