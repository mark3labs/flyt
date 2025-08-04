/// <reference types="mdast-util-to-hast" />
/// <reference types="mdast-util-directive" />
import { visit } from 'unist-util-visit';
export function remarkBlogPosts() {
    return (tree) => {
        visit(tree, (node) => {
            if (node.type !== 'leafDirective')
                return;
            if (node.name !== 'blog-posts')
                return;
            const data = node.data || (node.data = {});
            data.hName = 'div';
            data.hProperties = { 'data-blog-posts': true };
        });
    };
}
//# sourceMappingURL=blog-posts.js.map