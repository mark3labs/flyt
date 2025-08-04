import { jsx as _jsx, Fragment as _Fragment, jsxs as _jsxs } from "react/jsx-runtime";
import { Layout } from 'virtual:consumer-components';
import { MDXProvider } from '@mdx-js/react';
import { useEffect, useRef } from 'react';
import { ScrollRestoration, useLocation } from 'react-router';
import 'virtual:styles';
import { components } from './components/mdx/index.js';
import { useConfig } from './hooks/useConfig.js';
import { useOgImageUrl } from './hooks/useOgImageUrl.js';
import { PageDataContext } from './hooks/usePageData.js';
export function Root(props) {
    const { children, content, filePath, frontmatter, lastUpdatedAt, path } = props;
    const { pathname } = useLocation();
    const previousPathRef = useRef(undefined);
    useEffect(() => {
        previousPathRef.current = pathname;
    });
    return (_jsxs(_Fragment, { children: [_jsx(Head, { frontmatter: frontmatter }), typeof window !== 'undefined' && _jsx(ScrollRestoration, {}), _jsx(MDXProvider, { components: components, children: _jsx(Layout, { frontmatter: frontmatter, path: path, children: _jsx(PageDataContext.Provider, { value: {
                            content,
                            filePath,
                            frontmatter,
                            lastUpdatedAt,
                            previousPath: previousPathRef.current,
                        }, children: children }) }) })] }));
}
function Head({ frontmatter }) {
    const config = useConfig();
    const ogImageUrl = useOgImageUrl();
    const { baseUrl, font, iconUrl, logoUrl } = config;
    const description = frontmatter?.description ?? config.description;
    const title = frontmatter?.title ?? config.title;
    const titleTemplate = (() => {
        if (!config.title)
            return undefined;
        if (title.includes(config.title))
            return undefined;
        return config.titleTemplate;
    })();
    const isLocalhost = typeof window !== 'undefined' && window.location.hostname === 'localhost';
    const fullTitle = titleTemplate ? titleTemplate.replace('%s', title) : title;
    return (_jsxs(_Fragment, { children: [fullTitle && _jsx("title", { children: fullTitle }), baseUrl && import.meta.env.PROD && !isLocalhost && _jsx("base", { href: baseUrl }), description !== 'undefined' && _jsx("meta", { name: "description", content: description }), iconUrl && typeof iconUrl === 'string' && (_jsx("link", { rel: "icon", href: iconUrl, type: getIconType(iconUrl) })), iconUrl && typeof iconUrl !== 'string' && (_jsx("link", { rel: "icon", href: iconUrl.light, type: getIconType(iconUrl.light) })), iconUrl && typeof iconUrl !== 'string' && (_jsx("link", { rel: "icon", href: iconUrl.dark, type: getIconType(iconUrl.dark), media: "(prefers-color-scheme: dark)" })), _jsx("meta", { property: "og:type", content: "website" }), _jsx("meta", { property: "og:title", content: title || config.title }), baseUrl && _jsx("meta", { property: "og:url", content: baseUrl }), description !== 'undefined' && _jsx("meta", { property: "og:description", content: description }), ogImageUrl && (_jsx("meta", { property: "og:image", content: ogImageUrl
                    .replace('%logo', `${baseUrl ? baseUrl : ''}${typeof logoUrl === 'string' ? logoUrl : logoUrl?.dark || ''}`)
                    .replace('%title', title || '')
                    .replace('%description', (description !== 'undefined' ? description : '') || '') })), (font?.default?.google || font?.mono?.google) && (_jsx("link", { rel: "preconnect", href: "https://fonts.googleapis.com" })), (font?.default?.google || font?.mono?.google) && (_jsx("link", { rel: "preconnect", href: "https://fonts.gstatic.com", crossOrigin: "" })), font?.default?.google && (_jsx("link", { href: `https://fonts.googleapis.com/css2?family=${font.default.google}:wght@300;400;500&display=swap`, rel: "stylesheet" })), font?.mono?.google && (_jsx("link", { href: `https://fonts.googleapis.com/css2?family=${font.mono.google}:wght@300;400;500&display=swap`, rel: "stylesheet" })), _jsx("meta", { name: "twitter:card", content: "summary_large_image" }), ogImageUrl && (_jsx("meta", { property: "twitter:image", content: ogImageUrl
                    .replace('%logo', `${baseUrl ? baseUrl : ''}${typeof logoUrl === 'string' ? logoUrl : logoUrl?.dark || ''}`)
                    .replace('%title', title ? encodeURIComponent(title) : '')
                    .replace('%description', description ? encodeURIComponent(description) : '') }))] }));
}
function getIconType(iconUrl) {
    if (iconUrl.endsWith('.svg'))
        return 'image/svg+xml';
    if (iconUrl.endsWith('.png'))
        return 'image/png';
    if (iconUrl.endsWith('.jpg'))
        return 'image/jpeg';
    if (iconUrl.endsWith('.ico'))
        return 'image/x-icon';
    if (iconUrl.endsWith('.webp'))
        return 'image/webp';
    return undefined;
}
//# sourceMappingURL=root.js.map