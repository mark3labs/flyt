import { jsx as _jsx } from "react/jsx-runtime";
import clsx from 'clsx';
import * as styles from './Button.css.js';
import { Link } from './Link.js';
export function Button({ children, className, href, variant }) {
    return (_jsx(Link, { className: clsx(className, styles.button, variant === 'accent' && styles.button_accent), href: href, variant: "styleless", children: children }));
}
//# sourceMappingURL=Button.js.map