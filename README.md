# 🍥 tiramisu
> ` is tiramisu a cake or a pie?`

Build modern, cross-platform desktop apps in HTML + Go from one codebase. It uses the built-in OS webview (WebView2 on Windows, WebKitGTK on Linux, and WebKit on macOS) to render HTML to the screen. Tiramisu creates a connection between Go and the webview, allowing you manipulate the UI and calling Go methods from the webview seamlessly.

## Features
- 💻 **Cross-platform**: Write once, run everywhere. Tiramisu supports Windows, macOS, and Linux.
- 🪶 **Lightweight**: No need for a heavy framework. Tiramisu uses the built-in webview of the OS.
- ⚡**Fast development**: Use *ANY* web framework for  your UI. Tiramisu handles all the magic of making it work, for you.

## Installation
`go install git.iwakura.rip/grng/tiramisu`

## Example

```go
package main

import (
    "fmt"
    "git.iwakura.rip/grng/tiramisu"
)

func main() {
    // create the webview instance
    app := tiramisu.New(tiramisu.Options{
        Title:  "Tiramisu Example",
        Width:  800,
        Height: 600,
    })

    // bind a go function to the webview
    app.Bind("hello", func(name string) string {
        return fmt.Sprintf("Hello, %s!", name)
    })

    // set the HTML content of the webview
    app.HTML(`
        <h1>Tiramisu Example</h1>
        <button onclick="tiramisu.invoke('hello', 'world').then(alert)">Greet</button>
    `)

    // t.Run() also allows you to pass a func(), which is executed on the main thread
    // before the webview is shown, so you can do any setup you need.
    app.Run(/* func() {
        // This code runs on the main thread before the webview is shown
        fmt.Println("Webview is ready!")
    }*/)
}
```

## Building

### Windows
Windows users should use the `-ldflags='-H windowsgui'` flag to avoid showing a console window when running the app.
```
go build -ldflags='-H windowsgui' ./example
```

### *nix
```
go build ./example
```


## Contributing & Development
```
git clone ssh://git@git.iwakura.rip:6969/grng/tiramisu.git
cd tiramisu
air
```

## License
`tiramisu` is licensed under the GNU General Public License (GPL) v3.0. See the [LICENSE](LICENSE) file for more details.