import { jsx as _jsx, jsxs as _jsxs, Fragment as _Fragment } from "react/jsx-runtime";
import clsx from 'clsx';
import { useInView } from 'react-intersection-observer';
import { assignInlineVars } from '@vanilla-extract/dynamic';
import { bannerHeight } from '../components/Banner.css.js';
import { Banner } from '../components/Banner.js';
import { Content } from '../components/Content.js';
import { DesktopTopNav } from '../components/DesktopTopNav.js';
import { Footer } from '../components/Footer.js';
import { MobileTopNav } from '../components/MobileTopNav.js';
import { Outline } from '../components/Outline.js';
import { Sidebar } from '../components/Sidebar.js';
import { SkipLink, skipLinkId } from '../components/SkipLink.js';
import { useConfig } from '../hooks/useConfig.js';
import { useLayout } from '../hooks/useLayout.js';
import { useLocalStorage } from '../hooks/useLocalStorage.js';
import { usePageData } from '../hooks/usePageData.js';
import { contentVars, defaultFontFamily, fontFamilyVars } from '../styles/vars.css.js';
import * as styles from './DocsLayout.css.js';
export function DocsLayout({ children, }) {
    const { banner, font } = useConfig();
    const { frontmatter = {} } = usePageData();
    const { content } = frontmatter;
    const { layout, showOutline, showSidebar, showTopNav } = useLayout();
    const { ref, inView } = useInView({
        initialInView: true,
        rootMargin: '100px 0px 0px 0px',
    });
    const [showBanner, setShowBanner] = useLocalStorage('banner', true);
    return (_jsxs("div", { className: styles.root, "data-layout": layout, style: assignInlineVars({
            [bannerHeight]: showBanner ? banner?.height : undefined,
            [fontFamilyVars.default]: font?.default?.google
                ? `${font.default.google}, ${defaultFontFamily.default}`
                : undefined,
            [fontFamilyVars.mono]: font?.mono?.google
                ? `${font.mono.google}, ${defaultFontFamily.mono}`
                : undefined,
        }), children: [_jsx(SkipLink, {}), showBanner && _jsx(Banner, { hide: () => setShowBanner(false) }), showSidebar && (_jsx("div", { className: styles.gutterLeft, children: _jsx(Sidebar, { className: styles.sidebar }) })), showTopNav && (_jsxs(_Fragment, { children: [_jsxs("div", { ref: ref, className: clsx(styles.gutterTop, showSidebar && styles.gutterTop_offsetLeftGutter, (layout === 'minimal' || layout === 'landing') && styles.gutterTop_sticky), children: [_jsx(DesktopTopNav, {}), _jsx(MobileTopNav, {})] }), _jsxs("div", { className: clsx(styles.gutterTopCurtain, showSidebar && styles.gutterTopCurtain_withSidebar, (layout === 'minimal' || layout === 'landing') && styles.gutterTopCurtain_hidden), children: [_jsx(DesktopTopNav.Curtain, {}), _jsx(MobileTopNav.Curtain, { enableScrollToTop: !inView })] })] })), showOutline && (_jsx("div", { className: clsx(styles.gutterRight, showSidebar && styles.gutterRight_withSidebar), children: _jsx(Outline, {}) })), _jsxs("div", { id: skipLinkId, className: clsx(styles.content, showSidebar && styles.content_withSidebar, showTopNav && styles.content_withTopNav), style: assignInlineVars({
                    [contentVars.horizontalPadding]: content?.horizontalPadding,
                    [contentVars.width]: content?.width,
                    [contentVars.verticalPadding]: content?.verticalPadding,
                }), children: [_jsx(Content, { children: children }), _jsx(Footer, {})] }), _jsx("div", { "data-bottom-observer": true })] }));
}
//# sourceMappingURL=DocsLayout.js.map