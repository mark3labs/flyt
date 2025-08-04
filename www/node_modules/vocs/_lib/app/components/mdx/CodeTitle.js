import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { clsx } from 'clsx';
import { Icon } from '../Icon.js';
import { File } from '../icons/File.js';
import { Terminal } from '../icons/Terminal.js';
import * as styles from './CodeTitle.css.js';
export function CodeTitle({ children, className, language, ...props }) {
    return (_jsxs("div", { ...props, className: clsx(className, styles.root), children: [language === 'bash' ? (_jsx(Icon, { label: "Terminal", size: "14px", icon: Terminal, style: { marginTop: 3 } })) : children.match(/\.(.*)$/) ? (_jsx(Icon, { label: "File", size: "14px", icon: File, style: { marginTop: 1 } })) : null, children] }));
}
//# sourceMappingURL=CodeTitle.js.map