package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"

	v8 "rogchap.com/v8go"
)

type SSRPage struct {
	RenderedContent template.HTML
	ClientBundle    template.JS
}

func (ssr *SSRPage) Render(writer http.ResponseWriter) error {
	tmpl, err := template.New("page").Parse(htmlTemplate)
	if err != nil {
		log.Fatal("Error parsing template:", err)
		return errors.New("Error parsing template")
	}
	writer.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(writer, ssr)
	if err != nil {
		fmt.Println(err)
		return errors.New("error executing template")
	}
	return nil
}

const htmlTemplate = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<title>React App</title>
	</head>
	<body>
		<div id="app">{{ .RenderedContent }}</div>
		<script type="module">
		{{ .ClientBundle }}
		</script>
	</body>
	</html>
`

var iso = v8.NewIsolate()

func RenderHTML(backendBundle string, props string) string {
	ctx := v8.NewContext(iso)
	defer ctx.Close() // Close context
	_, err := ctx.RunScript(backendBundle, "bundle.js")
	if err != nil {
		log.Fatalf("Failed to evaluate bundled script: %v", err)
	}
	// Pass props to the renderApp function
	val, err := ctx.RunScript(fmt.Sprintf("renderApp(%s)", props), "render.js")
	if err != nil {
		log.Fatalf("Failed to render React component: %v", err)
	}

	return val.String()
}
