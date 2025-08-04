import { R as Router, y as RouteModules, D as DataStrategyFunction, T as To, z as NavigateOptions, B as BlockerFunction, E as Blocker, J as SerializeFrom, a as RelativeRoutingType, c as Location, K as ParamParseKey, O as Path, Q as PathPattern, V as PathMatch, U as UIMatch, W as Navigation, A as Action, h as Params, f as RouteObject, I as InitialEntry, H as HydrationState, b as IndexRouteObject, m as LoaderFunction, i as ActionFunction, M as MetaFunction, l as LinksFunction, N as NonIndexRouteObject, X as AppLoadContext, u as unstable_RouterContextProvider, Y as MiddlewareEnabled, Z as unstable_InitialContext, _ as Equal, e as RouterInit, $ as RouterState, P as PatchRoutesOnNavigationFunction, a0 as DataRouteObject, j as ClientLoaderFunction } from './route-data-CqEmXQub.mjs';
export { ab as ActionFunctionArgs, C as ClientActionFunction, aH as ClientActionFunctionArgs, aI as ClientLoaderFunctionArgs, aF as DataRouteMatch, ac as DataStrategyFunctionArgs, ad as DataStrategyMatch, ae as DataStrategyResult, ag as ErrorResponse, a5 as Fetcher, ah as FormEncType, ai as FormMethod, aN as Future, a2 as GetScrollPositionFunction, a3 as GetScrollRestorationKeyFunction, aj as HTMLFormMethod, aJ as HeadersArgs, k as HeadersFunction, aM as HtmlLinkDescriptor, av as IDLE_BLOCKER, au as IDLE_FETCHER, at as IDLE_NAVIGATION, L as LazyRouteFunction, o as LinkDescriptor, ak as LoaderFunctionArgs, aK as MetaArgs, q as MetaDescriptor, a6 as NavigationStates, d as Navigator, aL as PageLinkDescriptor, aG as PatchRoutesOnNavigationFunctionArgs, am as PathParam, an as RedirectFunction, aa as RevalidationState, g as RouteMatch, a9 as RouterFetchOptions, a8 as RouterNavigateOptions, a7 as RouterSubscriber, S as ShouldRevalidateFunction, ap as ShouldRevalidateFunctionArgs, a1 as StaticHandler, a4 as StaticHandlerContext, aT as UNSAFE_DataRouterContext, aU as UNSAFE_DataRouterStateContext, af as UNSAFE_DataWithResponseInit, aS as UNSAFE_ErrorResponseImpl, aV as UNSAFE_FetchersContext, aW as UNSAFE_LocationContext, aX as UNSAFE_NavigationContext, aY as UNSAFE_RouteContext, aZ as UNSAFE_ViewTransitionContext, aP as UNSAFE_createBrowserHistory, aR as UNSAFE_createRouter, aQ as UNSAFE_invariant, ar as createPath, aw as data, ax as generatePath, ay as isRouteErrorResponse, az as matchPath, aA as matchRoutes, as as parsePath, aB as redirect, aC as redirectDocument, aD as replace, aE as resolvePath, al as unstable_MiddlewareFunction, ao as unstable_RouterContext, aO as unstable_SerializesTo, aq as unstable_createContext } from './route-data-CqEmXQub.mjs';
import { a as RouteComponentType, H as HydrateFallbackType, E as ErrorBoundaryType } from './components-CjQijYga.mjs';
export { f as Await, A as AwaitProps, I as IndexRouteProps, L as LayoutRouteProps, g as MemoryRouter, M as MemoryRouterOpts, b as MemoryRouterProps, h as Navigate, N as NavigateProps, i as Outlet, O as OutletProps, P as PathRouteProps, j as Route, c as RouteProps, k as Router, d as RouterProps, l as RouterProvider, R as RouterProviderProps, m as Routes, e as RoutesProps, W as UNSAFE_WithComponentProps, v as UNSAFE_WithErrorBoundaryProps, t as UNSAFE_WithHydrateFallbackProps, q as UNSAFE_hydrationRouteProperties, s as UNSAFE_mapRouteProperties, w as UNSAFE_withComponentProps, x as UNSAFE_withErrorBoundaryProps, u as UNSAFE_withHydrateFallbackProps, n as createMemoryRouter, o as createRoutesFromChildren, p as createRoutesFromElements, r as renderMatches } from './components-CjQijYga.mjs';
import * as React from 'react';
import React__default, { ReactElement } from 'react';
import { A as AssetsManifest, E as EntryContext, F as FutureConfig, S as ServerBuild } from './index-react-server-client-KLg-U4nr.mjs';
export { l as BrowserRouter, B as BrowserRouterProps, D as DOMRouterOpts, a0 as DiscoverBehavior, c as FetcherFormProps, h as FetcherSubmitFunction, G as FetcherSubmitOptions, i as FetcherWithComponents, q as Form, d as FormProps, a1 as HandleDataRequestFunction, a2 as HandleDocumentRequestFunction, a3 as HandleErrorFunction, m as HashRouter, H as HashRouterProps, a as HistoryRouterProps, n as Link, L as LinkProps, X as Links, W as Meta, p as NavLink, N as NavLinkProps, b as NavLinkRenderProps, P as ParamKeyValuePair, $ as PrefetchBehavior, Z as PrefetchPageLinks, Y as Scripts, _ as ScriptsProps, r as ScrollRestoration, e as ScrollRestorationProps, a4 as ServerEntryModule, f as SetURLSearchParams, T as StaticRouter, M as StaticRouterProps, V as StaticRouterProvider, O as StaticRouterProviderProps, g as SubmitFunction, I as SubmitOptions, J as SubmitTarget, a5 as UNSAFE_FrameworkContext, a6 as UNSAFE_createClientRoutes, a7 as UNSAFE_createClientRoutesWithHMRRevalidationOptOut, a8 as UNSAFE_shouldHydrateRouteLoader, a9 as UNSAFE_useScrollRestoration, U as URLSearchParamsInit, j as createBrowserRouter, k as createHashRouter, K as createSearchParams, Q as createStaticHandler, R as createStaticRouter, o as unstable_HistoryRouter, z as unstable_usePrompt, y as useBeforeUnload, w as useFetcher, x as useFetchers, v as useFormAction, u as useLinkClickHandler, s as useSearchParams, t as useSubmit, C as useViewTransitionState } from './index-react-server-client-KLg-U4nr.mjs';
import { ParseOptions, SerializeOptions } from 'cookie';
export { ParseOptions as CookieParseOptions, SerializeOptions as CookieSerializeOptions } from 'cookie';
import { P as Pages } from './register-DiOIlEq5.mjs';
export { a as Register } from './register-DiOIlEq5.mjs';
import { R as RSCPayload, m as matchRSCServerRequest } from './browser-7LYX59NK.mjs';
export { D as unstable_DecodeActionFunction, a as unstable_DecodeFormStateFunction, b as unstable_DecodeReplyFunction, L as unstable_LoadServerActionFunction, c as unstable_RSCManifestPayload, d as unstable_RSCMatch, e as unstable_RSCRenderPayload, j as unstable_RSCRouteConfig, i as unstable_RSCRouteConfigEntry, f as unstable_RSCRouteManifest, h as unstable_RSCRouteMatch, g as unstable_getRSCStream } from './browser-7LYX59NK.mjs';
import 'node:async_hooks';

declare const SingleFetchRedirectSymbol: unique symbol;
declare function getTurboStreamSingleFetchDataStrategy(getRouter: () => Router, manifest: AssetsManifest, routeModules: RouteModules, ssr: boolean, basename: string | undefined): DataStrategyFunction;
declare function decodeViaTurboStream(body: ReadableStream<Uint8Array>, global: Window | typeof globalThis): Promise<{
    done: Promise<undefined>;
    value: unknown;
}>;

/**
 * The mode to use when running the server.
 */
declare enum ServerMode {
    Development = "development",
    Production = "production",
    Test = "test"
}

/**
 * Resolves a URL against the current {@link Location}.
 *
 * @example
 * import { useHref } from "react-router";
 *
 * function SomeComponent() {
 *   let href = useHref("some/where");
 *   // "/resolved/some/where"
 * }
 *
 * @public
 * @category Hooks
 * @param to The path to resolve
 * @param options Options
 * @param options.relative Defaults to `"route"` so routing is relative to the
 * route tree.
 * Set to `"path"` to make relative routing operate against path segments.
 * @returns The resolved href string
 */
declare function useHref(to: To, { relative }?: {
    relative?: RelativeRoutingType;
}): string;
/**
 * Returns `true` if this component is a descendant of a {@link Router}, useful
 * to ensure a component is used within a {@link Router}.
 *
 * @public
 * @category Hooks
 * @mode framework
 * @mode data
 * @returns Whether the component is within a {@link Router} context
 */
declare function useInRouterContext(): boolean;
/**
 * Returns the current {@link Location}. This can be useful if you'd like to
 * perform some side effect whenever it changes.
 *
 * @example
 * import * as React from 'react'
 * import { useLocation } from 'react-router'
 *
 * function SomeComponent() {
 *   let location = useLocation()
 *
 *   React.useEffect(() => {
 *     // Google Analytics
 *     ga('send', 'pageview')
 *   }, [location]);
 *
 *   return (
 *     // ...
 *   );
 * }
 *
 * @public
 * @category Hooks
 * @returns The current {@link Location} object
 */
declare function useLocation(): Location;
/**
 * Returns the current {@link Navigation} action which describes how the router
 * came to the current {@link Location}, either by a pop, push, or replace on
 * the [`History`](https://developer.mozilla.org/en-US/docs/Web/API/History) stack.
 *
 * @public
 * @category Hooks
 * @returns The current {@link NavigationType} type (`"PUSH"`, `"REPLACE"`, or `"POP"`)
 */
declare function useNavigationType(): Action;
/**
 * Returns a {@link PathMatch} object if the given pattern matches the current URL.
 * This is useful for components that need to know "active" state, e.g.
 * {@link NavLink | `<NavLink>`}.
 *
 * @public
 * @category Hooks
 * @param pattern The pattern to match against the current {@link Location}
 * @returns The path match object if the pattern matches, `null` otherwise
 */
declare function useMatch<ParamKey extends ParamParseKey<Path>, Path extends string>(pattern: PathPattern<Path> | Path): PathMatch<ParamKey> | null;
/**
 * The interface for the `navigate` function returned from {@link useNavigate}.
 */
interface NavigateFunction {
    (to: To, options?: NavigateOptions): void | Promise<void>;
    (delta: number): void | Promise<void>;
}
/**
 * Returns a function that lets you navigate programmatically in the browser in
 * response to user interactions or effects.
 *
 * It's often better to use {@link redirect} in [`action`](../../start/framework/route-module#action)/[`loader`](../../start/framework/route-module#loader)
 * functions than this hook.
 *
 * The returned function signature is `navigate(to, options?)`/`navigate(delta)` where:
 *
 * * `to` can be a string path, a {@link To} object, or a number (delta)
 * * `options` contains options for modifying the navigation
 *   * `flushSync`: Wrap the DOM updates in [`ReactDom.flushSync`](https://react.dev/reference/react-dom/flushSync)
 *   * `preventScrollReset`: Do not scroll back to the top of the page after navigation
 *   * `relative`: `"route"` or `"path"` to control relative routing logic
 *   * `replace`: Replace the current entry in the [`History`](https://developer.mozilla.org/en-US/docs/Web/API/History) stack
 *   * `state`: Optional [`history.state`](https://developer.mozilla.org/en-US/docs/Web/API/History/state) to include with the new {@link Location}
 *   * `viewTransition`: Enable [`document.startViewTransition`](https://developer.mozilla.org/en-US/docs/Web/API/Document/startViewTransition) for this navigation
 *
 * @example
 * import { useNavigate } from "react-router";
 *
 * function SomeComponent() {
 *   let navigate = useNavigate();
 *   return (
 *     <button onClick={() => navigate(-1)}>
 *       Go Back
 *     </button>
 *   );
 * }
 *
 * @additionalExamples
 * ### Navigate to another path
 *
 * ```tsx
 * navigate("/some/route");
 * navigate("/some/route?search=param");
 * ```
 *
 * ### Navigate with a {@link To} object
 *
 * All properties are optional.
 *
 * ```tsx
 * navigate({
 *   pathname: "/some/route",
 *   search: "?search=param",
 *   hash: "#hash",
 *   state: { some: "state" },
 * });
 * ```
 *
 * If you use `state`, that will be available on the {@link Location} object on
 * the next page. Access it with `useLocation().state` (see {@link useLocation}).
 *
 * ### Navigate back or forward in the history stack
 *
 * ```tsx
 * // back
 * // often used to close modals
 * navigate(-1);
 *
 * // forward
 * // often used in a multistep wizard workflows
 * navigate(1);
 * ```
 *
 * Be cautious with `navigate(number)`. If your application can load up to a
 * route that has a button that tries to navigate forward/back, there may not be
 * a `[`History`](https://developer.mozilla.org/en-US/docs/Web/API/History)
 * entry to go back or forward to, or it can go somewhere you don't expect
 * (like a different domain).
 *
 * Only use this if you're sure they will have an entry in the [`History`](https://developer.mozilla.org/en-US/docs/Web/API/History)
 * stack to navigate to.
 *
 * ### Replace the current entry in the history stack
 *
 * This will remove the current entry in the [`History`](https://developer.mozilla.org/en-US/docs/Web/API/History)
 * stack, replacing it with a new one, similar to a server side redirect.
 *
 * ```tsx
 * navigate("/some/route", { replace: true });
 * ```
 *
 * ### Prevent Scroll Reset
 *
 * [MODES: framework, data]
 *
 * <br/>
 * <br/>
 *
 * To prevent {@link ScrollRestoration | `<ScrollRestoration>`} from resetting
 * the scroll position, use the `preventScrollReset` option.
 *
 * ```tsx
 * navigate("?some-tab=1", { preventScrollReset: true });
 * ```
 *
 * For example, if you have a tab interface connected to search params in the
 * middle of a page, and you don't want it to scroll to the top when a tab is
 * clicked.
 *
 * @public
 * @category Hooks
 * @returns A navigate function for programmatic navigation
 */
declare function useNavigate(): NavigateFunction;
/**
 * Returns the parent route {@link Outlet | `<Outlet context>`}.
 *
 * Often parent routes manage state or other values you want shared with child
 * routes. You can create your own [context provider](https://react.dev/learn/passing-data-deeply-with-context)
 * if you like, but this is such a common situation that it's built-into
 * `<Outlet />`.
 *
 * ```tsx
 * // Parent route
 * function Parent() {
 *   const [count, setCount] = React.useState(0);
 *   return <Outlet context={[count, setCount]} />;
 * }
 * ```
 *
 * ```tsx
 * // Child route
 * import { useOutletContext } from "react-router-dom";
 *
 * function Child() {
 *   const [count, setCount] = useOutletContext();
 *   const increment = () => setCount((c) => c + 1);
 *   return <button onClick={increment}>{count}</button>;
 * }
 * ```
 *
 * If you're using TypeScript, we recommend the parent component provide a custom
 * hook for accessing the context value. This makes it easier for consumers to
 * get nice typings, control consumers, and know who's consuming the context value.
 *
 * Here's a more realistic example:
 *
 * ```tsx filename=src/routes/dashboard.tsx lines=[13,19]
 * import * as React from "react";
 * import type { User } from "./types";
 * import { Outlet, useOutletContext } from "react-router-dom";
 *
 * type ContextType = { user: User | null };
 *
 * export default function Dashboard() {
 *   const [user, setUser] = React.useState<User | null>(null);
 *
 *   return (
 *     <div>
 *       <h1>Dashboard</h1>
 *       <Outlet context={{ user } satisfies ContextType} />
 *     </div>
 *   );
 * }
 *
 * export function useUser() {
 *   return useOutletContext<ContextType>();
 * }
 * ```
 *
 * ```tsx filename=src/routes/dashboard/messages.tsx lines=[1,4]
 * import { useUser } from "../dashboard";
 *
 * export default function DashboardMessages() {
 *   const { user } = useUser();
 *   return (
 *     <div>
 *       <h2>Messages</h2>
 *       <p>Hello, {user.name}!</p>
 *     </div>
 *   );
 * }
 * ```
 *
 * @public
 * @category Hooks
 * @returns The context value passed to the parent {@link Outlet} component
 */
declare function useOutletContext<Context = unknown>(): Context;
/**
 * Returns the element for the child route at this level of the route
 * hierarchy. Used internally by {@link Outlet | `<Outlet>`} to render child
 * routes.
 *
 * @public
 * @category Hooks
 * @param context The context to pass to the outlet
 * @returns The child route element or `null` if no child routes match
 */
declare function useOutlet(context?: unknown): React.ReactElement | null;
/**
 * Returns an object of key/value-pairs of the dynamic params from the current
 * URL that were matched by the routes. Child routes inherit all params from
 * their parent routes.
 *
 * Assuming a route pattern like `/posts/:postId` is matched by `/posts/123`
 * then `params.postId` will be `"123"`.
 *
 * @example
 * import { useParams } from "react-router";
 *
 * function SomeComponent() {
 *   let params = useParams();
 *   params.postId;
 * }
 *
 * @additionalExamples
 * ### Basic Usage
 *
 * ```tsx
 * import { useParams } from "react-router";
 *
 * // given a route like:
 * <Route path="/posts/:postId" element={<Post />} />;
 *
 * // or a data route like:
 * createBrowserRouter([
 *   {
 *     path: "/posts/:postId",
 *     component: Post,
 *   },
 * ]);
 *
 * // or in routes.ts
 * route("/posts/:postId", "routes/post.tsx");
 * ```
 *
 * Access the params in a component:
 *
 * ```tsx
 * import { useParams } from "react-router";
 *
 * export default function Post() {
 *   let params = useParams();
 *   return <h1>Post: {params.postId}</h1>;
 * }
 * ```
 *
 * ### Multiple Params
 *
 * Patterns can have multiple params:
 *
 * ```tsx
 * "/posts/:postId/comments/:commentId";
 * ```
 *
 * All will be available in the params object:
 *
 * ```tsx
 * import { useParams } from "react-router";
 *
 * export default function Post() {
 *   let params = useParams();
 *   return (
 *     <h1>
 *       Post: {params.postId}, Comment: {params.commentId}
 *     </h1>
 *   );
 * }
 * ```
 *
 * ### Catchall Params
 *
 * Catchall params are defined with `*`:
 *
 * ```tsx
 * "/files/*";
 * ```
 *
 * The matched value will be available in the params object as follows:
 *
 * ```tsx
 * import { useParams } from "react-router";
 *
 * export default function File() {
 *   let params = useParams();
 *   let catchall = params["*"];
 *   // ...
 * }
 * ```
 *
 * You can destructure the catchall param:
 *
 * ```tsx
 * export default function File() {
 *   let { "*": catchall } = useParams();
 *   console.log(catchall);
 * }
 * ```
 *
 * @public
 * @category Hooks
 * @returns An object containing the dynamic route parameters
 */
declare function useParams<ParamsOrKey extends string | Record<string, string | undefined> = string>(): Readonly<[
    ParamsOrKey
] extends [string] ? Params<ParamsOrKey> : Partial<ParamsOrKey>>;
/**
 * Resolves the pathname of the given `to` value against the current
 * {@link Location}. Similar to {@link useHref}, but returns a
 * {@link Path} instead of a string.
 *
 * @example
 * import { useResolvedPath } from "react-router";
 *
 * function SomeComponent() {
 *   // if the user is at /dashboard/profile
 *   let path = useResolvedPath("../accounts");
 *   path.pathname; // "/dashboard/accounts"
 *   path.search; // ""
 *   path.hash; // ""
 * }
 *
 * @public
 * @category Hooks
 * @param to The path to resolve
 * @param options Options
 * @param options.relative Defaults to `"route"` so routing is relative to the route tree.
 *                         Set to `"path"` to make relative routing operate against path segments.
 * @returns The resolved {@link Path} object with `pathname`, `search`, and `hash`
 */
declare function useResolvedPath(to: To, { relative }?: {
    relative?: RelativeRoutingType;
}): Path;
/**
 * Hook version of {@link Routes | `<Routes>`} that uses objects instead of
 * components. These objects have the same properties as the component props.
 * The return value of `useRoutes` is either a valid React element you can use
 * to render the route tree, or `null` if nothing matched.
 *
 * @example
 * import { useRoutes } from "react-router";
 *
 * function App() {
 *   let element = useRoutes([
 *     {
 *       path: "/",
 *       element: <Dashboard />,
 *       children: [
 *         {
 *           path: "messages",
 *           element: <DashboardMessages />,
 *         },
 *         { path: "tasks", element: <DashboardTasks /> },
 *       ],
 *     },
 *     { path: "team", element: <AboutPage /> },
 *   ]);
 *
 *   return element;
 * }
 *
 * @public
 * @category Hooks
 * @param routes An array of {@link RouteObject}s that define the route hierarchy
 * @param locationArg An optional {@link Location} object or pathname string to
 * use instead of the current {@link Location}
 * @returns A React element to render the matched route, or `null` if no routes matched
 */
declare function useRoutes(routes: RouteObject[], locationArg?: Partial<Location> | string): React.ReactElement | null;
/**
 * Returns the current {@link Navigation}, defaulting to an "idle" navigation
 * when no navigation is in progress. You can use this to render pending UI
 * (like a global spinner) or read [`FormData`](https://developer.mozilla.org/en-US/docs/Web/API/FormData)
 * from a form navigation.
 *
 * @example
 * import { useNavigation } from "react-router";
 *
 * function SomeComponent() {
 *   let navigation = useNavigation();
 *   navigation.state;
 *   navigation.formData;
 *   // etc.
 * }
 *
 * @public
 * @category Hooks
 * @mode framework
 * @mode data
 * @returns The current {@link Navigation} object
 */
declare function useNavigation(): Navigation;
/**
 * Revalidate the data on the page for reasons outside of normal data mutations
 * like [`Window` focus](https://developer.mozilla.org/en-US/docs/Web/API/Window/focus_event)
 * or polling on an interval.
 *
 * Note that page data is already revalidated automatically after actions.
 * If you find yourself using this for normal CRUD operations on your data in
 * response to user interactions, you're probably not taking advantage of the
 * other APIs like {@link useFetcher}, {@link Form}, {@link useSubmit} that do
 * this automatically.
 *
 * @example
 * import { useRevalidator } from "react-router";
 *
 * function WindowFocusRevalidator() {
 *   const revalidator = useRevalidator();
 *
 *   useFakeWindowFocus(() => {
 *     revalidator.revalidate();
 *   });
 *
 *   return (
 *     <div hidden={revalidator.state === "idle"}>
 *       Revalidating...
 *     </div>
 *   );
 * }
 *
 * @public
 * @category Hooks
 * @mode framework
 * @mode data
 * @returns An object with a `revalidate` function and the current revalidation
 * `state`
 */
declare function useRevalidator(): {
    revalidate: () => Promise<void>;
    state: Router["state"]["revalidation"];
};
/**
 * Returns the active route matches, useful for accessing `loaderData` for
 * parent/child routes or the route [`handle`](../../start/framework/route-module#handle)
 * property
 *
 * @public
 * @category Hooks
 * @mode framework
 * @mode data
 * @returns An array of {@link UIMatch | UI matches} for the current route hierarchy
 */
declare function useMatches(): UIMatch[];
/**
 * Returns the data from the closest route
 * [`loader`](../../start/framework/route-module#loader) or
 * [`clientLoader`](../../start/framework/route-module#clientloader).
 *
 * @example
 * import { useLoaderData } from "react-router";
 *
 * export async function loader() {
 *   return await fakeDb.invoices.findAll();
 * }
 *
 * export default function Invoices() {
 *   let invoices = useLoaderData<typeof loader>();
 *   // ...
 * }
 *
 * @public
 * @category Hooks
 * @mode framework
 * @mode data
 * @returns The data returned from the route's [`loader`](../../start/framework/route-module#loader) or [`clientLoader`](../../start/framework/route-module#clientloader) function
 */
declare function useLoaderData<T = any>(): SerializeFrom<T>;
/**
 * Returns the [`loader`](../../start/framework/route-module#loader) data for a
 * given route by route ID.
 *
 * Route IDs are created automatically. They are simply the path of the route file
 * relative to the app folder without the extension.
 *
 * | Route Filename               | Route ID               |
 * | ---------------------------- | ---------------------- |
 * | `app/root.tsx`               | `"root"`               |
 * | `app/routes/teams.tsx`       | `"routes/teams"`       |
 * | `app/whatever/teams.$id.tsx` | `"whatever/teams.$id"` |
 *
 * @example
 * import { useRouteLoaderData } from "react-router";
 *
 * function SomeComponent() {
 *   const { user } = useRouteLoaderData("root");
 * }
 *
 * // You can also specify your own route ID's manually in your routes.ts file:
 * route("/", "containers/app.tsx", { id: "app" })
 * useRouteLoaderData("app");
 *
 * @public
 * @category Hooks
 * @mode framework
 * @mode data
 * @param routeId The ID of the route to return loader data from
 * @returns The data returned from the specified route's [`loader`](../../start/framework/route-module#loader)
 * function, or `undefined` if not found
 */
declare function useRouteLoaderData<T = any>(routeId: string): SerializeFrom<T> | undefined;
/**
 * Returns the [`action`](../../start/framework/route-module#action) data from
 * the most recent `POST` navigation form submission or `undefined` if there
 * hasn't been one.
 *
 * @example
 * import { Form, useActionData } from "react-router";
 *
 * export async function action({ request }) {
 *   const body = await request.formData();
 *   const name = body.get("visitorsName");
 *   return { message: `Hello, ${name}` };
 * }
 *
 * export default function Invoices() {
 *   const data = useActionData();
 *   return (
 *     <Form method="post">
 *       <input type="text" name="visitorsName" />
 *       {data ? data.message : "Waiting..."}
 *     </Form>
 *   );
 * }
 *
 * @public
 * @category Hooks
 * @mode framework
 * @mode data
 * @returns The data returned from the route's [`action`](../../start/framework/route-module#action)
 * function, or `undefined` if no [`action`](../../start/framework/route-module#action)
 * has been called
 */
declare function useActionData<T = any>(): SerializeFrom<T> | undefined;
/**
 * Accesses the error thrown during an
 * [`action`](../../start/framework/route-module#action),
 * [`loader`](../../start/framework/route-module#loader),
 * or component render to be used in a route module
 * [`ErrorBoundary`](../../start/framework/route-module#errorboundary).
 *
 * @example
 * export function ErrorBoundary() {
 *   const error = useRouteError();
 *   return <div>{error.message}</div>;
 * }
 *
 * @public
 * @category Hooks
 * @mode framework
 * @mode data
 * @returns The error that was thrown during route [loading](../../start/framework/route-module#loader),
 * [`action`](../../start/framework/route-module#action) execution, or rendering
 */
declare function useRouteError(): unknown;
/**
 * Returns the resolved promise value from the closest {@link Await | `<Await>`}.
 *
 * @example
 * function SomeDescendant() {
 *   const value = useAsyncValue();
 *   // ...
 * }
 *
 * // somewhere in your app
 * <Await resolve={somePromise}>
 *   <SomeDescendant />
 * </Await>;
 *
 * @public
 * @category Hooks
 * @mode framework
 * @mode data
 * @returns The resolved value from the nearest {@link Await} component
 */
declare function useAsyncValue(): unknown;
/**
 * Returns the rejection value from the closest {@link Await | `<Await>`}.
 *
 * @example
 * import { Await, useAsyncError } from "react-router";
 *
 * function ErrorElement() {
 *   const error = useAsyncError();
 *   return (
 *     <p>Uh Oh, something went wrong! {error.message}</p>
 *   );
 * }
 *
 * // somewhere in your app
 * <Await
 *   resolve={promiseThatRejects}
 *   errorElement={<ErrorElement />}
 * />;
 *
 * @public
 * @category Hooks
 * @mode framework
 * @mode data
 * @returns The error that was thrown in the nearest {@link Await} component
 */
declare function useAsyncError(): unknown;
/**
 * Allow the application to block navigations within the SPA and present the
 * user a confirmation dialog to confirm the navigation. Mostly used to avoid
 * using half-filled form data. This does not handle hard-reloads or
 * cross-origin navigations.
 *
 * The {@link Blocker} object returned by the hook has the following properties:
 *
 * - **`state`**
 *   - `unblocked` - the blocker is idle and has not prevented any navigation
 *   - `blocked` - the blocker has prevented a navigation
 *   - `proceeding` - the blocker is proceeding through from a blocked navigation
 * - **`location`**
 *   - When in a `blocked` state, this represents the {@link Location} to which we
 *     blocked a navigation. When in a `proceeding` state, this is the location
 *     being navigated to after a `blocker.proceed()` call.
 * - **`proceed()`**
 *   - When in a `blocked` state, you may call `blocker.proceed()` to proceed to the
 *     blocked location.
 * - **`reset()`**
 *   - When in a `blocked` state, you may call `blocker.reset()` to return the blocker
 *     back to an `unblocked` state and leave the user at the current location.
 *
 * @example
 * // Boolean version
 * const blocker = useBlocker(value !== "");
 *
 * // Function version
 * const blocker = useBlocker(
 *   ({ currentLocation, nextLocation, historyAction }) =>
 *     value !== "" &&
 *     currentLocation.pathname !== nextLocation.pathname
 * );
 *
 * @additionalExamples
 * ```tsx
 * import { useCallback, useState } from "react";
 * import { BlockerFunction, useBlocker } from "react-router";
 *
 * export function ImportantForm() {
 *   const [value, setValue] = useState("");
 *
 *   const shouldBlock = useCallback<BlockerFunction>(
 *     () => value !== "",
 *     [value]
 *   );
 *   const blocker = useBlocker(shouldBlock);
 *
 *   return (
 *     <form
 *       onSubmit={(e) => {
 *         e.preventDefault();
 *         setValue("");
 *         if (blocker.state === "blocked") {
 *           blocker.proceed();
 *         }
 *       }}
 *     >
 *       <input
 *         name="data"
 *         value={value}
 *         onChange={(e) => setValue(e.target.value)}
 *       />
 *
 *       <button type="submit">Save</button>
 *
 *       {blocker.state === "blocked" ? (
 *         <>
 *           <p style={{ color: "red" }}>
 *             Blocked the last navigation to
 *           </p>
 *           <button
 *             type="button"
 *             onClick={() => blocker.proceed()}
 *           >
 *             Let me through
 *           </button>
 *           <button
 *             type="button"
 *             onClick={() => blocker.reset()}
 *           >
 *             Keep me here
 *           </button>
 *         </>
 *       ) : blocker.state === "proceeding" ? (
 *         <p style={{ color: "orange" }}>
 *           Proceeding through blocked navigation
 *         </p>
 *       ) : (
 *         <p style={{ color: "green" }}>
 *           Blocker is currently unblocked
 *         </p>
 *       )}
 *     </form>
 *   );
 * }
 * ```
 *
 * @public
 * @category Hooks
 * @mode framework
 * @mode data
 * @param shouldBlock Either a boolean or a function returning a boolean which
 * indicates whether the navigation should be blocked. The function format
 * receives a single object parameter containing the `currentLocation`,
 * `nextLocation`, and `historyAction` of the potential navigation.
 * @returns A {@link Blocker} object with state and reset functionality
 */
declare function useBlocker(shouldBlock: boolean | BlockerFunction): Blocker;

interface ServerRouterProps {
    context: EntryContext;
    url: string | URL;
    nonce?: string;
}
/**
 * The server entry point for a React Router app in Framework Mode. This component
 * is used to generate the HTML in the response from the server.
 * See [`entry.server.tsx`](../api/framework-conventions/entry.server.tsx).
 *
 * @public
 * @category Framework Routers
 * @mode framework
 * @param props Props
 * @param props.context The entry context containing the manifest, route modules,
 * and other data needed for rendering.
 * @param props.url The URL of the request being handled.
 * @param props.nonce An optional nonce for Content Security Policy (CSP) compliance,
 * used to allow inline scripts to run safely.
 * @returns A React element that represents the server-rendered application.
 */
declare function ServerRouter({ context, url, nonce, }: ServerRouterProps): ReactElement;

interface StubRouteExtensions {
    Component?: RouteComponentType;
    HydrateFallback?: HydrateFallbackType;
    ErrorBoundary?: ErrorBoundaryType;
    loader?: LoaderFunction;
    action?: ActionFunction;
    children?: StubRouteObject[];
    meta?: MetaFunction;
    links?: LinksFunction;
}
interface StubIndexRouteObject extends Omit<IndexRouteObject, "Component" | "HydrateFallback" | "ErrorBoundary" | "loader" | "action" | "element" | "errorElement" | "children">, StubRouteExtensions {
}
interface StubNonIndexRouteObject extends Omit<NonIndexRouteObject, "Component" | "HydrateFallback" | "ErrorBoundary" | "loader" | "action" | "element" | "errorElement" | "children">, StubRouteExtensions {
}
type StubRouteObject = StubIndexRouteObject | StubNonIndexRouteObject;
interface RoutesTestStubProps {
    /**
     *  The initial entries in the history stack. This allows you to start a test with
     *  multiple locations already in the history stack (for testing a back navigation, etc.)
     *  The test will default to the last entry in initialEntries if no initialIndex is provided.
     *  e.g. initialEntries={["/home", "/about", "/contact"]}
     */
    initialEntries?: InitialEntry[];
    /**
     * The initial index in the history stack to render. This allows you to start a test at a specific entry.
     * It defaults to the last entry in initialEntries.
     * e.g.
     *   initialEntries: ["/", "/events/123"]
     *   initialIndex: 1 // start at "/events/123"
     */
    initialIndex?: number;
    /**
     *  Used to set the route's initial loader and action data.
     *  e.g. hydrationData={{
     *   loaderData: { "/contact": { locale: "en-US" } },
     *   actionData: { "/login": { errors: { email: "invalid email" } }}
     *  }}
     */
    hydrationData?: HydrationState;
    /**
     * Future flags mimicking the settings in react-router.config.ts
     */
    future?: Partial<FutureConfig>;
}
/**
 * @category Utils
 */
declare function createRoutesStub(routes: StubRouteObject[], _context?: AppLoadContext | unstable_RouterContextProvider): ({ initialEntries, initialIndex, hydrationData, future, }: RoutesTestStubProps) => React.JSX.Element;

interface CookieSignatureOptions {
    /**
     * An array of secrets that may be used to sign/unsign the value of a cookie.
     *
     * The array makes it easy to rotate secrets. New secrets should be added to
     * the beginning of the array. `cookie.serialize()` will always use the first
     * value in the array, but `cookie.parse()` may use any of them so that
     * cookies that were signed with older secrets still work.
     */
    secrets?: string[];
}
type CookieOptions = ParseOptions & SerializeOptions & CookieSignatureOptions;
/**
 * A HTTP cookie.
 *
 * A Cookie is a logical container for metadata about a HTTP cookie; its name
 * and options. But it doesn't contain a value. Instead, it has `parse()` and
 * `serialize()` methods that allow a single instance to be reused for
 * parsing/encoding multiple different values.
 *
 * @see https://remix.run/utils/cookies#cookie-api
 */
interface Cookie {
    /**
     * The name of the cookie, used in the `Cookie` and `Set-Cookie` headers.
     */
    readonly name: string;
    /**
     * True if this cookie uses one or more secrets for verification.
     */
    readonly isSigned: boolean;
    /**
     * The Date this cookie expires.
     *
     * Note: This is calculated at access time using `maxAge` when no `expires`
     * option is provided to `createCookie()`.
     */
    readonly expires?: Date;
    /**
     * Parses a raw `Cookie` header and returns the value of this cookie or
     * `null` if it's not present.
     */
    parse(cookieHeader: string | null, options?: ParseOptions): Promise<any>;
    /**
     * Serializes the given value to a string and returns the `Set-Cookie`
     * header.
     */
    serialize(value: any, options?: SerializeOptions): Promise<string>;
}
/**
 * Creates a logical container for managing a browser cookie from the server.
 */
declare const createCookie: (name: string, cookieOptions?: CookieOptions) => Cookie;
type IsCookieFunction = (object: any) => object is Cookie;
/**
 * Returns true if an object is a Remix cookie container.
 *
 * @see https://remix.run/utils/cookies#iscookie
 */
declare const isCookie: IsCookieFunction;

type RequestHandler = (request: Request, loadContext?: MiddlewareEnabled extends true ? unstable_InitialContext : AppLoadContext) => Promise<Response>;
type CreateRequestHandlerFunction = (build: ServerBuild | (() => ServerBuild | Promise<ServerBuild>), mode?: string) => RequestHandler;
declare const createRequestHandler: CreateRequestHandlerFunction;

/**
 * An object of name/value pairs to be used in the session.
 */
interface SessionData {
    [name: string]: any;
}
/**
 * Session persists data across HTTP requests.
 *
 * @see https://reactrouter.com/explanation/sessions-and-cookies#sessions
 */
interface Session<Data = SessionData, FlashData = Data> {
    /**
     * A unique identifier for this session.
     *
     * Note: This will be the empty string for newly created sessions and
     * sessions that are not backed by a database (i.e. cookie-based sessions).
     */
    readonly id: string;
    /**
     * The raw data contained in this session.
     *
     * This is useful mostly for SessionStorage internally to access the raw
     * session data to persist.
     */
    readonly data: FlashSessionData<Data, FlashData>;
    /**
     * Returns `true` if the session has a value for the given `name`, `false`
     * otherwise.
     */
    has(name: (keyof Data | keyof FlashData) & string): boolean;
    /**
     * Returns the value for the given `name` in this session.
     */
    get<Key extends (keyof Data | keyof FlashData) & string>(name: Key): (Key extends keyof Data ? Data[Key] : undefined) | (Key extends keyof FlashData ? FlashData[Key] : undefined) | undefined;
    /**
     * Sets a value in the session for the given `name`.
     */
    set<Key extends keyof Data & string>(name: Key, value: Data[Key]): void;
    /**
     * Sets a value in the session that is only valid until the next `get()`.
     * This can be useful for temporary values, like error messages.
     */
    flash<Key extends keyof FlashData & string>(name: Key, value: FlashData[Key]): void;
    /**
     * Removes a value from the session.
     */
    unset(name: keyof Data & string): void;
}
type FlashSessionData<Data, FlashData> = Partial<Data & {
    [Key in keyof FlashData as FlashDataKey<Key & string>]: FlashData[Key];
}>;
type FlashDataKey<Key extends string> = `__flash_${Key}__`;
type CreateSessionFunction = <Data = SessionData, FlashData = Data>(initialData?: Data, id?: string) => Session<Data, FlashData>;
/**
 * Creates a new Session object.
 *
 * Note: This function is typically not invoked directly by application code.
 * Instead, use a `SessionStorage` object's `getSession` method.
 */
declare const createSession: CreateSessionFunction;
type IsSessionFunction = (object: any) => object is Session;
/**
 * Returns true if an object is a React Router session.
 *
 * @see https://reactrouter.com/api/utils/isSession
 */
declare const isSession: IsSessionFunction;
/**
 * SessionStorage stores session data between HTTP requests and knows how to
 * parse and create cookies.
 *
 * A SessionStorage creates Session objects using a `Cookie` header as input.
 * Then, later it generates the `Set-Cookie` header to be used in the response.
 */
interface SessionStorage<Data = SessionData, FlashData = Data> {
    /**
     * Parses a Cookie header from a HTTP request and returns the associated
     * Session. If there is no session associated with the cookie, this will
     * return a new Session with no data.
     */
    getSession: (cookieHeader?: string | null, options?: ParseOptions) => Promise<Session<Data, FlashData>>;
    /**
     * Stores all data in the Session and returns the Set-Cookie header to be
     * used in the HTTP response.
     */
    commitSession: (session: Session<Data, FlashData>, options?: SerializeOptions) => Promise<string>;
    /**
     * Deletes all data associated with the Session and returns the Set-Cookie
     * header to be used in the HTTP response.
     */
    destroySession: (session: Session<Data, FlashData>, options?: SerializeOptions) => Promise<string>;
}
/**
 * SessionIdStorageStrategy is designed to allow anyone to easily build their
 * own SessionStorage using `createSessionStorage(strategy)`.
 *
 * This strategy describes a common scenario where the session id is stored in
 * a cookie but the actual session data is stored elsewhere, usually in a
 * database or on disk. A set of create, read, update, and delete operations
 * are provided for managing the session data.
 */
interface SessionIdStorageStrategy<Data = SessionData, FlashData = Data> {
    /**
     * The Cookie used to store the session id, or options used to automatically
     * create one.
     */
    cookie?: Cookie | (CookieOptions & {
        name?: string;
    });
    /**
     * Creates a new record with the given data and returns the session id.
     */
    createData: (data: FlashSessionData<Data, FlashData>, expires?: Date) => Promise<string>;
    /**
     * Returns data for a given session id, or `null` if there isn't any.
     */
    readData: (id: string) => Promise<FlashSessionData<Data, FlashData> | null>;
    /**
     * Updates data for the given session id.
     */
    updateData: (id: string, data: FlashSessionData<Data, FlashData>, expires?: Date) => Promise<void>;
    /**
     * Deletes data for a given session id from the data store.
     */
    deleteData: (id: string) => Promise<void>;
}
/**
 * Creates a SessionStorage object using a SessionIdStorageStrategy.
 *
 * Note: This is a low-level API that should only be used if none of the
 * existing session storage options meet your requirements.
 */
declare function createSessionStorage<Data = SessionData, FlashData = Data>({ cookie: cookieArg, createData, readData, updateData, deleteData, }: SessionIdStorageStrategy<Data, FlashData>): SessionStorage<Data, FlashData>;

interface CookieSessionStorageOptions {
    /**
     * The Cookie used to store the session data on the client, or options used
     * to automatically create one.
     */
    cookie?: SessionIdStorageStrategy["cookie"];
}
/**
 * Creates and returns a SessionStorage object that stores all session data
 * directly in the session cookie itself.
 *
 * This has the advantage that no database or other backend services are
 * needed, and can help to simplify some load-balanced scenarios. However, it
 * also has the limitation that serialized session data may not exceed the
 * browser's maximum cookie size. Trade-offs!
 */
declare function createCookieSessionStorage<Data = SessionData, FlashData = Data>({ cookie: cookieArg }?: CookieSessionStorageOptions): SessionStorage<Data, FlashData>;

interface MemorySessionStorageOptions {
    /**
     * The Cookie used to store the session id on the client, or options used
     * to automatically create one.
     */
    cookie?: SessionIdStorageStrategy["cookie"];
}
/**
 * Creates and returns a simple in-memory SessionStorage object, mostly useful
 * for testing and as a reference implementation.
 *
 * Note: This storage does not scale beyond a single process, so it is not
 * suitable for most production scenarios.
 */
declare function createMemorySessionStorage<Data = SessionData, FlashData = Data>({ cookie }?: MemorySessionStorageOptions): SessionStorage<Data, FlashData>;

type DevServerHooks = {
    getCriticalCss?: (pathname: string) => Promise<string | undefined>;
    processRequestError?: (error: unknown) => void;
};
declare function setDevServerHooks(devServerHooks: DevServerHooks): void;

type Args = {
    [K in keyof Pages]: ToArgs<Pages[K]["params"]>;
};
type ToArgs<Params extends Record<string, string | undefined>> = Equal<Params, {}> extends true ? [] : Partial<Params> extends Params ? [Params] | [] : [
    Params
];
/**
  Returns a resolved URL path for the specified route.

  ```tsx
  const h = href("/:lang?/about", { lang: "en" })
  // -> `/en/about`

  <Link to={href("/products/:id", { id: "abc123" })} />
  ```
 */
declare function href<Path extends keyof Args>(path: Path, ...args: Args[Path]): string;

type BrowserCreateFromReadableStreamFunction = (body: ReadableStream<Uint8Array>, { temporaryReferences, }: {
    temporaryReferences: unknown;
}) => Promise<unknown>;
type EncodeReplyFunction = (args: unknown[], options: {
    temporaryReferences: unknown;
}) => Promise<BodyInit>;
/**
 * Create a React `callServer` implementation for React Router.
 *
 * @example
 * import {
 *   createFromReadableStream,
 *   createTemporaryReferenceSet,
 *   encodeReply,
 *   setServerCallback,
 * } from "@vitejs/plugin-rsc/browser";
 * import { unstable_createCallServer as createCallServer } from "react-router";
 *
 * setServerCallback(
 *   createCallServer({
 *     createFromReadableStream,
 *     createTemporaryReferenceSet,
 *     encodeReply,
 *   })
 * );
 *
 * @name unstable_createCallServer
 * @public
 * @category RSC
 * @mode data
 * @param opts Options
 * @param opts.createFromReadableStream Your `react-server-dom-xyz/client`'s
 * `createFromReadableStream`. Used to decode payloads from the server.
 * @param opts.createTemporaryReferenceSet A function that creates a temporary
 * reference set for the RSC payload.
 * @param opts.encodeReply Your `react-server-dom-xyz/client`'s `encodeReply`.
 * Used when sending payloads to the server.
 * @param opts.fetch Optional fetch implementation. Defaults to global `fetch`.
 * @returns A function that can be used to call server actions.
 */
declare function createCallServer({ createFromReadableStream, createTemporaryReferenceSet, encodeReply, fetch: fetchImplementation, }: {
    createFromReadableStream: BrowserCreateFromReadableStreamFunction;
    createTemporaryReferenceSet: () => unknown;
    encodeReply: EncodeReplyFunction;
    fetch?: (request: Request) => Promise<Response>;
}): (id: string, args: unknown[]) => Promise<unknown>;
/**
 * Props for the `unstable_RSCHydratedRouter` component.
 * @name unstable_RSCHydratedRouterProps
 */
interface RSCHydratedRouterProps {
    /**
     * Your `react-server-dom-xyz/client`'s `createFromReadableStream` function,
     * used to decode payloads from the server.
     */
    createFromReadableStream: BrowserCreateFromReadableStreamFunction;
    /**
     * Optional fetch implementation.  Defaults to global `fetch`.
     */
    fetch?: (request: Request) => Promise<Response>;
    /**
     * The decoded `RSCPayload` to hydrate.
     */
    payload: RSCPayload;
    /**
     * `eager` or `lazy` - Determines if links are eagerly discovered, or delayed
     * until clicked.
     */
    routeDiscovery?: "eager" | "lazy";
    /**
     * A function that returns an `unstable_InitialContext` object
     * (`Map<RouterContext, unknown>`), for use in client loaders, actions and
     * middleware.
     */
    unstable_getContext?: RouterInit["unstable_getContext"];
}
/**
 * Hydrates a server rendered `RSCPayload` in the browser.
 *
 * @example
 * import { startTransition, StrictMode } from "react";
 * import { hydrateRoot } from "react-dom/client";
 * import {
 *   unstable_getRSCStream as getRSCStream,
 *   unstable_RSCHydratedRouter as RSCHydratedRouter,
 * } from "react-router";
 * import type { unstable_RSCPayload as RSCPayload } from "react-router";
 *
 * createFromReadableStream(getRSCStream()).then(
 *   (payload: RSCServerPayload) => {
 *     startTransition(async () => {
 *       hydrateRoot(
 *         document,
 *         <StrictMode>
 *           <RSCHydratedRouter
 *             createFromReadableStream={
 *               createFromReadableStream
 *             }
 *             payload={payload}
 *           />
 *         </StrictMode>,
 *         {
 *           formState: await getFormState(payload),
 *         }
 *       );
 *     });
 *   }
 * );
 *
 * @name unstable_RSCHydratedRouter
 * @public
 * @category RSC
 * @mode data
 * @param props Props
 * @param {unstable_RSCHydratedRouterProps.createFromReadableStream} props.createFromReadableStream n/a
 * @param {unstable_RSCHydratedRouterProps.fetch} props.fetch n/a
 * @param {unstable_RSCHydratedRouterProps.payload} props.payload n/a
 * @param {unstable_RSCHydratedRouterProps.routeDiscovery} props.routeDiscovery n/a
 * @param {unstable_RSCHydratedRouterProps.unstable_getContext} props.unstable_getContext n/a
 * @returns A hydrated router that can be used to navigate and render routes.
 */
declare function RSCHydratedRouter({ createFromReadableStream, fetch: fetchImplementation, payload, routeDiscovery, unstable_getContext, }: RSCHydratedRouterProps): React.JSX.Element;

type SSRCreateFromReadableStreamFunction = (body: ReadableStream<Uint8Array>) => Promise<unknown>;
/**
 * Routes the incoming request to the RSC server and appropriately proxies the
 * server response for data / resource requests, or renders to HTML for a document
 * request.
 *
 * @example
 * import { createFromReadableStream } from "@vitejs/plugin-rsc/ssr";
 * import * as ReactDomServer from "react-dom/server.edge";
 * import {
 *   unstable_RSCStaticRouter as RSCStaticRouter,
 *   unstable_routeRSCServerRequest as routeRSCServerRequest,
 * } from "react-router";
 *
 * routeRSCServerRequest({
 *   request,
 *   fetchServer,
 *   createFromReadableStream,
 *   async renderHTML(getPayload) {
 *     const payload = await getPayload();
 *
 *     return await renderHTMLToReadableStream(
 *       <RSCStaticRouter getPayload={getPayload} />,
 *       {
 *         bootstrapScriptContent,
 *         formState: await getFormState(payload),
 *       }
 *     );
 *   },
 * });
 *
 * @name unstable_routeRSCServerRequest
 * @public
 * @category RSC
 * @mode data
 * @param opts Options
 * @param opts.createFromReadableStream Your `react-server-dom-xyz/client`'s
 * `createFromReadableStream` function, used to decode payloads from the server.
 * @param opts.fetchServer A function that forwards a `Request` to the RSC handler
 * and returns a `Promise<Response>` containing a serialized `RSCPayload`.
 * @param opts.renderHTML A function that renders the `RSCPayload` to HTML, usually
 * using a `<RSCStaticRouter>`.
 * @param opts.request The request to route.
 * @param opts.hydrate Whether to hydrate the server response with the RSC payload.
 * Defaults to `true`.
 * @returns A `Response` that either contains the RSC payload for data requests, or
 * renders the HTML for document requests.
 */
declare function routeRSCServerRequest({ request, fetchServer, createFromReadableStream, renderHTML, hydrate, }: {
    request: Request;
    fetchServer: (request: Request) => Promise<Response>;
    createFromReadableStream: SSRCreateFromReadableStreamFunction;
    renderHTML: (getPayload: () => Promise<RSCPayload>) => ReadableStream<Uint8Array> | Promise<ReadableStream<Uint8Array>>;
    hydrate?: boolean;
}): Promise<Response>;
/**
 * Props for the `RSCStaticRouter` component.
 * @name unstable_RSCStaticRouterProps
 */
interface RSCStaticRouterProps {
    /**
     * A function that starts decoding of the `RSCPayload`. Usually passed through
     * from `routeRSCServerRequest`'s `renderHTML`.
     */
    getPayload: () => Promise<RSCPayload>;
}
/**
 * Pre-renders an `RSCPayload` to HTML. Usually used in `routeRSCServerRequest`'s
 * `renderHTML` callback.
 *
 * @example
 * import { createFromReadableStream } from "@vitejs/plugin-rsc/ssr";
 * import * as ReactDomServer from "react-dom/server.edge";
 * import {
 *   unstable_RSCStaticRouter as RSCStaticRouter,
 *   unstable_routeRSCServerRequest as routeRSCServerRequest,
 * } from "react-router";
 *
 * routeRSCServerRequest({
 *   request,
 *   fetchServer,
 *   createFromReadableStream,
 *   async renderHTML(getPayload) {
 *     const payload = await getPayload();
 *
 *     return await renderHTMLToReadableStream(
 *       <RSCStaticRouter getPayload={getPayload} />,
 *       {
 *         bootstrapScriptContent,
 *         formState: await getFormState(payload),
 *       }
 *     );
 *   },
 * });
 *
 * @name unstable_RSCStaticRouter
 * @public
 * @category RSC
 * @mode data
 * @param props Props
 * @param {unstable_RSCStaticRouterProps.getPayload} props.getPayload n/a
 * @returns A React component that renders the `RSCPayload` as HTML.
 */
declare function RSCStaticRouter({ getPayload }: RSCStaticRouterProps): React.JSX.Element | null;

declare function RSCDefaultRootErrorBoundary({ hasRootLayout, }: {
    hasRootLayout: boolean;
}): React__default.JSX.Element;

declare function deserializeErrors(errors: RouterState["errors"]): RouterState["errors"];

type RemixErrorBoundaryProps = React.PropsWithChildren<{
    location: Location;
    isOutsideRemixApp?: boolean;
    error?: Error;
}>;
type RemixErrorBoundaryState = {
    error: null | Error;
    location: Location;
};
declare class RemixErrorBoundary extends React.Component<RemixErrorBoundaryProps, RemixErrorBoundaryState> {
    constructor(props: RemixErrorBoundaryProps);
    static getDerivedStateFromError(error: Error): {
        error: Error;
    };
    static getDerivedStateFromProps(props: RemixErrorBoundaryProps, state: RemixErrorBoundaryState): {
        error: Error | null;
        location: Location<any>;
    };
    render(): string | number | boolean | Iterable<React.ReactNode> | React.JSX.Element | null | undefined;
}

declare function getPatchRoutesOnNavigationFunction(manifest: AssetsManifest, routeModules: RouteModules, ssr: boolean, routeDiscovery: ServerBuild["routeDiscovery"], isSpaMode: boolean, basename: string | undefined): PatchRoutesOnNavigationFunction | undefined;
declare function useFogOFWarDiscovery(router: Router, manifest: AssetsManifest, routeModules: RouteModules, ssr: boolean, routeDiscovery: ServerBuild["routeDiscovery"], isSpaMode: boolean): void;

declare function getHydrationData(state: {
    loaderData?: Router["state"]["loaderData"];
    actionData?: Router["state"]["actionData"];
    errors?: Router["state"]["errors"];
}, routes: DataRouteObject[], getRouteInfo: (routeId: string) => {
    clientLoader: ClientLoaderFunction | undefined;
    hasLoader: boolean;
    hasHydrateFallback: boolean;
}, location: Path, basename: string | undefined, isSpaMode: boolean): HydrationState;

declare const unstable_matchRSCServerRequest: typeof matchRSCServerRequest;

export { ActionFunction, AppLoadContext, Blocker, BlockerFunction, ClientLoaderFunction, type Cookie, type CookieOptions, type CookieSignatureOptions, type CreateRequestHandlerFunction, DataRouteObject, Router as DataRouter, DataStrategyFunction, EntryContext, type FlashSessionData, HydrationState, IndexRouteObject, InitialEntry, type IsCookieFunction, type IsSessionFunction, LinksFunction, LoaderFunction, Location, MetaFunction, type NavigateFunction, NavigateOptions, Navigation, Action as NavigationType, NonIndexRouteObject, ParamParseKey, Params, PatchRoutesOnNavigationFunction, Path, PathMatch, PathPattern, RelativeRoutingType, type RequestHandler, RouteObject, RouterInit, RouterState, type RoutesTestStubProps, ServerBuild, ServerRouter, type ServerRouterProps, type Session, type SessionData, type SessionIdStorageStrategy, type SessionStorage, To, UIMatch, AssetsManifest as UNSAFE_AssetsManifest, MiddlewareEnabled as UNSAFE_MiddlewareEnabled, RSCDefaultRootErrorBoundary as UNSAFE_RSCDefaultRootErrorBoundary, RemixErrorBoundary as UNSAFE_RemixErrorBoundary, RouteModules as UNSAFE_RouteModules, ServerMode as UNSAFE_ServerMode, SingleFetchRedirectSymbol as UNSAFE_SingleFetchRedirectSymbol, decodeViaTurboStream as UNSAFE_decodeViaTurboStream, deserializeErrors as UNSAFE_deserializeErrors, getHydrationData as UNSAFE_getHydrationData, getPatchRoutesOnNavigationFunction as UNSAFE_getPatchRoutesOnNavigationFunction, getTurboStreamSingleFetchDataStrategy as UNSAFE_getTurboStreamSingleFetchDataStrategy, useFogOFWarDiscovery as UNSAFE_useFogOFWarDiscovery, createCookie, createCookieSessionStorage, createMemorySessionStorage, createRequestHandler, createRoutesStub, createSession, createSessionStorage, href, isCookie, isSession, type BrowserCreateFromReadableStreamFunction as unstable_BrowserCreateFromReadableStreamFunction, type EncodeReplyFunction as unstable_EncodeReplyFunction, unstable_InitialContext, RSCHydratedRouter as unstable_RSCHydratedRouter, type RSCHydratedRouterProps as unstable_RSCHydratedRouterProps, RSCPayload as unstable_RSCPayload, RSCStaticRouter as unstable_RSCStaticRouter, type RSCStaticRouterProps as unstable_RSCStaticRouterProps, unstable_RouterContextProvider, type SSRCreateFromReadableStreamFunction as unstable_SSRCreateFromReadableStreamFunction, createCallServer as unstable_createCallServer, unstable_matchRSCServerRequest, routeRSCServerRequest as unstable_routeRSCServerRequest, setDevServerHooks as unstable_setDevServerHooks, useActionData, useAsyncError, useAsyncValue, useBlocker, useHref, useInRouterContext, useLoaderData, useLocation, useMatch, useMatches, useNavigate, useNavigation, useNavigationType, useOutlet, useOutletContext, useParams, useResolvedPath, useRevalidator, useRouteError, useRouteLoaderData, useRoutes };
