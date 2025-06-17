package tiramisu

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"

	toast "github.com/electricbubble/go-toast"
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
	o     TiramisuOptions
	w     wv.WebView
	funcs map[string]FuncHandler
}

func New(o TiramisuOptions) *Tiramisu {
	w := wv.New(o.Debug)
	w.SetSize(o.Width, o.Height, o.Hints)
	w.SetTitle(o.Title)

	t := &Tiramisu{
		o:     o,
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
	// internal
	t.bind("__TIRAMISU_INTERNAL_invoke", func(args ...any) (any, error) {
		name, err := ArgAs[string](args, 0)
		if err != nil {
			return nil, err
		}
		// pass through any extra args
		if len(args) == 1 {
			return t.invoke(name)
		}
		return t.invoke(name, args[1:]...)
	})

	// filesystem
	t.bind("__TIRAMISU_FILESYSTEM_readFile", func(args ...any) (any, error) {
		filename, err := ArgAs[string](args, 0)
		if err != nil {
			return nil, err
		}
		data, err := os.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("error reading file %s: %w", filename, err)
		}
		return string(data), nil
	})

	t.bind("__TIRAMISU_FILESYSTEM_readDir", func(args ...any) (any, error) {
		dirname, err := ArgAs[string](args, 0)
		if err != nil {
			return nil, err
		}
		entries, err := os.ReadDir(dirname)
		if err != nil {
			return nil, fmt.Errorf("error reading directory %s: %w", dirname, err)
		}
		var names []string
		for _, e := range entries {
			names = append(names, e.Name())
		}
		return names, nil
	})

	t.bind("__TIRAMISU_INTERNAL_exists", func(args ...any) (any, error) {
		path, err := ArgAs[string](args, 0)
		if err != nil {
			return nil, err
		}
		_, err = os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				return false, nil
			}
			return nil, fmt.Errorf("error checking existence of %s: %w", path, err)
		}
		return true, nil
	})

	// notifications
	t.bind("__TIRAMISU_NOTIFICATIONS_notify", func(args ...any) (any, error) {
		msg, err := ArgAs[string](args, 0)
		if err != nil {
			return nil, err
		}
		if err := toast.Push(msg); err != nil {
			return nil, fmt.Errorf("error sending notification: %w", err)
		}

		return nil, nil
	})
}

func Arg(args []any, i int) (any, error) {
	if i < 0 || i >= len(args) {
		return nil, fmt.Errorf("arg %d out of range [0..%d]", i, len(args)-1)
	}
	return args[i], nil
}

func ArgAs[T any](args []any, i int) (T, error) {
	var t T

	arg, err := Arg(args, i)
	if err != nil {
		return t, err
	}

	if v, ok := arg.(T); ok {
		return v, nil
	}

	j, err := json.Marshal(arg)
	if err != nil {
		return t, fmt.Errorf("cannot marshal arg[%d]: %w", i, err)
	}
	if err := json.Unmarshal(j, &t); err != nil {
		return t, fmt.Errorf("cannot unmarshal arg[%d] -> %T: %w", i, t, err)
	}
	return t, nil
}
