/// <reference types="mdast-util-to-hast" />
/// <reference types="mdast-util-directive" />
import { visit } from 'unist-util-visit';
export function remarkSponsors() {
    return (tree) => {
        visit(tree, (node, index, parent) => {
            if (node.type !== 'leafDirective')
                return;
            if (node.name !== 'sponsors')
                return;
            if (!index)
                return;
            (parent?.children[index]).children.push({
                type: 'paragraph',
                data: {
                    hName: 'div',
                    hProperties: { 'data-sponsors': true },
                },
            });
        });
    };
}
//# sourceMappingURL=sponsors.js.map