import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import clsx from 'clsx';
import { useConfig } from '../hooks/useConfig.js';
import { Button as Button_ } from './Button.js';
import * as styles from './HomePage.css.js';
import { Logo as Logo_ } from './Logo.js';
import * as Tabs from './Tabs.js';
export function Root({ children, className }) {
    return _jsx("div", { className: clsx(className, styles.root), children: children });
}
export function Logo({ className }) {
    const { logoUrl, title } = useConfig();
    return logoUrl ? (_jsx("div", { className: clsx(className, styles.logo), children: _jsx(Logo_, {}) })) : (_jsx("h1", { className: clsx(className, styles.title), children: title }));
}
export function Tagline({ children, className }) {
    return _jsx("div", { className: clsx(className, styles.tagline), children: children });
}
export function Description({ children, className }) {
    return _jsx("div", { className: clsx(className, styles.description), children: children });
}
export function Buttons({ children, className }) {
    return _jsx("div", { className: clsx(className, styles.buttons), children: children });
}
export function Button(props) {
    return _jsx(Button_, { ...props, className: clsx(styles.button, props.className) });
}
export function InstallPackage({ name, type = 'install', }) {
    return (_jsxs(Tabs.Root, { className: styles.tabs, defaultValue: "npm", children: [_jsxs(Tabs.List, { className: styles.tabsList, children: [_jsx(Tabs.Trigger, { value: "npm", children: "npm" }), _jsx(Tabs.Trigger, { value: "pnpm", children: "pnpm" }), _jsx(Tabs.Trigger, { value: "yarn", children: "yarn" })] }), _jsxs(Tabs.Content, { className: styles.tabsContent, value: "npm", children: [_jsx("span", { className: styles.packageManager, children: "npm" }), " ", type === 'init' ? 'init' : 'install', ' ', name] }), _jsxs(Tabs.Content, { className: styles.tabsContent, value: "pnpm", children: [_jsx("span", { className: styles.packageManager, children: "pnpm" }), " ", type === 'init' ? 'create' : 'add', ' ', name] }), _jsxs(Tabs.Content, { className: styles.tabsContent, value: "yarn", children: [_jsx("span", { className: styles.packageManager, children: "yarn" }), " ", type === 'init' ? 'create' : 'add', ' ', name] })] }));
}
//# sourceMappingURL=HomePage.js.map