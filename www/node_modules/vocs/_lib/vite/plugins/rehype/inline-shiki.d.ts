import type { RehypeShikiCoreOptions } from '@shikijs/rehype/core';
import type { Root } from 'hast';
import type { BuiltinLanguage } from 'shiki';
import type { LanguageInput } from 'shiki/core';
import type { Plugin } from 'unified';
export type RehypeInlineShikiOptions = RehypeShikiCoreOptions & {
    /**
     * Language names to include.
     *
     * @default Object.keys(bundledLanguages)
     */
    langs?: Array<LanguageInput | BuiltinLanguage>;
};
export declare const rehypeInlineShiki: Plugin<[RehypeInlineShikiOptions], Root>;
//# sourceMappingURL=inline-shiki.d.js.map