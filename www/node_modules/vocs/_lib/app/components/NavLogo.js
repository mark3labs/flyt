import { jsx as _jsx } from "react/jsx-runtime";
import { useConfig } from '../hooks/useConfig.js';
import { Logo } from './Logo.js';
import * as styles from './NavLogo.css.js';
export function NavLogo() {
    const config = useConfig();
    if (config.logoUrl)
        return _jsx(Logo, { className: styles.logoImage });
    return _jsx("div", { className: styles.title, children: config.title });
}
//# sourceMappingURL=NavLogo.js.map