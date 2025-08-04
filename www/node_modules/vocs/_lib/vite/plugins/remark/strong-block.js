/// <reference types="mdast-util-to-hast" />
/// <reference types="mdast-util-directive" />
import { visit } from 'unist-util-visit';
export function remarkStrongBlock() {
    return (tree) => {
        visit(tree, 'strong', (node, _, parent) => {
            if (!parent)
                return;
            if (parent.type !== 'paragraph')
                return;
            if (parent.children.length > 1)
                return;
            parent.type = 'strong';
            parent.children = node.children;
        });
    };
}
//# sourceMappingURL=strong-block.js.map