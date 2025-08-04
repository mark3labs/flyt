import { jsx as _jsx } from "react/jsx-runtime";
import { clsx } from 'clsx';
import { Footnotes } from './Footnotes.js';
import * as styles from './Section.css.js';
export function Section(props) {
    if ('data-footnotes' in props)
        return _jsx(Footnotes, { ...props, className: clsx(props.className, styles.root) });
    return _jsx("section", { ...props, className: clsx(props.className, styles.root) });
}
//# sourceMappingURL=Section.js.map