import { jsx as _jsx } from "react/jsx-runtime";
import { clsx } from 'clsx';
import { Callout } from '../Callout.js';
import * as styles from './Aside.css.js';
export function Aside(props) {
    const className = clsx(props.className, styles.root);
    if ('data-callout' in props)
        return (_jsx(Callout, { className: className, type: props['data-callout'], children: props.children }));
    return _jsx("aside", { ...props, className: className });
}
//# sourceMappingURL=Aside.js.map