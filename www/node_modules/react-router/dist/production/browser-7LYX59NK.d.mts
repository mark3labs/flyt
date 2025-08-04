import { AsyncLocalStorage } from 'node:async_hooks';
import * as React from 'react';
import { u as unstable_RouterContextProvider, i as ActionFunction, C as ClientActionFunction, j as ClientLoaderFunction, k as HeadersFunction, l as LinksFunction, m as LoaderFunction, M as MetaFunction, S as ShouldRevalidateFunction, c as Location, h as Params } from './route-data-CqEmXQub.mjs';

type ServerContext = {
    redirect?: Response;
};
declare global {
    var ___reactRouterServerStorage___: AsyncLocalStorage<ServerContext> | undefined;
}
type RSCRouteConfigEntryBase = {
    action?: ActionFunction;
    clientAction?: ClientActionFunction;
    clientLoader?: ClientLoaderFunction;
    ErrorBoundary?: React.ComponentType<any>;
    handle?: any;
    headers?: HeadersFunction;
    HydrateFallback?: React.ComponentType<any>;
    Layout?: React.ComponentType<any>;
    links?: LinksFunction;
    loader?: LoaderFunction;
    meta?: MetaFunction;
    shouldRevalidate?: ShouldRevalidateFunction;
};
type RSCRouteConfigEntry = RSCRouteConfigEntryBase & {
    id: string;
    path?: string;
    Component?: React.ComponentType<any>;
    lazy?: () => Promise<RSCRouteConfigEntryBase & ({
        default?: React.ComponentType<any>;
        Component?: never;
    } | {
        default?: never;
        Component?: React.ComponentType<any>;
    })>;
} & ({
    index: true;
} | {
    children?: RSCRouteConfigEntry[];
});
type RSCRouteConfig = Array<RSCRouteConfigEntry>;
type RSCRouteManifest = {
    clientAction?: ClientActionFunction;
    clientLoader?: ClientLoaderFunction;
    element?: React.ReactElement | false;
    errorElement?: React.ReactElement;
    handle?: any;
    hasAction: boolean;
    hasComponent: boolean;
    hasErrorBoundary: boolean;
    hasLoader: boolean;
    hydrateFallbackElement?: React.ReactElement;
    id: string;
    index?: boolean;
    links?: LinksFunction;
    meta?: MetaFunction;
    parentId?: string;
    path?: string;
    shouldRevalidate?: ShouldRevalidateFunction;
};
type RSCRouteMatch = RSCRouteManifest & {
    params: Params;
    pathname: string;
    pathnameBase: string;
};
type RSCRenderPayload = {
    type: "render";
    actionData: Record<string, any> | null;
    basename: string | undefined;
    errors: Record<string, any> | null;
    loaderData: Record<string, any>;
    location: Location;
    matches: RSCRouteMatch[];
    patches?: RSCRouteManifest[];
    nonce?: string;
    formState?: unknown;
};
type RSCManifestPayload = {
    type: "manifest";
    patches: RSCRouteManifest[];
};
type RSCActionPayload = {
    type: "action";
    actionResult: Promise<unknown>;
    rerender?: Promise<RSCRenderPayload | RSCRedirectPayload>;
};
type RSCRedirectPayload = {
    type: "redirect";
    status: number;
    location: string;
    replace: boolean;
    reload: boolean;
    actionResult?: Promise<unknown>;
};
type RSCPayload = RSCRenderPayload | RSCManifestPayload | RSCActionPayload | RSCRedirectPayload;
type RSCMatch = {
    statusCode: number;
    headers: Headers;
    payload: RSCPayload;
};
type DecodeActionFunction = (formData: FormData) => Promise<() => Promise<unknown>>;
type DecodeFormStateFunction = (result: unknown, formData: FormData) => unknown;
type DecodeReplyFunction = (reply: FormData | string, { temporaryReferences }: {
    temporaryReferences: unknown;
}) => Promise<unknown[]>;
type LoadServerActionFunction = (id: string) => Promise<Function>;
/**
 * Matches the given routes to a Request and returns a RSC Response encoding an
 * `RSCPayload` for consumption by a RSC enabled client router.
 *
 * @example
 * import {
 *   createTemporaryReferenceSet,
 *   decodeAction,
 *   decodeReply,
 *   loadServerAction,
 *   renderToReadableStream,
 * } from "@vitejs/plugin-rsc/rsc";
 * import { unstable_matchRSCServerRequest as matchRSCServerRequest } from "react-router";
 *
 * matchRSCServerRequest({
 *   createTemporaryReferenceSet,
 *   decodeAction,
 *   decodeFormState,
 *   decodeReply,
 *   loadServerAction,
 *   request,
 *   routes: routes(),
 *   generateResponse(match) {
 *     return new Response(
 *       renderToReadableStream(match.payload),
 *       {
 *         status: match.statusCode,
 *         headers: match.headers,
 *       }
 *     );
 *   },
 * });
 *
 * @name unstable_matchRSCServerRequest
 * @public
 * @category RSC
 * @mode data
 * @param opts Options
 * @param opts.basename The basename to use when matching the request.
 * @param opts.decodeAction Your `react-server-dom-xyz/server`'s `decodeAction`
 * function, responsible for loading a server action.
 * @param opts.decodeReply Your `react-server-dom-xyz/server`'s `decodeReply`
 * function, used to decode the server function's arguments and bind them to the
 * implementation for invocation by the router.
 * @param opts.decodeFormState A function responsible for decoding form state for
 * progressively enhanceable forms with `useActionState` using your
 * `react-server-dom-xyz/server`'s `decodeFormState`.
 * @param opts.generateResponse A function responsible for using your
 * `renderToReadableStream` to generate a Response encoding the `RSCPayload`.
 * @param opts.loadServerAction Your `react-server-dom-xyz/server`'s
 * `loadServerAction` function, used to load a server action by ID.
 * @param opts.request The request to match against.
 * @param opts.requestContext An instance of `unstable_RouterContextProvider`
 * that should be created per request, to be passed to loaders, actions and middleware.
 * @param opts.routes Your route definitions.
 * @param opts.createTemporaryReferenceSet A function that returns a temporary
 * reference set for the request, used to track temporary references in the RSC stream.
 * @param opts.onError An optional error handler that will be called with any
 * errors that occur during the request processing.
 * @returns A Response that contains the RSC data for hydration.
 */
declare function matchRSCServerRequest({ createTemporaryReferenceSet, basename, decodeReply, requestContext, loadServerAction, decodeAction, decodeFormState, onError, request, routes, generateResponse, }: {
    createTemporaryReferenceSet: () => unknown;
    basename?: string;
    decodeReply?: DecodeReplyFunction;
    decodeAction?: DecodeActionFunction;
    decodeFormState?: DecodeFormStateFunction;
    requestContext?: unstable_RouterContextProvider;
    loadServerAction?: LoadServerActionFunction;
    onError?: (error: unknown) => void;
    request: Request;
    routes: RSCRouteConfigEntry[];
    generateResponse: (match: RSCMatch, { temporaryReferences, }: {
        temporaryReferences: unknown;
    }) => Response;
}): Promise<Response>;

declare global {
    interface Window {
        __FLIGHT_DATA: any[];
    }
}
/**
 * Get the prerendered RSC stream for hydration. Usually passed directly to your
 * `react-server-dom-xyz/client`'s `createFromReadableStream`.
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
 *           <RSCHydratedRouter ...props />
 *         </StrictMode>,
 *         {
 *           // Options
 *         }
 *       );
 *     });
 *   }
 * );
 *
 * @name unstable_getRSCStream
 * @public
 * @category RSC
 * @mode data
 * @returns A `ReadableStream` that contains the RSC data for hydration.
 */
declare function getRSCStream(): ReadableStream<any>;

export { type DecodeActionFunction as D, type LoadServerActionFunction as L, type RSCPayload as R, type DecodeFormStateFunction as a, type DecodeReplyFunction as b, type RSCManifestPayload as c, type RSCMatch as d, type RSCRenderPayload as e, type RSCRouteManifest as f, getRSCStream as g, type RSCRouteMatch as h, type RSCRouteConfigEntry as i, type RSCRouteConfig as j, matchRSCServerRequest as m };
