import { jsx as _jsx, Fragment as _Fragment, jsxs as _jsxs } from "react/jsx-runtime";
import clsx from 'clsx';
import { useConfig } from '../hooks/useConfig.js';
import * as styles from './Logo.css.js';
export function Logo({ className }) {
    const { logoUrl } = useConfig();
    if (!logoUrl)
        return null;
    return (_jsx(_Fragment, { children: typeof logoUrl === 'string' ? (_jsx("img", { alt: "Logo", className: clsx(className, styles.root), src: logoUrl })) : (_jsxs(_Fragment, { children: [_jsx("img", { alt: "Logo", className: clsx(className, styles.root, styles.logoDark), src: logoUrl.dark }), _jsx("img", { alt: "Logo", className: clsx(className, styles.root, styles.logoLight), src: logoUrl.light })] })) }));
}
//# sourceMappingURL=Logo.js.map