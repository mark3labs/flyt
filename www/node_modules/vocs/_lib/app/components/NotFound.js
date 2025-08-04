import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { spaceVars } from '../styles/vars.css.js';
import { Link } from './Link.js';
import * as styles from './NotFound.css.js';
import { H1 } from './mdx/H1.js';
import { Paragraph } from './mdx/Paragraph.js';
export function NotFound() {
    return (_jsxs("div", { className: styles.root, children: [_jsx(H1, { children: "Page Not Found" }), _jsx("div", { style: { height: spaceVars['24'] } }), _jsx("hr", { className: styles.divider }), _jsx("div", { style: { height: spaceVars['24'] } }), _jsx(Paragraph, { children: "The page you were looking for could not be found." }), _jsx("div", { style: { height: spaceVars['8'] } }), _jsx(Link, { href: "/", children: "Go to Home Page" })] }));
}
//# sourceMappingURL=NotFound.js.map