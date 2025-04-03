package handlers

import (
	"matterpoll-bot/config"
	"matterpoll-bot/internal/storage"
	"net/http"
)

// TokenValidatorMiddleware проверяет полученный токен из тела запроса (только для режима "database").
func TokenValidatorMiddleware(store storage.StoreInterface, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		cmdPath := r.Form.Get("command")
		token := r.Form.Get("token")
		if cmdPath == "" || token == "" {
			http.Error(w, "'command' or 'token' are empty in the form data", http.StatusBadRequest)
			return
		}

		if config.Mode == "database" && !store.ValidateCmdToken(cmdPath, token) {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		
		next(w, r)
	}
}
