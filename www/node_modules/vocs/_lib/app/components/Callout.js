import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { clsx } from 'clsx';
import * as styles from './Callout.css.js';
import { CheckCircle } from './icons/CheckCircle.js';
import { ExclamationTriangle } from './icons/ExclamationTriangle.js';
import { InfoCircled } from './icons/InfoCircled.js';
import { LightningBolt } from './icons/LightningBolt.js';
export function Callout({ className, children, type }) {
    return (_jsxs("aside", { className: clsx(className, styles.root, styles[type]), children: [_jsxs("div", { className: styles.icon, children: [type === 'note' && _jsx(InfoCircled, {}), type === 'info' && _jsx(InfoCircled, {}), type === 'warning' && _jsx(ExclamationTriangle, {}), type === 'danger' && _jsx(ExclamationTriangle, {}), type === 'tip' && _jsx(LightningBolt, {}), type === 'success' && _jsx(CheckCircle, {})] }), _jsx("div", { className: styles.content, children: children })] }));
}
//# sourceMappingURL=Callout.js.map