import { jsx as _jsx } from "react/jsx-runtime";
import { clsx } from 'clsx';
import * as styles from './ListItem.css.js';
export function ListItem(props) {
    return _jsx("li", { ...props, className: clsx(props.className, styles.root) });
}
//# sourceMappingURL=ListItem.js.map