import { jsx as _jsx } from "react/jsx-runtime";
import * as Tabs from '@radix-ui/react-tabs';
import clsx from 'clsx';
import * as styles from './Tabs.css.js';
export function Root(props) {
    return _jsx(Tabs.Root, { ...props, className: clsx(props.className, styles.root) });
}
export function List(props) {
    return _jsx(Tabs.List, { ...props, className: clsx(props.className, styles.list) });
}
export function Trigger(props) {
    return _jsx(Tabs.Trigger, { ...props, className: clsx(props.className, styles.trigger) });
}
export function Content(props) {
    return _jsx(Tabs.Content, { ...props, className: clsx(props.className, styles.content) });
}
//# sourceMappingURL=Tabs.js.map