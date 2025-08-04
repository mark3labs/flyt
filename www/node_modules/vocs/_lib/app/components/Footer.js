import { jsx as _jsx, jsxs as _jsxs, Fragment as _Fragment } from "react/jsx-runtime";
import { Footer as ConsumerFooter } from 'virtual:consumer-components';
import { Pencil2Icon } from '@radix-ui/react-icons';
import { assignInlineVars } from '@vanilla-extract/dynamic';
import clsx from 'clsx';
import { useEffect, useMemo } from 'react';
import { useLocation, useNavigate } from 'react-router';
import { useEditLink } from '../hooks/useEditLink.js';
import { useLayout } from '../hooks/useLayout.js';
import { useMounted } from '../hooks/useMounted.js';
import { usePageData } from '../hooks/usePageData.js';
import { useSidebar } from '../hooks/useSidebar.js';
import * as styles from './Footer.css.js';
import { sizeVar } from './Icon.css.js';
import { Icon } from './Icon.js';
import { KeyboardShortcut } from './KeyboardShortcut.js';
import { Link } from './Link.js';
import { Socials } from './Socials.js';
import { ArrowLeft } from './icons/ArrowLeft.js';
import { ArrowRight } from './icons/ArrowRight.js';
export function Footer() {
    const { layout } = useLayout();
    const mounted = useMounted();
    const pageData = usePageData();
    const lastUpdatedAtDate = useMemo(() => (pageData.lastUpdatedAt ? new Date(pageData.lastUpdatedAt) : undefined), [pageData.lastUpdatedAt]);
    const lastUpdatedAtISOString = useMemo(() => lastUpdatedAtDate?.toISOString(), [lastUpdatedAtDate]);
    return (_jsxs("footer", { className: styles.root, "data-layout": layout, children: [layout === 'docs' && (_jsxs(_Fragment, { children: [_jsxs("div", { className: styles.container, children: [_jsx(EditLink, {}), mounted && pageData.lastUpdatedAt && (_jsxs("div", { className: styles.lastUpdated, children: ["Last updated:", ' ', _jsx("time", { dateTime: lastUpdatedAtISOString, children: new Intl.DateTimeFormat(undefined, {
                                            dateStyle: 'short',
                                            timeStyle: 'short',
                                        }).format(lastUpdatedAtDate) })] }))] }), _jsx(Navigation, {})] })), layout !== 'docs' && _jsx(Socials, {}), _jsx(ConsumerFooter, {})] }));
}
function EditLink() {
    const editLink = useEditLink();
    if (!editLink.url)
        return null;
    return (_jsx("div", { children: _jsxs(Link, { className: styles.editLink, href: editLink.url, children: [_jsx(Pencil2Icon, {}), " ", editLink.text] }) }));
}
function Navigation() {
    const mounted = useMounted();
    const sidebar = useSidebar();
    const { pathname } = useLocation();
    const flattenedSidebar = useMemo(() => flattenSidebar(sidebar.items || []).filter((item) => item.link), [sidebar]);
    const currentPageIndex = useMemo(() => flattenedSidebar.findIndex((item) => item.link === pathname), [flattenedSidebar, pathname]);
    const [prevPage, nextPage] = useMemo(() => {
        if (currentPageIndex < 0)
            return [];
        if (currentPageIndex === 0)
            return [null, flattenedSidebar[currentPageIndex + 1]];
        if (currentPageIndex === flattenedSidebar.length - 1)
            return [flattenedSidebar[currentPageIndex - 1], null];
        return [flattenedSidebar[currentPageIndex - 1], flattenedSidebar[currentPageIndex + 1]];
    }, [currentPageIndex, flattenedSidebar]);
    const navigate = useNavigate();
    // biome-ignore lint/correctness/useExhaustiveDependencies: <explanation>
    useEffect(() => {
        let index = currentPageIndex;
        let isListening = false;
        const keydown = (event) => {
            if (event.code === 'ShiftLeft')
                isListening = true;
            if (isListening) {
                const nextPage = flattenedSidebar[index + 1];
                const prevPage = flattenedSidebar[index - 1];
                if (event.code === 'ArrowRight' && nextPage?.link) {
                    const isExternalLink = !nextPage.link.match(/^(\.*\/|#)/);
                    if (isExternalLink)
                        window.open(nextPage.link, 'noopener,noreferrer');
                    else
                        navigate(nextPage.link);
                    index++;
                }
                if (event.code === 'ArrowLeft' && prevPage?.link) {
                    const isExternalLink = !prevPage.link.match(/^(\.*\/|#)/);
                    if (isExternalLink)
                        window.open(prevPage.link, 'noopener,noreferrer');
                    else
                        navigate(prevPage.link);
                    index--;
                }
            }
        };
        const keyup = (event) => {
            if (event.code === 'ShiftLeft')
                isListening = false;
        };
        window.addEventListener('keydown', keydown);
        window.addEventListener('keyup', keyup);
        return () => {
            window.removeEventListener('keydown', keydown);
            window.removeEventListener('keyup', keyup);
        };
    }, []);
    if (!mounted)
        return null;
    return (_jsxs("div", { className: styles.navigation, children: [prevPage ? (_jsxs(Link, { className: clsx(styles.navigationItem, styles.navigationItem_left), href: prevPage.link, variant: "styleless", children: [_jsxs("div", { className: styles.navigationText, children: [_jsx("div", { className: clsx(styles.navigationIcon, styles.navigationIcon_left), style: assignInlineVars({ [sizeVar]: '0.75em' }), children: _jsx(Icon, { label: "Previous", icon: ArrowLeft }) }), _jsx("div", { className: styles.navigationTextInner, children: prevPage.text })] }), _jsx(KeyboardShortcut, { description: "Previous", keys: ['shift', '←'] })] })) : (_jsx("div", {})), nextPage ? (_jsxs(Link, { className: clsx(styles.navigationItem, styles.navigationItem_right), href: nextPage.link, variant: "styleless", children: [_jsxs("div", { className: styles.navigationText, children: [_jsx("div", { className: styles.navigationTextInner, style: { textAlign: 'right' }, children: nextPage.text }), _jsx("div", { className: clsx(styles.navigationIcon, styles.navigationIcon_right), style: assignInlineVars({ [sizeVar]: '0.75em' }), children: _jsx(Icon, { label: "Next", icon: ArrowRight }) })] }), _jsx(KeyboardShortcut, { description: "Next", keys: ['shift', '→'] })] })) : (_jsx("div", {}))] }));
}
function flattenSidebar(sidebar) {
    const items = [];
    for (const item of sidebar) {
        if (item.link) {
            items.push(item);
        }
        if (item.items) {
            items.push(...flattenSidebar(item.items));
        }
    }
    return items;
}
//# sourceMappingURL=Footer.js.map