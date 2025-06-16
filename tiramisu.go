package tiramisu

import (
	"embed"
	"fmt"

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
	t := &Tiramisu{
		w:     w,
		funcs: make(map[string]FuncHandler),
	}

	w.SetSize(o.Width, o.Height, o.Hints)
	w.SetTitle(o.Title)

	return t
}

func (t *Tiramisu) Run(fn func()) {
	defer t.w.Destroy()
	t.w.Dispatch(func() {
		t.injectJS()

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
	t.injectJS()
}

//go:embed runtime/out/*
var runtimeFS embed.FS

func (t *Tiramisu) injectJS() {
	js, err := runtimeFS.ReadFile("runtime/out/preload.js")
	if err != nil {
		panic(fmt.Sprintf("failed to read preload.js: %v", err))
	}

	t.w.Eval(string(js))
	t.bind("invoke", func(args ...any) (any, error) {
		name, ok := args[0].(string)
		if !ok {
			return nil, fmt.Errorf("first argument must be a string, got %T", args[0])
		}
		if len(args) == 1 {
			return t.invoke(name)
		}
		return t.invoke(name, args[1:]...)
	})
}
