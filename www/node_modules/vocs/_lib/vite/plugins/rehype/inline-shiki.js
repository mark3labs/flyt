import { bundledLanguages, getSingletonHighlighter } from 'shiki';
import { visit } from 'unist-util-visit';
const inlineShikiRegex = /(.*){:(.*)}$/;
let promise;
export const rehypeInlineShiki = (options = {}) => {
    const themeNames = ('themes' in options ? Object.values(options.themes) : [options.theme]).filter(Boolean);
    const langs = options.langs || Object.keys(bundledLanguages);
    return async (tree) => {
        if (!promise)
            promise = getSingletonHighlighter({
                themes: themeNames,
                langs,
            });
        const highlighter = await promise;
        return visit(tree, 'element', (node, index, parent) => {
            if (node.tagName !== 'code')
                return;
            const match = node.children[0]?.value?.match(inlineShikiRegex);
            if (!match)
                return;
            const [, code, lang] = match;
            const hast = highlighter.codeToHast(code, { ...options, lang });
            const inlineCode = hast.children[0].children[0];
            if (!inlineCode)
                return;
            parent?.children.splice(index ?? 0, 1, inlineCode);
        });
    };
};
//# sourceMappingURL=inline-shiki.js.map