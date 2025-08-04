import type { PluggableList } from 'unified';
import type { PluginOption } from 'vite';
import type { ParsedConfig } from '../../config.js';
type RemarkPluginsParameters = {
    markdown?: ParsedConfig['markdown'];
};
export declare const getRemarkPlugins: ({ markdown }?: RemarkPluginsParameters) => PluggableList;
type RehypePluginsParameters = {
    markdown?: ParsedConfig['markdown'];
    cacheDir?: ParsedConfig['cacheDir'];
    rootDir?: ParsedConfig['rootDir'];
    twoslash?: ParsedConfig['twoslash'] | false;
};
export declare const getRehypePlugins: ({ markdown, cacheDir, rootDir, twoslash, }?: RehypePluginsParameters) => PluggableList;
export declare function mdx(): Promise<PluginOption[]>;
export {};
//# sourceMappingURL=mdx.d.js.map