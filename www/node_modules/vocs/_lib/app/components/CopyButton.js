import { jsx as _jsx } from "react/jsx-runtime";
import * as styles from './CopyButton.css.js';
import { Icon } from './Icon.js';
import { Checkmark } from './icons/Checkmark.js';
import { Copy } from './icons/Copy.js';
export function CopyButton({ copy, copied }) {
    return (_jsx("button", { className: styles.root, "data-copied": copied, onClick: copy, type: "button", children: copied ? (_jsx(Icon, { label: "Copied", size: "14px", className: styles.copied, icon: Checkmark })) : (_jsx(Icon, { label: "Copy", size: "18px", icon: Copy })) }));
}
//# sourceMappingURL=CopyButton.js.map