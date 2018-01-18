package weeb

import (
	"encoding/json"
	"net/http"
)

// SendText sends the given text back as response (with given status code)
func (app *App) SendText(w http.ResponseWriter, code int, text string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
	w.Write([]byte(text))
}

// SendJSON sends the given entity back as json (with given status code)
func (app *App) SendJSON(w http.ResponseWriter, code int, value interface{}) {
	text, err := json.Marshal(value)
	if err != nil {
		message := "error encoding response as json"
		app.Log.Error(err.Error(), L{"value": value})
		app.SendText(w, 500, message)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(text)
}

// Bind parses the request body into a given entity
func (app *App) Bind(r *http.Request, entity interface{}) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(entity)
}
