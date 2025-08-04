import { jsx as _jsx } from "react/jsx-runtime";
import { clsx } from 'clsx';
import * as styles from './TableCell.css.js';
export function TableCell(props) {
    return _jsx("td", { ...props, className: clsx(props.className, styles.root) });
}
//# sourceMappingURL=TableCell.js.map