declare global {
    interface Window {
        invoke: (name: string, ...args: any[]) => Promise<any>;
        tiramisu: {
            invoke: (name: string, ...args: any[]) => Promise<any>;
        };
    }
}

export { };

