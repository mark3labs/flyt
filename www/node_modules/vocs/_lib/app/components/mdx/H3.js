import { jsx as _jsx } from "react/jsx-runtime";
import { clsx } from 'clsx';
import * as styles from './H3.css.js';
import { Heading } from './Heading.js';
export function H3(props) {
    return _jsx(Heading, { ...props, className: clsx(props.className, styles.root), level: 3 });
}
//# sourceMappingURL=H3.js.map