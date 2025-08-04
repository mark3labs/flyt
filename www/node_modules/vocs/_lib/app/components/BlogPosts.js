import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { posts } from 'virtual:blog';
import { Fragment } from 'react';
import { Authors } from './Authors.js';
import * as styles from './BlogPosts.css.js';
import { RouterLink } from './RouterLink.js';
export function BlogPosts() {
    return (_jsx("div", { className: styles.root, children: posts.map((post, index) => (_jsxs(Fragment, { children: [_jsx("div", { className: styles.post, children: _jsxs(RouterLink, { to: post.path, children: [_jsx("h2", { className: styles.title, children: post.title }), _jsx(Authors, { authors: post.authors, date: post.date }), _jsxs("p", { className: styles.description, children: [post.description, " ", _jsx("span", { className: styles.readMore, children: "[\u2192]" })] })] }) }), index < posts.length - 1 && _jsx("hr", { className: styles.divider })] }, index))) }));
}
//# sourceMappingURL=BlogPosts.js.map