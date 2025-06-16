package tiramisu

import (
	"embed"
	"fmt"
	"os"

	wv "github.com/webview/webview_go"
)

type TiramisuOptions struct {
	Debug  bool
	Width  int
	Height int
	Title  string
	Hints  wv.Hint
}

type FuncHandler func(args ...any) (any, error)

type Tiramisu struct {
	w     wv.WebView
	funcs map[string]FuncHandler
}

func New(o TiramisuOptions) *Tiramisu {
	w := wv.New(o.Debug)
	w.SetSize(o.Width, o.Height, o.Hints)
	w.SetTitle(o.Title)

	t := &Tiramisu{
		w:     w,
		funcs: make(map[string]FuncHandler),
	}

	return t
}

func (t *Tiramisu) Run(fn func()) {
	defer t.w.Destroy()
	t.w.Dispatch(func() {
		t.loadJSRuntime()
		t.loadGoRuntime()

		if fn != nil {
			fn()
		}
	})
	t.w.Run()
}

func (t *Tiramisu) bind(name string, fn FuncHandler) {
	t.w.Bind(name, fn)
}

func (t *Tiramisu) Bind(name string, fn FuncHandler) {
	if _, exists := t.funcs[name]; exists {
		panic(fmt.Sprintf("function %s is already bound", name))
	}

	t.funcs[name] = fn
	t.bind(name, func(args ...any) (any, error) {
		return t.invoke(name, args...)
	})
}

func (t *Tiramisu) invoke(name string, args ...any) (any, error) {
	fn, ok := t.funcs[name]
	if !ok {
		return nil, fmt.Errorf("function %s not found", name)
	}
	result, err := fn(args...)
	if err != nil {
		return nil, fmt.Errorf("error invoking function %s: %w", name, err)
	}
	return result, nil
}

func (t *Tiramisu) Eval(js string) {
	t.w.Eval(js)
}

func (t *Tiramisu) Evalf(js string, args ...any) {
	fmt := fmt.Sprintf(js, args...)
	t.w.Eval(fmt)
}

func (t *Tiramisu) HTML(html string) {
	t.w.SetHtml(html)
	t.loadJSRuntime()
	t.loadGoRuntime()
}

//go:embed runtime/out/*
var runtimeFS embed.FS

func (t *Tiramisu) loadJSRuntime() {
	js, err := runtimeFS.ReadFile("runtime/out/preload.js")
	if err != nil {
		panic(fmt.Sprintf("failed to read preload.js: %v", err))
	}

	t.w.Eval(string(js))
}

func (t *Tiramisu) loadGoRuntime() {
	t.bind("__TIRAMISU_INTERNAL_invoke", func(args ...any) (any, error) {
		name, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("first argument must be a string, got %T", args[0])
		}
		if len(args) == 1 {
			return t.invoke(name)
		}
		return t.invoke(name, args[1:]...)
	})

	t.bind("__TIRAMISU_INTERNAL_readFile", func(args ...any) (any, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("readFile expects exactly one argument, got %d", len(args))
		}

		filename, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("readFile expects a string argument, got %T", args[0])
		}

		data, err := os.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("error reading file %s: %w", filename, err)
		}
		return string(data), nil
	})

	t.bind("__TIRAMISU_INTERNAL_readDir", func(args ...any) (any, error) {
		if len(args) != 1 {
			return nil, fmt.Errorf("readDir expects exactly one argument, got %d", len(args))
		}

		dirname, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("readDir expects a string argument, got %T", args[0])
		}

		files, err := os.ReadDir(dirname)
		if err != nil {
			return nil, fmt.Errorf("error reading directory %s: %w", dirname, err)
		}

		var fileNames []string
		for _, file := range files {
			fileNames = append(fileNames, file.Name())
		}
		return fileNames, nil
	})
}
