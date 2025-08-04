import { jsx as _jsx } from "react/jsx-runtime";
import { clsx } from 'clsx';
import * as styles from './H5.css.js';
import { Heading } from './Heading.js';
export function H5(props) {
    return _jsx(Heading, { ...props, className: clsx(props.className, styles.root), level: 5 });
}
//# sourceMappingURL=H5.js.map