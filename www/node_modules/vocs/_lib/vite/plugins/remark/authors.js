/// <reference types="mdast-util-to-hast" />
/// <reference types="mdast-util-directive" />
import { visit } from 'unist-util-visit';
export function remarkAuthors() {
    return (tree) => {
        visit(tree, (node, index, parent) => {
            if (node.type !== 'leafDirective')
                return;
            if (node.name !== 'authors')
                return;
            if (!index)
                return;
            (parent?.children[index - 1]).children.push({
                type: 'paragraph',
                data: {
                    hName: 'div',
                    hProperties: { 'data-authors': true },
                },
            });
        });
    };
}
//# sourceMappingURL=authors.js.map