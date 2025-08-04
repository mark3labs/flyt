import * as vite from 'vite';
export type BuildParameters = {
    clean?: boolean;
    logger?: vite.Logger;
    hooks?: {
        onBundleStart?: () => void;
        onBundleEnd?: ({ error }: {
            error?: Error;
        }) => void;
        onPrerenderStart?: () => void;
        onPrerenderEnd?: ({ error }: {
            error?: Error;
        }) => void;
        onScriptsStart?: () => void;
        onScriptsEnd?: ({ error }: {
            error?: Error;
        }) => void;
    };
    logLevel?: vite.LogLevel;
    outDir?: string;
    publicDir?: string;
    searchIndex?: boolean;
};
export declare function build({ clean, logger, hooks, logLevel, outDir, publicDir, searchIndex, }?: BuildParameters): Promise<void>;
//# sourceMappingURL=build.d.js.map