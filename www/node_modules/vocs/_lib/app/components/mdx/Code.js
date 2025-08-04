import { jsx as _jsx } from "react/jsx-runtime";
import { clsx } from 'clsx';
import * as styles from './Code.css.js';
export function Code(props) {
    const children = filterEmptyLines(props.children);
    return (_jsx("code", { ...props, className: clsx(props.className, styles.root), children: children }));
}
function filterEmptyLines(nodes) {
    if (!Array.isArray(nodes))
        return nodes;
    return nodes
        .map((child, index) => child.props &&
        'data-line' in child.props &&
        typeof child.props.children === 'string' &&
        child.props.children.trim() === '' &&
        nodes[index + 1]?.props?.className?.includes('twoslash-tag-line')
        ? null
        : child)
        .filter(Boolean);
}
//# sourceMappingURL=Code.js.map