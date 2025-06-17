declare global {
    interface Window {
        // internal
        __TIRAMISU_INTERNAL_invoke: (name: string, ...args: any[]) => Promise<any>;

        // filesystem
        __TIRAMISU_FILESYSTEM_readFile: (path: string) => Promise<string>;
        __TIRAMISU_FILESYSTEM_readDir: (path: string) => Promise<string[]>;
        __TIRAMISU_INTERNAL_exists: (path: string) => Promise<boolean>;

        // notifications
        __TIRAMISU_NOTIFICATIONS_notify: (message: string, ico?: string) => Promise<void>;

        tiramisu: {
            invoke: (name: string, ...args: any[]) => Promise<any>;
            fs: {
                readFile(path: string): Promise<string>;
                readDir(path: string): Promise<string[]>;
                exists(path: string): Promise<boolean>;
            },
            notifications: {
                notify(message: string, ico?: string): Promise<void>;
            }
        };
    }
}

export { };

