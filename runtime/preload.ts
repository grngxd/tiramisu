const tiramisu = {
    invoke: window.__TIRAMISU_INTERNAL_invoke,
    fs: {
        readFile: window.__TIRAMISU_INTERNAL_readFile,
        readDir: window.__TIRAMISU_INTERNAL_readDir
    }
}

window.tiramisu = tiramisu