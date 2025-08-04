import * as React from 'react';
import { R as RouterProviderProps$1 } from './components-CjQijYga.mjs';
import './browser-7LYX59NK.mjs';
import { e as RouterInit } from './route-data-CqEmXQub.mjs';
import 'node:async_hooks';

type RouterProviderProps = Omit<RouterProviderProps$1, "flushSync">;
declare function RouterProvider(props: Omit<RouterProviderProps, "flushSync">): React.JSX.Element;

interface HydratedRouterProps {
    /**
     * Context object to passed through to `createBrowserRouter` and made available
     * to `clientLoader`/`clientActon` functions
     */
    unstable_getContext?: RouterInit["unstable_getContext"];
}
/**
 * Framework-mode router component to be used to to hydrate a router from a
 * `ServerRouter`. See [`entry.client.tsx`](../api/framework-conventions/entry.client.tsx).
 *
 * @public
 * @category Framework Routers
 * @mode framework
 * @param props Props
 * @param props.unstable_getContext Context object to passed through to
 * {@link createBrowserRouter} and made available to `clientLoader`/`clientAction`
 * functions
 * @returns A React element that represents the hydrated application.
 */
declare function HydratedRouter(props: HydratedRouterProps): React.JSX.Element;

export { HydratedRouter, RouterProvider, type RouterProviderProps };
