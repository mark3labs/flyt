import { jsx as _jsx } from "react/jsx-runtime";
import { routes as routes_virtual } from 'virtual:routes';
import { NotFound } from './components/NotFound.js';
import { DocsLayout } from './layouts/DocsLayout.js';
import { Root } from './root.js';
const notFoundRoute = (() => {
    const virtualRoute = routes_virtual.find(({ path }) => path === '*');
    if (virtualRoute)
        return {
            path: virtualRoute.path,
            lazy: async () => {
                const { frontmatter, ...route } = await virtualRoute.lazy();
                return {
                    ...route,
                    element: (_jsx(Root, { frontmatter: frontmatter, path: virtualRoute.path, children: _jsx(DocsLayout, { children: _jsx(route.default, {}) }) })),
                };
            },
        };
    return {
        path: '*', // 404
        lazy: undefined,
        element: (_jsx(Root, { frontmatter: { layout: 'minimal' }, path: "*", children: _jsx(DocsLayout, { children: _jsx(NotFound, {}) }) })),
    };
})();
export const routes = [
    ...routes_virtual
        .filter(({ path }) => path !== '*')
        .map((route_virtual) => ({
        path: route_virtual.path,
        lazy: async () => {
            const { frontmatter, ...route } = await route_virtual.lazy();
            return {
                ...route,
                element: (_jsx(Root, { content: decodeURIComponent(route_virtual.content ?? ''), filePath: route_virtual.filePath, frontmatter: frontmatter, lastUpdatedAt: route_virtual.lastUpdatedAt, path: route_virtual.path, children: _jsx(DocsLayout, { children: _jsx(route.default, {}) }) })),
            };
        },
    })),
    notFoundRoute,
];
//# sourceMappingURL=routes.js.map