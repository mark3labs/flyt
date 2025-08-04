import { jsx as _jsx } from "react/jsx-runtime";
import { clsx } from 'clsx';
import * as styles from './H1.css.js';
import { Heading } from './Heading.js';
export function H1(props) {
    return _jsx(Heading, { ...props, className: clsx(props.className, styles.root), level: 1 });
}
//# sourceMappingURL=H1.js.map