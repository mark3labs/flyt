import { jsx as _jsx } from "react/jsx-runtime";
import { clsx } from 'clsx';
import * as styles from './Header.css.js';
export function Header(props) {
    return _jsx("header", { ...props, className: clsx(props.className, styles.root) });
}
//# sourceMappingURL=Header.js.map