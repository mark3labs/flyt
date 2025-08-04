import { jsx as _jsx } from "react/jsx-runtime";
import { clsx } from 'clsx';
import * as styles from './TableRow.css.js';
export function TableRow(props) {
    return _jsx("tr", { ...props, className: clsx(props.className, styles.root) });
}
//# sourceMappingURL=TableRow.js.map