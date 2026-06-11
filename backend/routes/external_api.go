package routes

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

const (
	wc26BaseURL       = "https://worldcup26.ir"
	wc26SettingsID    = "tournsettings26"
	wc26SettingsEmail = "wc26_api_email"
	wc26SettingsPass  = "wc26_api_password"
	wc26SettingsToken = "wc26_api_token"
)

var (
	cachedToken string
	tokenMu     sync.RWMutex
)

type wc26AuthResponse struct {
	Token string `json:"token"`
}

type ExternalMatchData struct {
	MatchNumber int                    `json:"match_number"`
	Phase       string                 `json:"phase"`
	GroupCode   string                 `json:"group_code"`
	External    map[string]interface{} `json:"external"`
	Local       map[string]interface{} `json:"local"`
}

func RegisterExternalAPI(app core.App, se *core.ServeEvent) {
	se.Router.GET("/api/wc/external/matches", func(e *core.RequestEvent) error {
		if !e.HasSuperuserAuth() && !isAdminUser(e) {
			return apis.NewForbiddenError("Admin access required", nil)
		}

		token, err := getWC26Token(app)
		if err != nil {
			return apis.NewBadRequestError("Failed to authenticate with external API", err)
		}

		req, _ := http.NewRequest("GET", wc26BaseURL+"/get/games", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return apis.NewBadRequestError("Failed to fetch from external API", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusUnauthorized {
			token, err = refreshWC26Token(app)
			if err != nil {
				return apis.NewBadRequestError("Failed to re-authenticate with external API", err)
			}
			req.Header.Set("Authorization", "Bearer "+token)
			resp, err = http.DefaultClient.Do(req)
			if err != nil {
				return apis.NewBadRequestError("Failed to fetch from external API after re-auth", err)
			}
			defer resp.Body.Close()
		}

		body, _ := io.ReadAll(resp.Body)

		var externalData struct {
			Games []map[string]interface{} `json:"games"`
		}
		if err := json.Unmarshal(body, &externalData); err != nil {
			var gamesArray []map[string]interface{}
			if err2 := json.Unmarshal(body, &gamesArray); err2 != nil {
				return apis.NewBadRequestError("Failed to parse external API response", err)
			}
			externalData.Games = gamesArray
		}

		localMatches, err := app.FindAllRecords("matches_id")
		if err != nil {
			return apis.NewBadRequestError("Failed to fetch local matches", err)
		}

		localByNumber := make(map[int]*core.Record)
		for _, m := range localMatches {
			mn := m.GetInt("match_number")
			localByNumber[mn] = m
		}

		teams, _ := app.FindAllRecords("teams_id")
		teamById := make(map[string]string)
		teamByCode := make(map[string]string)
		for _, t := range teams {
			teamById[t.Id] = t.GetString("name")
			teamByCode[t.GetString("code")] = t.GetString("name")
		}

		var result []ExternalMatchData
		for _, game := range externalData.Games {
			matchNum := 0
			if id, ok := game["id"]; ok {
				switch v := id.(type) {
				case string:
					matchNum, _ = strconv.Atoi(v)
				case float64:
					matchNum = int(v)
				}
			}

			localMatch, hasLocal := localByNumber[matchNum]

			phase := "group"
			groupCode := ""
			if t, ok := game["type"]; ok {
				phase = fmt.Sprintf("%v", t)
			}
			if g, ok := game["group"]; ok {
				groupCode = fmt.Sprintf("%v", g)
				if groupCode == "R32" || groupCode == "R16" || groupCode == "QF" || groupCode == "SF" || groupCode == "FINAL" || groupCode == "3RD" {
					groupCode = ""
				}
			}

			local := map[string]interface{}{
				"status":     "unknown",
				"score_home": nil,
				"score_away": nil,
			}
			if hasLocal {
				local["status"] = localMatch.GetString("status")
				local["score_home"] = localMatch.Get("score_home")
				local["score_away"] = localMatch.Get("score_away")

				homeTeamID := localMatch.GetString("home_team")
				awayTeamID := localMatch.GetString("away_team")
				local["home_team_name"] = teamById[homeTeamID]
				local["away_team_name"] = teamById[awayTeamID]
			}

			result = append(result, ExternalMatchData{
				MatchNumber: matchNum,
				Phase:       phase,
				GroupCode:   groupCode,
				External:    game,
				Local:       local,
			})
		}

		return e.JSON(http.StatusOK, result)
	})

	se.Router.POST("/api/wc/external/sync", func(e *core.RequestEvent) error {
		if !e.HasSuperuserAuth() && !isAdminUser(e) {
			return apis.NewForbiddenError("Admin access required", nil)
		}

		var payloads []SyncRequestPayload
		if err := json.NewDecoder(e.Request.Body).Decode(&payloads); err != nil {
			return apis.NewBadRequestError("Invalid request body", err)
		}

		synced := 0
		err := app.RunInTransaction(func(txApp core.App) error {
			for _, payload := range payloads {
				if payload.Status == "upcoming" {
					continue
				}

				match, err := txApp.FindFirstRecordByFilter(
					"matches_id",
					"match_number = {:matchNum}",
					map[string]any{"matchNum": payload.MatchNumber},
				)
				if err != nil {
					continue
				}

				match.Set("score_home", payload.ScoreHome)
				match.Set("score_away", payload.ScoreAway)
				match.Set("status", payload.Status)

				if err := txApp.Save(match); err != nil {
					return fmt.Errorf("failed updating match #%d: %w", payload.MatchNumber, err)
				}
				synced++
			}
			return nil
		})

		if err != nil {
			return apis.NewBadRequestError("Sync failed", err)
		}

		return e.JSON(http.StatusOK, map[string]any{
			"success": true,
			"message": fmt.Sprintf("Successfully synced %d matches", synced),
		})
	})
}

func isAdminUser(e *core.RequestEvent) bool {
	if e.Auth == nil {
		return false
	}
	return e.Auth.GetBool("is_admin")
}

func getWC26Token(app core.App) (string, error) {
	tokenMu.RLock()
	if cachedToken != "" {
		t := cachedToken
		tokenMu.RUnlock()
		return t, nil
	}
	tokenMu.RUnlock()

	settings, err := app.FindRecordById("settings_id", wc26SettingsID)
	if err != nil {
		return "", fmt.Errorf("settings record not found: %w", err)
	}

	token := settings.GetString(wc26SettingsToken)
	if token != "" {
		tokenMu.Lock()
		cachedToken = token
		tokenMu.Unlock()
		return token, nil
	}

	return registerAndStoreWC26(app, settings)
}

func refreshWC26Token(app core.App) (string, error) {
	settings, err := app.FindRecordById("settings_id", wc26SettingsID)
	if err != nil {
		return "", fmt.Errorf("settings record not found: %w", err)
	}

	email := settings.GetString(wc26SettingsEmail)
	pass := settings.GetString(wc26SettingsPass)

	if email == "" || pass == "" {
		return registerAndStoreWC26(app, settings)
	}

	return authenticateAndStoreWC26(app, settings, email, pass)
}

func registerAndStoreWC26(app core.App, settings *core.Record) (string, error) {
	randBytes := make([]byte, 8)
	_, _ = rand.Read(randBytes)
	randStr := hex.EncodeToString(randBytes)

	email := fmt.Sprintf("wc2026-prode-%s@autogenerated.local", randStr[:8])
	pass := fmt.Sprintf("AutoPass_%s!", randStr)

	regPayload := map[string]string{
		"name":     "Prode WC2026 Bot",
		"email":    email,
		"password": pass,
	}
	regBody, _ := json.Marshal(regPayload)

	resp, err := http.Post(wc26BaseURL+"/auth/register", "application/json", bytes.NewReader(regBody))
	if err != nil {
		return "", fmt.Errorf("registration request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var authResp wc26AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		return "", fmt.Errorf("failed to parse registration response: %w", err)
	}

	if authResp.Token == "" {
		return authenticateAndStoreWC26(app, settings, email, pass)
	}

	settings.Set(wc26SettingsEmail, email)
	settings.Set(wc26SettingsPass, pass)
	settings.Set(wc26SettingsToken, authResp.Token)
	_ = app.Save(settings)

	tokenMu.Lock()
	cachedToken = authResp.Token
	tokenMu.Unlock()

	return authResp.Token, nil
}

func authenticateAndStoreWC26(app core.App, settings *core.Record, email, pass string) (string, error) {
	authPayload := map[string]string{
		"email":    email,
		"password": pass,
	}
	authBody, _ := json.Marshal(authPayload)

	resp, err := http.Post(wc26BaseURL+"/auth/authenticate", "application/json", bytes.NewReader(authBody))
	if err != nil {
		return "", fmt.Errorf("authentication request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var authResp wc26AuthResponse
	if err := json.Unmarshal(body, &authResp); err != nil {
		return "", fmt.Errorf("failed to parse auth response: %w", err)
	}

	if authResp.Token == "" {
		return "", fmt.Errorf("external API returned empty token")
	}

	settings.Set(wc26SettingsToken, authResp.Token)
	_ = app.Save(settings)

	tokenMu.Lock()
	cachedToken = authResp.Token
	tokenMu.Unlock()

	return authResp.Token, nil
}
