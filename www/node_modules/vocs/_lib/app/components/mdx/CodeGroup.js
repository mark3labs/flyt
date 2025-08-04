import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import * as Tabs from '../Tabs.js';
import * as styles from './CodeGroup.css.js';
export function CodeGroup({ children }) {
    if (!Array.isArray(children))
        return null;
    const tabs = children.map((child_) => {
        const child = child_.props['data-title'] ? child_ : child_.props.children;
        const { props } = child;
        const title = props['data-title'];
        const content = props.children;
        return { title, content };
    });
    return (_jsxs(Tabs.Root, { className: styles.root, defaultValue: tabs[0].title, children: [_jsx(Tabs.List, { "aria-label": "Code group", children: tabs.map(({ title }, i) => (_jsx(Tabs.Trigger, { value: title || i.toString(), children: title }, title || i.toString()))) }), tabs.map(({ title, content }, i) => {
                const isShiki = content.props?.children?.props?.className?.includes('shiki');
                return (_jsx(Tabs.Content, { "data-shiki": isShiki, value: title || i.toString(), children: content }, title || i.toString()));
            })] }));
}
//# sourceMappingURL=CodeGroup.js.map