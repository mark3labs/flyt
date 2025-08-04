import { jsx as _jsx } from "react/jsx-runtime";
import * as NavigationMenu from '@radix-ui/react-navigation-menu';
import { assignInlineVars } from '@vanilla-extract/dynamic';
import clsx from 'clsx';
import { useConfig } from '../hooks/useConfig.js';
import { Link as Link_ } from './Link.js';
import * as styles from './NavigationMenu.css.js';
export const Root = (props) => (_jsx(NavigationMenu.Root, { ...props, className: clsx(props.className, styles.root) }));
export const List = (props) => (_jsx(NavigationMenu.List, { ...props, className: clsx(props.className, styles.list) }));
export const Link = ({ active, children, className, href, }) => (_jsx(NavigationMenu.Link, { asChild: true, children: _jsx(Link_, { "data-active": active, className: clsx(className, styles.link), href: href, variant: "styleless", children: children }) }));
export const Item = (props) => (_jsx(NavigationMenu.Item, { ...props, className: clsx(props.className, styles.item) }));
export const Trigger = ({ active, className, ...props }) => {
    const { basePath } = useConfig();
    const assetBasePath = import.meta.env.PROD ? basePath : '';
    return (_jsx(NavigationMenu.Trigger, { ...props, "data-active": active, className: clsx(className, styles.trigger), style: assignInlineVars({
            [styles.chevronDownIcon]: `url(${assetBasePath}/.vocs/icons/chevron-down.svg)`,
        }) }));
};
export const Content = (props) => (_jsx(NavigationMenu.Content, { ...props, className: clsx(props.className, styles.content) }));
//# sourceMappingURL=NavigationMenu.js.map