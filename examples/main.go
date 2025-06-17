package main

import (
	"fmt"

	_ "embed"

	t "git.iwakura.rip/grng/tiramisu"
	webview "github.com/webview/webview_go"
)

//go:embed index.html
var html string

func main() {
	app := t.New(t.TiramisuOptions{
		Debug:  true,
		Width:  800,
		Height: 600,
		Title:  "Tiramisu",
		Hints:  webview.HintFixed,
	})

	app.Run(func() {
		app.Bind("hello", func(args ...any) (any, error) {
			if len(args) == 0 {
				return "Hello, World!", nil
			}
			if len(args) == 1 {
				return fmt.Sprintf("Hello, %s!", args[0]), nil
			}
			return "Hello, unknown!", nil
		})

		app.HTML(html)
	})
}
