package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		backBundle := ServerBundle()
		props := CreateProps(map[string]interface{}{
			"Name":          "Go React SSR",
			"InitialNumber": 99,
		})
		html := RenderHTML(backBundle, props)
		page := SSRPage{
			RenderedContent: template.HTML(html),
			Props:           template.JS(props),
		}
		err := page.Render(w)
		if err != nil {
			fmt.Println(err)
			w.Write([]byte(err.Error()))
		}
	})
	fmt.Println("Server is running at http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
