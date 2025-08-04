import { jsx as _jsx } from "react/jsx-runtime";
import * as Popover_ from '@radix-ui/react-popover';
import clsx from 'clsx';
import * as styles from './Popover.css.js';
Popover.Root = Popover_.Root;
Popover.Trigger = Popover_.Trigger;
export function Popover({ children, className }) {
    return (_jsx(Popover_.Portal, { children: _jsx(Popover_.Content, { className: clsx(styles.root, className), sideOffset: 12, children: children }) }));
}
//# sourceMappingURL=Popover.js.map