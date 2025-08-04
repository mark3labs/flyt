import type { BuildParameters as BuildParameters_ } from '../../vite/build.js';
export type BuildParameters = Pick<BuildParameters_, 'clean' | 'logLevel' | 'outDir' | 'publicDir' | 'searchIndex'> & {
    clearScreen?: boolean;
};
export declare function build({ clean, clearScreen, logLevel, outDir, publicDir, searchIndex, }: BuildParameters): Promise<void>;
//# sourceMappingURL=build.d.js.map