import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { Kbd } from './mdx/Kbd.js';
import * as styles from './KeyboardShortcut.css.js';
export function KeyboardShortcut(props) {
    const { description, keys } = props;
    return (_jsxs("span", { className: styles.root, children: [description, _jsx("span", { className: styles.kbdGroup, children: keys.map((key) => (_jsx(Kbd, { children: key }, key))) })] }));
}
//# sourceMappingURL=KeyboardShortcut.js.map