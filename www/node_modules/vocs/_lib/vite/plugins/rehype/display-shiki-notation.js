/// <reference types="mdast-util-to-hast" />
/// <reference types="mdast-util-directive" />
import { visit } from 'unist-util-visit';
export function rehypeShikiDisplayNotation() {
    return (tree) => {
        visit(tree, 'text', (node) => {
            if (node.value.includes('// [\\!'))
                node.value = node.value.replace('// [\\!', '// [!').replaceAll('// //', '');
        });
    };
}
//# sourceMappingURL=display-shiki-notation.js.map