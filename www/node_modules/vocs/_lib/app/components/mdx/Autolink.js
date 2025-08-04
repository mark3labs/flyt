import { jsx as _jsx } from "react/jsx-runtime";
import { clsx } from 'clsx';
import { Link } from 'react-router';
import * as styles from './Autolink.css.js';
export function Autolink(props) {
    if (!props.href)
        return null;
    return _jsx(Link, { ...props, className: clsx(props.className, styles.root), to: props.href });
}
//# sourceMappingURL=Autolink.js.map