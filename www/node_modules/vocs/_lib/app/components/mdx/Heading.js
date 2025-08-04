import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { clsx } from 'clsx';
import { root, slugTarget } from './Heading.css.js';
export function Heading({ level, ...props }) {
    const Component = `h${level}`;
    return (_jsxs(Component, { ...props, id: undefined, className: clsx(props.className, root), children: [_jsx("div", { id: props.id, className: slugTarget }), props.children] }));
}
//# sourceMappingURL=Heading.js.map