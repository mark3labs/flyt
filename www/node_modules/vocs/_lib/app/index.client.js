import { jsx as _jsx } from "react/jsx-runtime";
import './styles/index.css.js';
import { hydrateRoot } from 'react-dom/client';
import { RouterProvider, createBrowserRouter } from 'react-router';
import { ConfigProvider, getConfig } from './hooks/useConfig.js';
import { routes } from './routes.js';
import { hydrateLazyRoutes } from './utils/hydrateLazyRoutes.js';
import { removeTempStyles } from './utils/removeTempStyles.js';
hydrate();
async function hydrate() {
    const basePath = getConfig().basePath;
    await hydrateLazyRoutes(routes, basePath);
    removeTempStyles();
    const router = createBrowserRouter(routes, { basename: basePath });
    hydrateRoot(document.getElementById('app'), _jsx(ConfigProvider, { children: _jsx(RouterProvider, { router: router }) }));
}
//# sourceMappingURL=index.client.js.map