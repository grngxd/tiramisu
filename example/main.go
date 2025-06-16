package main

import (
	"fmt"

	t "git.iwakura.rip/grng/tiramisu"
	webview "github.com/webview/webview_go"
)

func main() {
	app := t.New(t.TiramisuOptions{
		Debug:  true,
		Width:  1200,
		Height: 800,
		Title:  "Tiramisu Example",
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

		app.HTML(`
		<!DOCTYPE html>
		<html>
			<body>
				<h1>Tiramisu Example</h1>
				<p>Click the button to see a greeting:</p>
				<button onclick="greet()">Greet</button>
				
				<script>
				function greet() {
					tiramisu.invoke("hello", "world")
						.then(response => alert(response))
						.catch(console.error)
				}
				</script>
			</body>
		</html>
		`)
	})
}
