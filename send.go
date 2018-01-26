package weeb

import (
	"bytes"
	"encoding/json"
)

// SendText sends the given text back as response (with given status code)
func (ctx *Context) SendText(code int, text string) error {
	ctx.Response.Header().Set("Content-Type", "text/plain")
	ctx.Response.WriteHeader(code)
	ctx.Response.Write([]byte(text))
	return nil
}

// SendJSON sends the given entity back as json (with given status code)
func (ctx *Context) SendJSON(code int, value interface{}) error {
	text, err := json.Marshal(value)
	if err != nil {
		message := "error encoding response as json"
		ctx.Log.Error(message, L{"value": value, "err": err.Error()})
		ctx.SendText(500, message)
		return nil
	}

	ctx.Response.Header().Set("Content-Type", "application/json")
	ctx.Response.WriteHeader(code)
	ctx.Response.Write(text)
	return nil
}

// SendHTML sends a rendered template back as html
func (ctx *Context) SendHTML(code int, template string, value interface{}) error {
	var b bytes.Buffer
	err := ctx.app.templates.ExecuteTemplate(&b, template, value)
	if err != nil {
		message := "error executing template"
		ctx.Log.Error(message, L{"template": template, "value": value, "err": err.Error()})
		ctx.SendText(500, message)
		return nil
	}

	ctx.Response.Header().Set("Content-Type", "text/html")
	ctx.Response.WriteHeader(code)
	ctx.Response.Write(b.Bytes())
	return nil
}

// Bind parses the request body into a given entity
func (ctx *Context) Bind(entity interface{}) error {
	defer ctx.Request.Body.Close()
	return json.NewDecoder(ctx.Request.Body).Decode(entity)
}
