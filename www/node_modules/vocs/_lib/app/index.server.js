import { jsx as _jsx } from "react/jsx-runtime";
import './styles/index.css.js';
import { Route, Routes } from 'react-router';
import { StaticRouter, StaticRouterProvider, createStaticHandler, createStaticRouter, } from 'react-router';
import { resolveVocsConfig } from '../vite/utils/resolveVocsConfig.js';
import { ConfigProvider } from './hooks/useConfig.js';
import { routes } from './routes.js';
import { createFetchRequest } from './utils/createFetchRequest.js';
export async function prerender(location) {
    const unwrappedRoutes = (await Promise.all(routes.map(async (route) => {
        const location_ = location === '/' ? '/' : location.replace(/\/$/, '');
        const path = route.path.replace(/\.html$/, '');
        if (path !== location_ && path !== '*')
            return null;
        const element = route.lazy ? (await route.lazy()).element : route.element;
        return {
            path: route.path,
            element,
        };
    }))).filter(Boolean);
    const { config } = await resolveVocsConfig();
    const { basePath } = config;
    return (_jsx(ConfigProvider, { config: config, children: _jsx(StaticRouter, { location: location, basename: basePath, children: _jsx(Routes, { children: unwrappedRoutes.map((route) => (_jsx(Route, { path: route.path, element: route.element }, route.path))) }) }) }));
}
export async function render(req) {
    const { config } = await resolveVocsConfig();
    const { basePath } = config;
    const { query, dataRoutes } = createStaticHandler(routes, { basename: basePath });
    const fetchRequest = createFetchRequest(req);
    const context = (await query(fetchRequest));
    if (context instanceof Response)
        throw context;
    const router = createStaticRouter(dataRoutes, context);
    return (_jsx(ConfigProvider, { config: config, children: _jsx(StaticRouterProvider, { router: router, context: context }) }));
}
//# sourceMappingURL=index.server.js.map