import { jsx as _jsx } from "react/jsx-runtime";
import { clsx } from 'clsx';
import * as styles from './H2.css.js';
import { Heading } from './Heading.js';
export function H2(props) {
    return _jsx(Heading, { ...props, className: clsx(props.className, styles.root), level: 2 });
}
//# sourceMappingURL=H2.js.map