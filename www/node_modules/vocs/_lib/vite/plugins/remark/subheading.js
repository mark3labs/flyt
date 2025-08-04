/// <reference types="mdast-util-to-hast" />
/// <reference types="mdast-util-directive" />
import { visit } from 'unist-util-visit';
export function remarkSubheading() {
    return (tree) => {
        visit(tree, 'heading', (node, index, parent) => {
            if (!index)
                return;
            if (node.depth !== 1)
                return;
            if (node.children.length === 0)
                return;
            const subheadingRegex = / \[(.*)\]$/;
            const subheadingChild = node.children.find((child) => 'value' in child && typeof child.value === 'string' && child.value.match(subheadingRegex));
            const [match, subheading] = subheadingChild?.value?.match(subheadingRegex) ?? [];
            if (subheadingChild)
                subheadingChild.value = subheadingChild?.value?.replace(match, '');
            // remove original heading
            parent?.children.splice(index, 1);
            const header = {
                type: 'paragraph',
                data: {
                    hName: 'header',
                },
                children: [
                    node,
                    subheading
                        ? {
                            type: 'paragraph',
                            children: [{ type: 'text', value: subheading }],
                            data: {
                                hName: 'div',
                                hProperties: {
                                    role: 'doc-subtitle',
                                },
                            },
                        }
                        : undefined,
                ].filter(Boolean),
            };
            parent?.children.splice(index, 0, header);
        });
    };
}
//# sourceMappingURL=subheading.js.map