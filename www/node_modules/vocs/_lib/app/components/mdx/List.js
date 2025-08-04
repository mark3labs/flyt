import { jsx as _jsx } from "react/jsx-runtime";
import { clsx } from 'clsx';
import * as styles from './List.css.js';
export function List({ ordered, ...props }) {
    const Element = ordered ? 'ol' : 'ul';
    return (_jsx(Element, { ...props, className: clsx(props.className, styles.root, ordered ? styles.ordered : styles.unordered) }));
}
//# sourceMappingURL=List.js.map