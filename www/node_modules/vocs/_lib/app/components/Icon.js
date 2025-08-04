import { jsx as _jsx } from "react/jsx-runtime";
import { assignInlineVars } from '@vanilla-extract/dynamic';
import clsx from 'clsx';
import * as styles from './Icon.css.js';
export function Icon({ className, label, icon: Icon, size, style }) {
    return (_jsx("div", { "aria-label": label, className: clsx(styles.root, className), role: "img", style: {
            ...style,
            ...assignInlineVars({ [styles.sizeVar]: size }),
        }, children: _jsx(Icon, { height: size, width: size }) }));
}
//# sourceMappingURL=Icon.js.map