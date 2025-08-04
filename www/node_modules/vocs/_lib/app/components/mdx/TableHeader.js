import { jsx as _jsx } from "react/jsx-runtime";
import { clsx } from 'clsx';
import * as styles from './TableHeader.css.js';
export function TableHeader(props) {
    return _jsx("th", { ...props, className: clsx(props.className, styles.root) });
}
//# sourceMappingURL=TableHeader.js.map