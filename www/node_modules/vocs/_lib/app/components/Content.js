import { jsx as _jsx } from "react/jsx-runtime";
import { clsx } from 'clsx';
import * as styles from './Content.css.js';
export function Content({ children, className, }) {
    return _jsx("article", { className: clsx(className, styles.root), children: children });
}
//# sourceMappingURL=Content.js.map