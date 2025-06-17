const tiramisu = {
    invoke: window.__TIRAMISU_INTERNAL_invoke,
    fs: {
        readFile: window.__TIRAMISU_FILESYSTEM_readFile,
        readDir: window.__TIRAMISU_FILESYSTEM_readDir,
        exists: window.__TIRAMISU_INTERNAL_exists,
    },
    notifications: {
        notify: window.__TIRAMISU_NOTIFICATIONS_notify,
    },
}

window.tiramisu = tiramisu