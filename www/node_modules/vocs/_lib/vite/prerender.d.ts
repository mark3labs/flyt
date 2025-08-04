import type { Logger } from 'vite';
type PrerenderParameters = {
    logger?: Logger;
    outDir?: string;
};
export declare function prerender({ logger, outDir }: PrerenderParameters): Promise<void>;
export {};
//# sourceMappingURL=prerender.d.js.map