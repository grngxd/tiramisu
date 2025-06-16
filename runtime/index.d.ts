declare global {
    interface Window {
        __TIRAMISU_INTERNAL_invoke: (name: string, ...args: any[]) => Promise<any>;
        __TIRAMISU_INTERNAL_readFile: (path: string) => Promise<string>;
        __TIRAMISU_INTERNAL_readDir: (path: string) => Promise<string[]>;

        tiramisu: {
            invoke: (name: string, ...args: any[]) => Promise<any>;
            fs: {
                readFile(path: string): Promise<string>;
                readDir(path: string): Promise<string[]>;
            }
        };
    }
}

export { };

