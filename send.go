package weeb

import (
	"encoding/json"
)

// Text sends the given text back as response (with given status code)
func (ctx *Context) Text(code int, text string) error {
	ctx.SetHeader("Content-Type", "text/plain; charset=utf-8")
	ctx.SetStatusCode(code)
	ctx.SetBody(text)
	return nil
}

// JSON sends the given entity back as json (with given status code)
func (ctx *Context) JSON(code int, value interface{}) error {
	text, err := json.Marshal(value)
	if err != nil {
		message := "error encoding response as json"
		ctx.Log.Error(message, L{"value": value, "err": err.Error()})
		ctx.Text(500, message)
		return nil
	}

	ctx.SetHeader("Content-Type", "application/json; charset=utf-8")
	ctx.SetStatusCode(code)
	ctx.SetBody(string(text))
	return nil
}

// HTML sends a rendered template back as html
func (ctx *Context) HTML(code int, template string, value J) error {
	contents, err := ctx.Template(template, value)
	if err != nil {
		message := "error executing template"
		ctx.Log.Error(message, L{"template": template, "value": value, "err": err.Error()})
		ctx.Error(500, message)
		return nil
	}

	ctx.SetHeader("Content-Type", "text/html; charset=utf-8")
	ctx.SetStatusCode(code)
	ctx.SetBody(contents)
	return nil
}
