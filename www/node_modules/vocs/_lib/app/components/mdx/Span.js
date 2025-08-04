import { jsx as _jsx } from "react/jsx-runtime";
import { clsx } from 'clsx';
import * as styles from './Span.css.js';
import { TwoslashPopover } from './TwoslashPopover.js';
export function Span(props) {
    const className = clsx(props.className, styles.root);
    if (props.className?.includes('twoslash-hover'))
        return _jsx(TwoslashPopover, { ...props, className: className });
    return _jsx("span", { ...props, className: clsx(props.className, styles.root) });
}
//# sourceMappingURL=Span.js.map