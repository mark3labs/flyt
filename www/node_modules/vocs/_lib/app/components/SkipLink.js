import { jsx as _jsx } from "react/jsx-runtime";
import clsx from 'clsx';
import { useLocation } from 'react-router';
import { visuallyHidden } from '../styles/utils.css.js';
import * as styles from './SkipLink.css.js';
export const skipLinkId = 'vocs-content';
export function SkipLink() {
    const { pathname } = useLocation();
    return (_jsx("a", { className: clsx(styles.root, visuallyHidden), href: `${pathname}#${skipLinkId}`, children: "Skip to content" }));
}
//# sourceMappingURL=SkipLink.js.map