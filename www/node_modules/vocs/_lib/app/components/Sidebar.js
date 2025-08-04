import { jsx as _jsx, jsxs as _jsxs, Fragment as _Fragment } from "react/jsx-runtime";
import clsx from 'clsx';
import { useCallback, useEffect, useMemo, useRef, useState, } from 'react';
import { matchPath, useLocation, useMatch } from 'react-router';
import { useConfig } from '../hooks/useConfig.js';
import { usePageData } from '../hooks/usePageData.js';
import { useSidebar } from '../hooks/useSidebar.js';
import { Icon } from './Icon.js';
import { Link } from './Link.js';
import { NavLogo } from './NavLogo.js';
import { RouterLink } from './RouterLink.js';
import * as styles from './Sidebar.css.js';
import { Socials } from './Socials.js';
import { ThemeToggle } from './ThemeToggle.js';
import { ChevronRight } from './icons/ChevronRight.js';
export function Sidebar(props) {
    const { className, onClickItem } = props;
    const { theme } = useConfig();
    const { previousPath } = usePageData();
    const sidebarRef = useRef(null);
    const sidebar = useSidebar();
    const [backPath, setBackPath] = useState('/');
    // biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
    useEffect(() => {
        if (typeof window === 'undefined')
            return;
        if (!previousPath)
            return;
        setBackPath(previousPath);
    }, [sidebar.key, sidebar.backLink]);
    if (!sidebar)
        return null;
    const groups = getSidebarGroups(sidebar.items);
    return (_jsxs("aside", { ref: sidebarRef, className: clsx(styles.root, className), children: [_jsxs("div", { children: [_jsxs("div", { className: styles.logoWrapper, children: [_jsx("div", { className: styles.logo, children: _jsx(RouterLink, { to: "/", style: { alignItems: 'center', display: 'flex', height: '100%' }, children: _jsx(NavLogo, {}) }) }), _jsx("div", { className: styles.divider })] }), _jsx("nav", { className: styles.navigation, children: _jsxs("div", { className: styles.group, children: [sidebar.backLink && (_jsx("section", { className: styles.section, children: _jsx("div", { className: styles.items, children: _jsxs(RouterLink, { className: clsx(styles.item, styles.backLink), to: backPath, children: ["\u2190", ' ', typeof history !== 'undefined' && history.state?.key && backPath !== '/'
                                                    ? 'Back'
                                                    : 'Home'] }) }) })), groups.map((group, i) => (_jsx(SidebarItem, { depth: 0, item: group, onClick: onClickItem, sidebarRef: sidebarRef }, `${group.text}${i}`)))] }) })] }), _jsxs("div", { className: styles.footer, children: [_jsx("div", { className: styles.footerCurtain }), _jsxs("div", { className: styles.footerContent, children: [_jsx(Socials, {}), !theme?.colorScheme ? _jsx(ThemeToggle, {}) : null] })] })] }, sidebar.key));
}
function getSidebarGroups(sidebar) {
    const groups = [];
    let lastGroupIndex = 0;
    for (const item of sidebar) {
        if (item.items) {
            lastGroupIndex = groups.push(item);
            continue;
        }
        if (!groups[lastGroupIndex])
            groups.push({ text: '', items: [item] });
        else
            groups[lastGroupIndex].items.push(item);
    }
    return groups;
}
function getActiveChildItem(items, pathname) {
    return items.find((item) => {
        if (matchPath(pathname, item.link ?? ''))
            return true;
        if (item.link === pathname)
            return true;
        if (!item.items)
            return false;
        return getActiveChildItem(item.items, pathname);
    });
}
function SidebarItem(props) {
    const { depth, item, onClick, sidebarRef } = props;
    const itemRef = useRef(null);
    const { pathname } = useLocation();
    const match = useMatch(item.link || '');
    const hasActiveChildItem = useMemo(() => (item.items ? Boolean(getActiveChildItem(item.items, pathname)) : false), [item.items, pathname]);
    const [collapsed, setCollapsed] = useState(() => {
        if (item.link && match)
            return false;
        if (!item.items)
            return false;
        if (hasActiveChildItem)
            return false;
        return Boolean(item.collapsed);
    });
    const isCollapsable = item.collapsed !== undefined && item.items !== undefined;
    const onCollapseInteraction = useCallback((event) => {
        if ('key' in event && event.key !== 'Enter')
            return;
        setCollapsed((x) => !x);
    }, []);
    const onCollapseTriggerInteraction = useCallback((event) => {
        if ('key' in event && event.key !== 'Enter')
            return;
        setCollapsed((x) => !x);
    }, []);
    const active = useRef(true);
    useEffect(() => {
        if (!active.current)
            return;
        active.current = false;
        const match = matchPath(pathname, item.link ?? '');
        if (!match)
            return;
        requestAnimationFrame(() => {
            const offsetTop = itemRef.current?.offsetTop ?? 0;
            const sidebarHeight = sidebarRef?.current?.clientHeight ?? 0;
            if (offsetTop < sidebarHeight)
                return;
            sidebarRef?.current?.scrollTo({ top: offsetTop - 100 });
        });
    }, [item, pathname, sidebarRef]);
    if (item.items)
        return (_jsxs("section", { ref: itemRef, className: clsx(styles.section, depth === 0 && item.text && styles.level, depth === 0 && item.text && collapsed && styles.levelCollapsed), children: [item.text && (_jsxs("div", { className: styles.sectionHeader, ...(isCollapsable && !item.link
                        ? {
                            role: 'button',
                            tabIndex: 0,
                            onClick: onCollapseInteraction,
                            onKeyDown: onCollapseInteraction,
                        }
                        : {}), children: [item.text &&
                            (item.link ? (_jsx(Link, { "data-active": Boolean(match), onClick: (e) => {
                                    onClick?.(e);
                                    onCollapseInteraction(e);
                                }, className: clsx(depth === 0 ? [styles.sectionTitle, styles.sectionTitleLink] : styles.item, hasActiveChildItem && styles.sectionHeaderActive, item.disabled && styles.disabledItem), href: item.link, variant: "styleless", children: item.text })) : (_jsx("div", { className: clsx(depth === 0 ? styles.sectionTitle : styles.item, item.disabled && styles.disabledItem), children: item.text }))), isCollapsable && (_jsx("div", { 
                            // biome-ignore lint/a11y/useSemanticElements:
                            role: "button", tabIndex: 0, onClick: onCollapseTriggerInteraction, onKeyDown: onCollapseTriggerInteraction, children: _jsx(Icon, { className: clsx(styles.sectionCollapse, collapsed && styles.sectionCollapseActive), label: "toggle section", icon: ChevronRight, size: "16px" }) }))] })), !collapsed && (_jsx("div", { className: clsx(styles.items, depth !== 0 && styles.levelInset), children: item.items &&
                        item.items.length > 0 &&
                        depth < 5 &&
                        item.items.map((item, i) => (_jsx(SidebarItem, { depth: depth + 1, item: item, onClick: onClick, sidebarRef: sidebarRef }, `${item.text}${i}`))) }))] }));
    return (_jsx(_Fragment, { children: item.link ? (_jsx(Link, { ref: itemRef, "data-active": Boolean(match), onClick: onClick, className: clsx(styles.item, item.disabled && styles.disabledItem), href: item.link, variant: "styleless", children: item.text })) : (_jsx("div", { className: clsx(styles.item, item.disabled && styles.disabledItem), children: item.text })) }));
}
//# sourceMappingURL=Sidebar.js.map