import { jsx as _jsx } from "react/jsx-runtime";
import { clsx } from 'clsx';
import * as styles from './Figcaption.css.js';
export function Figcaption(props) {
    const className = clsx(props.className, styles.root);
    return _jsx("figcaption", { ...props, className: className });
}
//# sourceMappingURL=Figcaption.js.map