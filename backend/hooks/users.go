package hooks

import (
	"strings"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/security"
)

func RegisterUserHooks(app core.App) {

	app.OnRecordAuthWithOAuth2Request().BindFunc(func(e *core.RecordAuthWithOAuth2RequestEvent) error {

		// Si e.Record es nil o IsNewRecord es true, significa que es una cuenta nueva creándose por Google
		if e.Record == nil || e.IsNewRecord {
			email := e.OAuth2User.Email

			if email != "" {
				// 1. Extraer texto antes del '@' (añadiendo [0] al final de parts)
				parts := strings.Split(email, "@")
				baseUsername := parts[0] // <-- Añade [0] aquí para obtener el string, no el slice

				baseUsername = strings.ReplaceAll(baseUsername, " ", "_")
				finalUsername := baseUsername

				// 2. Comprobación segura de duplicados en la base de datos
				for {
					_, err := e.App.FindFirstRecordByFilter("users", "username = {:username}", map[string]any{
						"username": finalUsername,
					})
					if err != nil {
						// Si da error, significa que el username NO existe (está libre)
						break
					}
					// Si ya existe, concatenamos un sufijo aleatorio
					finalUsername = baseUsername + "_" + security.RandomString(4)
				}

				// 3. Inyectar el username en los datos de creación
				if e.CreateData == nil {
					e.CreateData = map[string]any{}
				}
				e.CreateData["username"] = finalUsername
			}
		}

		return e.Next()
	})
}
