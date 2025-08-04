import type { MiddlewareHandler } from 'hono';
export type ServeStaticOptions = {
    /**
     * Root path, relative to current working directory. (absolute paths are not supported)
     */
    root?: string;
    path?: string;
    index?: string;
    rewriteRequestPath?: (path: string) => string;
};
export declare const serveStatic: (options?: ServeStaticOptions) => MiddlewareHandler;
//# sourceMappingURL=serveStatic.d.js.map