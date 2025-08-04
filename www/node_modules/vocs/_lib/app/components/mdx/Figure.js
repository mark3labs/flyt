import { jsx as _jsx } from "react/jsx-runtime";
import { clsx } from 'clsx';
import * as styles from './Figure.css.js';
export function Figure(props) {
    const className = clsx(props.className, styles.root);
    return _jsx("figure", { ...props, className: className });
}
//# sourceMappingURL=Figure.js.map