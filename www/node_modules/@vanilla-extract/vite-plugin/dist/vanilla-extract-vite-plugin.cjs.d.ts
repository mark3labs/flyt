import { Plugin } from 'vite';
import { IdentifierOption } from '@vanilla-extract/integration';

type PluginFilter = (filterProps: {
    /** The name of the plugin */
    name: string;
    /**
     * The `mode` vite is running in.
     * @see https://vite.dev/guide/env-and-mode.html#modes
     */
    mode: string;
}) => boolean;
interface Options {
    identifiers?: IdentifierOption;
    unstable_pluginFilter?: PluginFilter;
    unstable_mode?: 'transform' | 'emitCss' | 'inlineCssInDev';
}
declare function vanillaExtractPlugin({ identifiers, unstable_pluginFilter: pluginFilter, unstable_mode, }?: Options): Plugin[];

export { vanillaExtractPlugin };
