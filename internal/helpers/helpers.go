package helpers

import (
	"booking/internal/config"
	"fmt"
	"net/http"
	"runtime/debug"
)

var app *config.AppConfig

func NewHelpers(a *config.AppConfig) {
	app = a
}
func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Printf("Client error with status of ", status)
	http.Error(w, http.StatusText(status), status)
}
func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Printf(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
func IsAuthenticated(r *http.Request) bool {
	exist := app.Session.Exists(r.Context(), "user_id")
	return exist
}
