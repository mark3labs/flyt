import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { TopNavEnd } from 'virtual:consumer-components';
import clsx from 'clsx';
import { useLocation } from 'react-router';
import { useActiveNavIds } from '../hooks/useActiveNavIds.js';
import { useConfig } from '../hooks/useConfig.js';
import { useLayout } from '../hooks/useLayout.js';
import { deserializeElement } from '../utils/deserializeElement.js';
import { DesktopSearch } from './DesktopSearch.js';
import * as styles from './DesktopTopNav.css.js';
import { NavLogo } from './NavLogo.js';
import * as NavigationMenu from './NavigationMenu.js';
import { RouterLink } from './RouterLink.js';
DesktopTopNav.Curtain = Curtain;
export function DesktopTopNav() {
    const config = useConfig();
    const { showLogo, showSidebar } = useLayout();
    return (_jsxs("div", { className: clsx(styles.root, showLogo && !showSidebar && styles.withLogo), children: [_jsx(DesktopSearch, {}), showLogo && (_jsx("div", { className: styles.logoWrapper, children: _jsx("div", { className: styles.logo, children: _jsx(RouterLink, { to: "/", style: { alignItems: 'center', display: 'flex', height: '56px', marginTop: '4px' }, children: _jsx(NavLogo, {}) }) }) })), _jsx("div", { className: styles.section }), _jsx("div", { className: styles.section, children: (config.topNav?.length || 0) > 0 && (_jsx("div", { className: styles.group, children: _jsx(Navigation, {}) })) })] }));
}
export function Curtain() {
    return _jsx("div", { className: styles.curtain });
}
function Navigation() {
    const { topNav } = useConfig();
    if (!topNav)
        return null;
    const { pathname } = useLocation();
    const activeIds = useActiveNavIds({ pathname, items: topNav });
    return (_jsx(NavigationMenu.Root, { delayDuration: 0, children: _jsxs(NavigationMenu.List, { children: [topNav.map((item, i) => {
                    if (item.element)
                        return deserializeElement(item.element);
                    if (item.link)
                        return (_jsx(NavigationMenu.Link, { active: activeIds.includes(item.id), className: styles.item, href: item.link, children: item.text }, i));
                    if (item.items)
                        return (_jsxs(NavigationMenu.Item, { className: styles.item, children: [_jsx(NavigationMenu.Trigger, { active: activeIds.includes(item.id), children: item.text }), _jsx(NavigationMenu.Content, { className: styles.content, children: _jsx(NavigationMenuContent, { items: item.items }) })] }, i));
                    return null;
                }), _jsx(TopNavEnd, {})] }) }));
}
function NavigationMenuContent({ items }) {
    const { pathname } = useLocation();
    const activeIds = useActiveNavIds({ pathname, items });
    return (_jsx("ul", { children: items?.map((item, i) => (_jsx(NavigationMenu.Link, { active: activeIds.includes(item.id), href: item.link, children: item.text }, i))) }));
}
//# sourceMappingURL=DesktopTopNav.js.map