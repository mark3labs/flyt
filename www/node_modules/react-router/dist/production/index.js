"use strict";Object.defineProperty(exports, "__esModule", {value: true}); function _interopRequireWildcard(obj) { if (obj && obj.__esModule) { return obj; } else { var newObj = {}; if (obj != null) { for (var key in obj) { if (Object.prototype.hasOwnProperty.call(obj, key)) { newObj[key] = obj[key]; } } } newObj.default = obj; return newObj; } } function _nullishCoalesce(lhs, rhsFn) { if (lhs != null) { return lhs; } else { return rhsFn(); } } function _optionalChain(ops) { let lastAccessLHS = undefined; let value = ops[0]; let i = 1; while (i < ops.length) { const op = ops[i]; const fn = ops[i + 1]; i += 2; if ((op === 'optionalAccess' || op === 'optionalCall') && value == null) { return undefined; } if (op === 'access' || op === 'optionalAccess') { lastAccessLHS = value; value = fn(value); } else if (op === 'call' || op === 'optionalCall') { value = fn((...args) => value.call(lastAccessLHS, ...args)); lastAccessLHS = undefined; } } return value; }/**
 * react-router v7.7.1
 *
 * Copyright (c) Remix Software Inc.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE.md file in the root directory of this source tree.
 *
 * @license MIT
 */
"use client";













































var _chunk4DGLNKXFjs = require('./chunk-4DGLNKXF.js');























































































var _chunk7OQROU2Djs = require('./chunk-7OQROU2D.js');

// lib/dom/ssr/server.tsx
var _react = require('react'); var React = _interopRequireWildcard(_react); var React2 = _interopRequireWildcard(_react); var React4 = _interopRequireWildcard(_react); var React5 = _interopRequireWildcard(_react);
function ServerRouter({
  context,
  url,
  nonce
}) {
  if (typeof url === "string") {
    url = new URL(url);
  }
  let { manifest, routeModules, criticalCss, serverHandoffString } = context;
  let routes = _chunk7OQROU2Djs.createServerRoutes.call(void 0, 
    manifest.routes,
    routeModules,
    context.future,
    context.isSpaMode
  );
  context.staticHandlerContext.loaderData = {
    ...context.staticHandlerContext.loaderData
  };
  for (let match of context.staticHandlerContext.matches) {
    let routeId = match.route.id;
    let route = routeModules[routeId];
    let manifestRoute = context.manifest.routes[routeId];
    if (route && manifestRoute && _chunk7OQROU2Djs.shouldHydrateRouteLoader.call(void 0, 
      routeId,
      route.clientLoader,
      manifestRoute.hasLoader,
      context.isSpaMode
    ) && (route.HydrateFallback || !manifestRoute.hasLoader)) {
      delete context.staticHandlerContext.loaderData[routeId];
    }
  }
  let router = _chunk4DGLNKXFjs.createStaticRouter.call(void 0, routes, context.staticHandlerContext);
  return /* @__PURE__ */ React.createElement(React.Fragment, null, /* @__PURE__ */ React.createElement(
    _chunk7OQROU2Djs.FrameworkContext.Provider,
    {
      value: {
        manifest,
        routeModules,
        criticalCss,
        serverHandoffString,
        future: context.future,
        ssr: context.ssr,
        isSpaMode: context.isSpaMode,
        routeDiscovery: context.routeDiscovery,
        serializeError: context.serializeError,
        renderMeta: context.renderMeta
      }
    },
    /* @__PURE__ */ React.createElement(_chunk7OQROU2Djs.RemixErrorBoundary, { location: router.state.location }, /* @__PURE__ */ React.createElement(
      _chunk4DGLNKXFjs.StaticRouterProvider,
      {
        router,
        context: context.staticHandlerContext,
        hydrate: false
      }
    ))
  ), context.serverHandoffStream ? /* @__PURE__ */ React.createElement(React.Suspense, null, /* @__PURE__ */ React.createElement(
    _chunk7OQROU2Djs.StreamTransfer,
    {
      context,
      identifier: 0,
      reader: context.serverHandoffStream.getReader(),
      textDecoder: new TextDecoder(),
      nonce
    }
  )) : null);
}

// lib/dom/ssr/routes-test-stub.tsx

function createRoutesStub(routes, _context) {
  return function RoutesTestStub({
    initialEntries,
    initialIndex,
    hydrationData,
    future
  }) {
    let routerRef = React2.useRef();
    let frameworkContextRef = React2.useRef();
    if (routerRef.current == null) {
      frameworkContextRef.current = {
        future: {
          unstable_subResourceIntegrity: _optionalChain([future, 'optionalAccess', _2 => _2.unstable_subResourceIntegrity]) === true,
          unstable_middleware: _optionalChain([future, 'optionalAccess', _3 => _3.unstable_middleware]) === true
        },
        manifest: {
          routes: {},
          entry: { imports: [], module: "" },
          url: "",
          version: ""
        },
        routeModules: {},
        ssr: false,
        isSpaMode: false,
        routeDiscovery: { mode: "lazy", manifestPath: "/__manifest" }
      };
      let patched = processRoutes(
        // @ts-expect-error `StubRouteObject` is stricter about `loader`/`action`
        // types compared to `AgnosticRouteObject`
        _chunk7OQROU2Djs.convertRoutesToDataRoutes.call(void 0, routes, (r) => r),
        _context !== void 0 ? _context : _optionalChain([future, 'optionalAccess', _4 => _4.unstable_middleware]) ? new (0, _chunk7OQROU2Djs.unstable_RouterContextProvider)() : {},
        frameworkContextRef.current.manifest,
        frameworkContextRef.current.routeModules
      );
      routerRef.current = _chunk4DGLNKXFjs.createMemoryRouter.call(void 0, patched, {
        initialEntries,
        initialIndex,
        hydrationData
      });
    }
    return /* @__PURE__ */ React2.createElement(_chunk7OQROU2Djs.FrameworkContext.Provider, { value: frameworkContextRef.current }, /* @__PURE__ */ React2.createElement(_chunk4DGLNKXFjs.RouterProvider, { router: routerRef.current }));
  };
}
function processRoutes(routes, context, manifest, routeModules, parentId) {
  return routes.map((route) => {
    if (!route.id) {
      throw new Error(
        "Expected a route.id in react-router processRoutes() function"
      );
    }
    let newRoute = {
      id: route.id,
      path: route.path,
      index: route.index,
      Component: route.Component ? _chunk4DGLNKXFjs.withComponentProps.call(void 0, route.Component) : void 0,
      HydrateFallback: route.HydrateFallback ? _chunk4DGLNKXFjs.withHydrateFallbackProps.call(void 0, route.HydrateFallback) : void 0,
      ErrorBoundary: route.ErrorBoundary ? _chunk4DGLNKXFjs.withErrorBoundaryProps.call(void 0, route.ErrorBoundary) : void 0,
      action: route.action ? (args) => route.action({ ...args, context }) : void 0,
      loader: route.loader ? (args) => route.loader({ ...args, context }) : void 0,
      handle: route.handle,
      shouldRevalidate: route.shouldRevalidate
    };
    let entryRoute = {
      id: route.id,
      path: route.path,
      index: route.index,
      parentId,
      hasAction: route.action != null,
      hasLoader: route.loader != null,
      // When testing routes, you should be stubbing loader/action/middleware,
      // not trying to re-implement the full loader/clientLoader/SSR/hydration
      // flow. That is better tested via E2E tests.
      hasClientAction: false,
      hasClientLoader: false,
      hasClientMiddleware: false,
      hasErrorBoundary: route.ErrorBoundary != null,
      // any need for these?
      module: "build/stub-path-to-module.js",
      clientActionModule: void 0,
      clientLoaderModule: void 0,
      clientMiddlewareModule: void 0,
      hydrateFallbackModule: void 0
    };
    manifest.routes[newRoute.id] = entryRoute;
    routeModules[route.id] = {
      default: newRoute.Component || _chunk4DGLNKXFjs.Outlet,
      ErrorBoundary: newRoute.ErrorBoundary || void 0,
      handle: route.handle,
      links: route.links,
      meta: route.meta,
      shouldRevalidate: route.shouldRevalidate
    };
    if (route.children) {
      newRoute.children = processRoutes(
        route.children,
        context,
        manifest,
        routeModules,
        newRoute.id
      );
    }
    return newRoute;
  });
}

// lib/server-runtime/cookies.ts
var _cookie = require('cookie');

// lib/server-runtime/crypto.ts
var encoder = /* @__PURE__ */ new TextEncoder();
var sign = async (value, secret) => {
  let data2 = encoder.encode(value);
  let key = await createKey(secret, ["sign"]);
  let signature = await crypto.subtle.sign("HMAC", key, data2);
  let hash = btoa(String.fromCharCode(...new Uint8Array(signature))).replace(
    /=+$/,
    ""
  );
  return value + "." + hash;
};
var unsign = async (cookie, secret) => {
  let index = cookie.lastIndexOf(".");
  let value = cookie.slice(0, index);
  let hash = cookie.slice(index + 1);
  let data2 = encoder.encode(value);
  let key = await createKey(secret, ["verify"]);
  try {
    let signature = byteStringToUint8Array(atob(hash));
    let valid = await crypto.subtle.verify("HMAC", key, signature, data2);
    return valid ? value : false;
  } catch (error) {
    return false;
  }
};
var createKey = async (secret, usages) => crypto.subtle.importKey(
  "raw",
  encoder.encode(secret),
  { name: "HMAC", hash: "SHA-256" },
  false,
  usages
);
function byteStringToUint8Array(byteString) {
  let array = new Uint8Array(byteString.length);
  for (let i = 0; i < byteString.length; i++) {
    array[i] = byteString.charCodeAt(i);
  }
  return array;
}

// lib/server-runtime/cookies.ts
var createCookie = (name, cookieOptions = {}) => {
  let { secrets = [], ...options } = {
    path: "/",
    sameSite: "lax",
    ...cookieOptions
  };
  warnOnceAboutExpiresCookie(name, options.expires);
  return {
    get name() {
      return name;
    },
    get isSigned() {
      return secrets.length > 0;
    },
    get expires() {
      return typeof options.maxAge !== "undefined" ? new Date(Date.now() + options.maxAge * 1e3) : options.expires;
    },
    async parse(cookieHeader, parseOptions) {
      if (!cookieHeader) return null;
      let cookies = _cookie.parse.call(void 0, cookieHeader, { ...options, ...parseOptions });
      if (name in cookies) {
        let value = cookies[name];
        if (typeof value === "string" && value !== "") {
          let decoded = await decodeCookieValue(value, secrets);
          return decoded;
        } else {
          return "";
        }
      } else {
        return null;
      }
    },
    async serialize(value, serializeOptions) {
      return _cookie.serialize.call(void 0, 
        name,
        value === "" ? "" : await encodeCookieValue(value, secrets),
        {
          ...options,
          ...serializeOptions
        }
      );
    }
  };
};
var isCookie = (object) => {
  return object != null && typeof object.name === "string" && typeof object.isSigned === "boolean" && typeof object.parse === "function" && typeof object.serialize === "function";
};
async function encodeCookieValue(value, secrets) {
  let encoded = encodeData(value);
  if (secrets.length > 0) {
    encoded = await sign(encoded, secrets[0]);
  }
  return encoded;
}
async function decodeCookieValue(value, secrets) {
  if (secrets.length > 0) {
    for (let secret of secrets) {
      let unsignedValue = await unsign(value, secret);
      if (unsignedValue !== false) {
        return decodeData(unsignedValue);
      }
    }
    return null;
  }
  return decodeData(value);
}
function encodeData(value) {
  return btoa(myUnescape(encodeURIComponent(JSON.stringify(value))));
}
function decodeData(value) {
  try {
    return JSON.parse(decodeURIComponent(myEscape(atob(value))));
  } catch (error) {
    return {};
  }
}
function myEscape(value) {
  let str = value.toString();
  let result = "";
  let index = 0;
  let chr, code;
  while (index < str.length) {
    chr = str.charAt(index++);
    if (/[\w*+\-./@]/.exec(chr)) {
      result += chr;
    } else {
      code = chr.charCodeAt(0);
      if (code < 256) {
        result += "%" + hex(code, 2);
      } else {
        result += "%u" + hex(code, 4).toUpperCase();
      }
    }
  }
  return result;
}
function hex(code, length) {
  let result = code.toString(16);
  while (result.length < length) result = "0" + result;
  return result;
}
function myUnescape(value) {
  let str = value.toString();
  let result = "";
  let index = 0;
  let chr, part;
  while (index < str.length) {
    chr = str.charAt(index++);
    if (chr === "%") {
      if (str.charAt(index) === "u") {
        part = str.slice(index + 1, index + 5);
        if (/^[\da-f]{4}$/i.exec(part)) {
          result += String.fromCharCode(parseInt(part, 16));
          index += 5;
          continue;
        }
      } else {
        part = str.slice(index, index + 2);
        if (/^[\da-f]{2}$/i.exec(part)) {
          result += String.fromCharCode(parseInt(part, 16));
          index += 2;
          continue;
        }
      }
    }
    result += chr;
  }
  return result;
}
function warnOnceAboutExpiresCookie(name, expires) {
  _chunk7OQROU2Djs.warnOnce.call(void 0, 
    !expires,
    `The "${name}" cookie has an "expires" property set. This will cause the expires value to not be updated when the session is committed. Instead, you should set the expires value when serializing the cookie. You can use \`commitSession(session, { expires })\` if using a session storage object, or \`cookie.serialize("value", { expires })\` if you're using the cookie directly.`
  );
}

// lib/server-runtime/entry.ts
function createEntryRouteModules(manifest) {
  return Object.keys(manifest).reduce((memo, routeId) => {
    let route = manifest[routeId];
    if (route) {
      memo[routeId] = route.module;
    }
    return memo;
  }, {});
}

// lib/server-runtime/mode.ts
var ServerMode = /* @__PURE__ */ ((ServerMode2) => {
  ServerMode2["Development"] = "development";
  ServerMode2["Production"] = "production";
  ServerMode2["Test"] = "test";
  return ServerMode2;
})(ServerMode || {});
function isServerMode(value) {
  return value === "development" /* Development */ || value === "production" /* Production */ || value === "test" /* Test */;
}

// lib/server-runtime/errors.ts
function sanitizeError(error, serverMode) {
  if (error instanceof Error && serverMode !== "development" /* Development */) {
    let sanitized = new Error("Unexpected Server Error");
    sanitized.stack = void 0;
    return sanitized;
  }
  return error;
}
function sanitizeErrors(errors, serverMode) {
  return Object.entries(errors).reduce((acc, [routeId, error]) => {
    return Object.assign(acc, { [routeId]: sanitizeError(error, serverMode) });
  }, {});
}
function serializeError(error, serverMode) {
  let sanitized = sanitizeError(error, serverMode);
  return {
    message: sanitized.message,
    stack: sanitized.stack
  };
}
function serializeErrors(errors, serverMode) {
  if (!errors) return null;
  let entries = Object.entries(errors);
  let serialized = {};
  for (let [key, val] of entries) {
    if (_chunk7OQROU2Djs.isRouteErrorResponse.call(void 0, val)) {
      serialized[key] = { ...val, __type: "RouteErrorResponse" };
    } else if (val instanceof Error) {
      let sanitized = sanitizeError(val, serverMode);
      serialized[key] = {
        message: sanitized.message,
        stack: sanitized.stack,
        __type: "Error",
        // If this is a subclass (i.e., ReferenceError), send up the type so we
        // can re-create the same type during hydration.  This will only apply
        // in dev mode since all production errors are sanitized to normal
        // Error instances
        ...sanitized.name !== "Error" ? {
          __subType: sanitized.name
        } : {}
      };
    } else {
      serialized[key] = val;
    }
  }
  return serialized;
}

// lib/server-runtime/routeMatching.ts
function matchServerRoutes(routes, pathname, basename) {
  let matches = _chunk7OQROU2Djs.matchRoutes.call(void 0, 
    routes,
    pathname,
    basename
  );
  if (!matches) return null;
  return matches.map((match) => ({
    params: match.params,
    pathname: match.pathname,
    route: match.route
  }));
}

// lib/server-runtime/data.ts
async function callRouteHandler(handler, args) {
  let result = await handler({
    request: stripRoutesParam(stripIndexParam2(args.request)),
    params: args.params,
    context: args.context
  });
  if (_chunk7OQROU2Djs.isDataWithResponseInit.call(void 0, result) && result.init && result.init.status && _chunk7OQROU2Djs.isRedirectStatusCode.call(void 0, result.init.status)) {
    throw new Response(null, result.init);
  }
  return result;
}
function stripIndexParam2(request) {
  let url = new URL(request.url);
  let indexValues = url.searchParams.getAll("index");
  url.searchParams.delete("index");
  let indexValuesToKeep = [];
  for (let indexValue of indexValues) {
    if (indexValue) {
      indexValuesToKeep.push(indexValue);
    }
  }
  for (let toKeep of indexValuesToKeep) {
    url.searchParams.append("index", toKeep);
  }
  let init = {
    method: request.method,
    body: request.body,
    headers: request.headers,
    signal: request.signal
  };
  if (init.body) {
    init.duplex = "half";
  }
  return new Request(url.href, init);
}
function stripRoutesParam(request) {
  let url = new URL(request.url);
  url.searchParams.delete("_routes");
  let init = {
    method: request.method,
    body: request.body,
    headers: request.headers,
    signal: request.signal
  };
  if (init.body) {
    init.duplex = "half";
  }
  return new Request(url.href, init);
}

// lib/server-runtime/invariant.ts
function invariant2(value, message) {
  if (value === false || value === null || typeof value === "undefined") {
    console.error(
      "The following error is a bug in React Router; please open an issue! https://github.com/remix-run/react-router/issues/new/choose"
    );
    throw new Error(message);
  }
}

// lib/server-runtime/dev.ts
var globalDevServerHooksKey = "__reactRouterDevServerHooks";
function setDevServerHooks(devServerHooks) {
  globalThis[globalDevServerHooksKey] = devServerHooks;
}
function getDevServerHooks() {
  return globalThis[globalDevServerHooksKey];
}
function getBuildTimeHeader(request, headerName) {
  if (typeof process !== "undefined") {
    try {
      if (_optionalChain([process, 'access', _5 => _5.env, 'optionalAccess', _6 => _6.IS_RR_BUILD_REQUEST]) === "yes") {
        return request.headers.get(headerName);
      }
    } catch (e) {
    }
  }
  return null;
}

// lib/server-runtime/routes.ts
function groupRoutesByParentId(manifest) {
  let routes = {};
  Object.values(manifest).forEach((route) => {
    if (route) {
      let parentId = route.parentId || "";
      if (!routes[parentId]) {
        routes[parentId] = [];
      }
      routes[parentId].push(route);
    }
  });
  return routes;
}
function createRoutes(manifest, parentId = "", routesByParentId = groupRoutesByParentId(manifest)) {
  return (routesByParentId[parentId] || []).map((route) => ({
    ...route,
    children: createRoutes(manifest, route.id, routesByParentId)
  }));
}
function createStaticHandlerDataRoutes(manifest, future, parentId = "", routesByParentId = groupRoutesByParentId(manifest)) {
  return (routesByParentId[parentId] || []).map((route) => {
    let commonRoute = {
      // Always include root due to default boundaries
      hasErrorBoundary: route.id === "root" || route.module.ErrorBoundary != null,
      id: route.id,
      path: route.path,
      unstable_middleware: route.module.unstable_middleware,
      // Need to use RR's version in the param typed here to permit the optional
      // context even though we know it'll always be provided in remix
      loader: route.module.loader ? async (args) => {
        let preRenderedData = getBuildTimeHeader(
          args.request,
          "X-React-Router-Prerender-Data"
        );
        if (preRenderedData != null) {
          let encoded = preRenderedData ? decodeURI(preRenderedData) : preRenderedData;
          invariant2(encoded, "Missing prerendered data for route");
          let uint8array = new TextEncoder().encode(encoded);
          let stream = new ReadableStream({
            start(controller) {
              controller.enqueue(uint8array);
              controller.close();
            }
          });
          let decoded = await _chunk7OQROU2Djs.decodeViaTurboStream.call(void 0, stream, global);
          let data2 = decoded.value;
          if (data2 && _chunk7OQROU2Djs.SingleFetchRedirectSymbol in data2) {
            let result = data2[_chunk7OQROU2Djs.SingleFetchRedirectSymbol];
            let init = { status: result.status };
            if (result.reload) {
              throw _chunk7OQROU2Djs.redirectDocument.call(void 0, result.redirect, init);
            } else if (result.replace) {
              throw _chunk7OQROU2Djs.replace.call(void 0, result.redirect, init);
            } else {
              throw _chunk7OQROU2Djs.redirect.call(void 0, result.redirect, init);
            }
          } else {
            invariant2(
              data2 && route.id in data2,
              "Unable to decode prerendered data"
            );
            let result = data2[route.id];
            invariant2(
              "data" in result,
              "Unable to process prerendered data"
            );
            return result.data;
          }
        }
        let val = await callRouteHandler(route.module.loader, args);
        return val;
      } : void 0,
      action: route.module.action ? (args) => callRouteHandler(route.module.action, args) : void 0,
      handle: route.module.handle
    };
    return route.index ? {
      index: true,
      ...commonRoute
    } : {
      caseSensitive: route.caseSensitive,
      children: createStaticHandlerDataRoutes(
        manifest,
        future,
        route.id,
        routesByParentId
      ),
      ...commonRoute
    };
  });
}

// lib/server-runtime/markup.ts
var ESCAPE_LOOKUP = {
  "&": "\\u0026",
  ">": "\\u003e",
  "<": "\\u003c",
  "\u2028": "\\u2028",
  "\u2029": "\\u2029"
};
var ESCAPE_REGEX = /[&><\u2028\u2029]/g;
function escapeHtml(html) {
  return html.replace(ESCAPE_REGEX, (match) => ESCAPE_LOOKUP[match]);
}

// lib/server-runtime/serverHandoff.ts
function createServerHandoffString(serverHandoff) {
  return escapeHtml(JSON.stringify(serverHandoff));
}

// lib/server-runtime/headers.ts
var _setcookieparser = require('set-cookie-parser');
function getDocumentHeaders(context, build) {
  return getDocumentHeadersImpl(context, (m) => {
    let route = build.routes[m.route.id];
    invariant2(route, `Route with id "${m.route.id}" not found in build`);
    return route.module.headers;
  });
}
function getDocumentHeadersImpl(context, getRouteHeadersFn) {
  let boundaryIdx = context.errors ? context.matches.findIndex((m) => context.errors[m.route.id]) : -1;
  let matches = boundaryIdx >= 0 ? context.matches.slice(0, boundaryIdx + 1) : context.matches;
  let errorHeaders;
  if (boundaryIdx >= 0) {
    let { actionHeaders, actionData, loaderHeaders, loaderData } = context;
    context.matches.slice(boundaryIdx).some((match) => {
      let id = match.route.id;
      if (actionHeaders[id] && (!actionData || !actionData.hasOwnProperty(id))) {
        errorHeaders = actionHeaders[id];
      } else if (loaderHeaders[id] && !loaderData.hasOwnProperty(id)) {
        errorHeaders = loaderHeaders[id];
      }
      return errorHeaders != null;
    });
  }
  return matches.reduce((parentHeaders, match, idx) => {
    let { id } = match.route;
    let loaderHeaders = context.loaderHeaders[id] || new Headers();
    let actionHeaders = context.actionHeaders[id] || new Headers();
    let includeErrorHeaders = errorHeaders != null && idx === matches.length - 1;
    let includeErrorCookies = includeErrorHeaders && errorHeaders !== loaderHeaders && errorHeaders !== actionHeaders;
    let headersFn = getRouteHeadersFn(match);
    if (headersFn == null) {
      let headers2 = new Headers(parentHeaders);
      if (includeErrorCookies) {
        prependCookies(errorHeaders, headers2);
      }
      prependCookies(actionHeaders, headers2);
      prependCookies(loaderHeaders, headers2);
      return headers2;
    }
    let headers = new Headers(
      typeof headersFn === "function" ? headersFn({
        loaderHeaders,
        parentHeaders,
        actionHeaders,
        errorHeaders: includeErrorHeaders ? errorHeaders : void 0
      }) : headersFn
    );
    if (includeErrorCookies) {
      prependCookies(errorHeaders, headers);
    }
    prependCookies(actionHeaders, headers);
    prependCookies(loaderHeaders, headers);
    prependCookies(parentHeaders, headers);
    return headers;
  }, new Headers());
}
function prependCookies(parentHeaders, childHeaders) {
  let parentSetCookieString = parentHeaders.get("Set-Cookie");
  if (parentSetCookieString) {
    let cookies = _setcookieparser.splitCookiesString.call(void 0, parentSetCookieString);
    let childCookies = new Set(childHeaders.getSetCookie());
    cookies.forEach((cookie) => {
      if (!childCookies.has(cookie)) {
        childHeaders.append("Set-Cookie", cookie);
      }
    });
  }
}

// lib/server-runtime/single-fetch.ts
var SERVER_NO_BODY_STATUS_CODES = /* @__PURE__ */ new Set([
  ..._chunk7OQROU2Djs.NO_BODY_STATUS_CODES,
  304
]);
async function singleFetchAction(build, serverMode, staticHandler, request, handlerUrl, loadContext, handleError) {
  try {
    let respond2 = function(context) {
      let headers = getDocumentHeaders(context, build);
      if (_chunk7OQROU2Djs.isRedirectStatusCode.call(void 0, context.statusCode) && headers.has("Location")) {
        return generateSingleFetchResponse(request, build, serverMode, {
          result: getSingleFetchRedirect(
            context.statusCode,
            headers,
            build.basename
          ),
          headers,
          status: _chunk7OQROU2Djs.SINGLE_FETCH_REDIRECT_STATUS
        });
      }
      if (context.errors) {
        Object.values(context.errors).forEach((err) => {
          if (!_chunk7OQROU2Djs.isRouteErrorResponse.call(void 0, err) || err.error) {
            handleError(err);
          }
        });
        context.errors = sanitizeErrors(context.errors, serverMode);
      }
      let singleFetchResult;
      if (context.errors) {
        singleFetchResult = { error: Object.values(context.errors)[0] };
      } else {
        singleFetchResult = {
          data: Object.values(context.actionData || {})[0]
        };
      }
      return generateSingleFetchResponse(request, build, serverMode, {
        result: singleFetchResult,
        headers,
        status: context.statusCode
      });
    };
    var respond = respond2;
    let handlerRequest = new Request(handlerUrl, {
      method: request.method,
      body: request.body,
      headers: request.headers,
      signal: request.signal,
      ...request.body ? { duplex: "half" } : void 0
    });
    let result = await staticHandler.query(handlerRequest, {
      requestContext: loadContext,
      skipLoaderErrorBubbling: true,
      skipRevalidation: true,
      unstable_respond: respond2
    });
    if (!_chunk7OQROU2Djs.isResponse.call(void 0, result)) {
      result = respond2(result);
    }
    if (_chunk7OQROU2Djs.isRedirectResponse.call(void 0, result)) {
      return generateSingleFetchResponse(request, build, serverMode, {
        result: getSingleFetchRedirect(
          result.status,
          result.headers,
          build.basename
        ),
        headers: result.headers,
        status: _chunk7OQROU2Djs.SINGLE_FETCH_REDIRECT_STATUS
      });
    }
    return result;
  } catch (error) {
    handleError(error);
    return generateSingleFetchResponse(request, build, serverMode, {
      result: { error },
      headers: new Headers(),
      status: 500
    });
  }
}
async function singleFetchLoaders(build, serverMode, staticHandler, request, handlerUrl, loadContext, handleError) {
  try {
    let respond2 = function(context) {
      let headers = getDocumentHeaders(context, build);
      if (_chunk7OQROU2Djs.isRedirectStatusCode.call(void 0, context.statusCode) && headers.has("Location")) {
        return generateSingleFetchResponse(request, build, serverMode, {
          result: {
            [_chunk7OQROU2Djs.SingleFetchRedirectSymbol]: getSingleFetchRedirect(
              context.statusCode,
              headers,
              build.basename
            )
          },
          headers,
          status: _chunk7OQROU2Djs.SINGLE_FETCH_REDIRECT_STATUS
        });
      }
      if (context.errors) {
        Object.values(context.errors).forEach((err) => {
          if (!_chunk7OQROU2Djs.isRouteErrorResponse.call(void 0, err) || err.error) {
            handleError(err);
          }
        });
        context.errors = sanitizeErrors(context.errors, serverMode);
      }
      let results = {};
      let loadedMatches = new Set(
        context.matches.filter(
          (m) => loadRouteIds ? loadRouteIds.has(m.route.id) : m.route.loader != null
        ).map((m) => m.route.id)
      );
      if (context.errors) {
        for (let [id, error] of Object.entries(context.errors)) {
          results[id] = { error };
        }
      }
      for (let [id, data2] of Object.entries(context.loaderData)) {
        if (!(id in results) && loadedMatches.has(id)) {
          results[id] = { data: data2 };
        }
      }
      return generateSingleFetchResponse(request, build, serverMode, {
        result: results,
        headers,
        status: context.statusCode
      });
    };
    var respond = respond2;
    let handlerRequest = new Request(handlerUrl, {
      headers: request.headers,
      signal: request.signal
    });
    let routesParam = new URL(request.url).searchParams.get("_routes");
    let loadRouteIds = routesParam ? new Set(routesParam.split(",")) : null;
    let result = await staticHandler.query(handlerRequest, {
      requestContext: loadContext,
      filterMatchesToLoad: (m) => !loadRouteIds || loadRouteIds.has(m.route.id),
      skipLoaderErrorBubbling: true,
      unstable_respond: respond2
    });
    if (!_chunk7OQROU2Djs.isResponse.call(void 0, result)) {
      result = respond2(result);
    }
    if (_chunk7OQROU2Djs.isRedirectResponse.call(void 0, result)) {
      return generateSingleFetchResponse(request, build, serverMode, {
        result: {
          [_chunk7OQROU2Djs.SingleFetchRedirectSymbol]: getSingleFetchRedirect(
            result.status,
            result.headers,
            build.basename
          )
        },
        headers: result.headers,
        status: _chunk7OQROU2Djs.SINGLE_FETCH_REDIRECT_STATUS
      });
    }
    return result;
  } catch (error) {
    handleError(error);
    return generateSingleFetchResponse(request, build, serverMode, {
      result: { root: { error } },
      headers: new Headers(),
      status: 500
    });
  }
}
function generateSingleFetchResponse(request, build, serverMode, {
  result,
  headers,
  status
}) {
  let resultHeaders = new Headers(headers);
  resultHeaders.set("X-Remix-Response", "yes");
  if (SERVER_NO_BODY_STATUS_CODES.has(status)) {
    return new Response(null, { status, headers: resultHeaders });
  }
  resultHeaders.set("Content-Type", "text/x-script");
  resultHeaders.delete("Content-Length");
  return new Response(
    encodeViaTurboStream(
      result,
      request.signal,
      build.entry.module.streamTimeout,
      serverMode
    ),
    {
      status: status || 200,
      headers: resultHeaders
    }
  );
}
function getSingleFetchRedirect(status, headers, basename) {
  let redirect2 = headers.get("Location");
  if (basename) {
    redirect2 = _chunk7OQROU2Djs.stripBasename.call(void 0, redirect2, basename) || redirect2;
  }
  return {
    redirect: redirect2,
    status,
    revalidate: (
      // Technically X-Remix-Revalidate isn't needed here - that was an implementation
      // detail of ?_data requests as our way to tell the front end to revalidate when
      // we didn't have a response body to include that information in.
      // With single fetch, we tell the front end via this revalidate boolean field.
      // However, we're respecting it for now because it may be something folks have
      // used in their own responses
      // TODO(v3): Consider removing or making this official public API
      headers.has("X-Remix-Revalidate") || headers.has("Set-Cookie")
    ),
    reload: headers.has("X-Remix-Reload-Document"),
    replace: headers.has("X-Remix-Replace")
  };
}
function encodeViaTurboStream(data2, requestSignal, streamTimeout, serverMode) {
  let controller = new AbortController();
  let timeoutId = setTimeout(
    () => controller.abort(new Error("Server Timeout")),
    typeof streamTimeout === "number" ? streamTimeout : 4950
  );
  requestSignal.addEventListener("abort", () => clearTimeout(timeoutId));
  return _chunk7OQROU2Djs.encode.call(void 0, data2, {
    signal: controller.signal,
    plugins: [
      (value) => {
        if (value instanceof Error) {
          let { name, message, stack } = serverMode === "production" /* Production */ ? sanitizeError(value, serverMode) : value;
          return ["SanitizedError", name, message, stack];
        }
        if (value instanceof _chunk7OQROU2Djs.ErrorResponseImpl) {
          let { data: data3, status, statusText } = value;
          return ["ErrorResponse", data3, status, statusText];
        }
        if (value && typeof value === "object" && _chunk7OQROU2Djs.SingleFetchRedirectSymbol in value) {
          return ["SingleFetchRedirect", value[_chunk7OQROU2Djs.SingleFetchRedirectSymbol]];
        }
      }
    ],
    postPlugins: [
      (value) => {
        if (!value) return;
        if (typeof value !== "object") return;
        return [
          "SingleFetchClassInstance",
          Object.fromEntries(Object.entries(value))
        ];
      },
      () => ["SingleFetchFallback"]
    ]
  });
}

// lib/server-runtime/server.ts
function derive(build, mode) {
  let routes = createRoutes(build.routes);
  let dataRoutes = createStaticHandlerDataRoutes(build.routes, build.future);
  let serverMode = isServerMode(mode) ? mode : "production" /* Production */;
  let staticHandler = _chunk7OQROU2Djs.createStaticHandler.call(void 0, dataRoutes, {
    basename: build.basename
  });
  let errorHandler = build.entry.module.handleError || ((error, { request }) => {
    if (serverMode !== "test" /* Test */ && !request.signal.aborted) {
      console.error(
        // @ts-expect-error This is "private" from users but intended for internal use
        _chunk7OQROU2Djs.isRouteErrorResponse.call(void 0, error) && error.error ? error.error : error
      );
    }
  });
  return {
    routes,
    dataRoutes,
    serverMode,
    staticHandler,
    errorHandler
  };
}
var createRequestHandler = (build, mode) => {
  let _build;
  let routes;
  let serverMode;
  let staticHandler;
  let errorHandler;
  return async function requestHandler(request, initialContext) {
    _build = typeof build === "function" ? await build() : build;
    if (typeof build === "function") {
      let derived = derive(_build, mode);
      routes = derived.routes;
      serverMode = derived.serverMode;
      staticHandler = derived.staticHandler;
      errorHandler = derived.errorHandler;
    } else if (!routes || !serverMode || !staticHandler || !errorHandler) {
      let derived = derive(_build, mode);
      routes = derived.routes;
      serverMode = derived.serverMode;
      staticHandler = derived.staticHandler;
      errorHandler = derived.errorHandler;
    }
    let params = {};
    let loadContext;
    let handleError = (error) => {
      if (mode === "development" /* Development */) {
        _optionalChain([getDevServerHooks, 'call', _7 => _7(), 'optionalAccess', _8 => _8.processRequestError, 'optionalCall', _9 => _9(error)]);
      }
      errorHandler(error, {
        context: loadContext,
        params,
        request
      });
    };
    if (_build.future.unstable_middleware) {
      if (initialContext == null) {
        loadContext = new (0, _chunk7OQROU2Djs.unstable_RouterContextProvider)();
      } else {
        try {
          loadContext = new (0, _chunk7OQROU2Djs.unstable_RouterContextProvider)(
            initialContext
          );
        } catch (e) {
          let error = new Error(
            `Unable to create initial \`unstable_RouterContextProvider\` instance. Please confirm you are returning an instance of \`Map<unstable_routerContext, unknown>\` from your \`getLoadContext\` function.

Error: ${e instanceof Error ? e.toString() : e}`
          );
          handleError(error);
          return returnLastResortErrorResponse(error, serverMode);
        }
      }
    } else {
      loadContext = initialContext || {};
    }
    let url = new URL(request.url);
    let normalizedBasename = _build.basename || "/";
    let normalizedPath = url.pathname;
    if (_chunk7OQROU2Djs.stripBasename.call(void 0, normalizedPath, normalizedBasename) === "/_root.data") {
      normalizedPath = normalizedBasename;
    } else if (normalizedPath.endsWith(".data")) {
      normalizedPath = normalizedPath.replace(/\.data$/, "");
    }
    if (_chunk7OQROU2Djs.stripBasename.call(void 0, normalizedPath, normalizedBasename) !== "/" && normalizedPath.endsWith("/")) {
      normalizedPath = normalizedPath.slice(0, -1);
    }
    let isSpaMode = getBuildTimeHeader(request, "X-React-Router-SPA-Mode") === "yes";
    if (!_build.ssr) {
      let decodedPath = decodeURI(normalizedPath);
      if (_build.prerender.length === 0) {
        isSpaMode = true;
      } else if (!_build.prerender.includes(decodedPath) && !_build.prerender.includes(decodedPath + "/")) {
        if (url.pathname.endsWith(".data")) {
          errorHandler(
            new (0, _chunk7OQROU2Djs.ErrorResponseImpl)(
              404,
              "Not Found",
              `Refusing to SSR the path \`${decodedPath}\` because \`ssr:false\` is set and the path is not included in the \`prerender\` config, so in production the path will be a 404.`
            ),
            {
              context: loadContext,
              params,
              request
            }
          );
          return new Response("Not Found", {
            status: 404,
            statusText: "Not Found"
          });
        } else {
          isSpaMode = true;
        }
      }
    }
    let manifestUrl = _chunk7OQROU2Djs.getManifestPath.call(void 0, 
      _build.routeDiscovery.manifestPath,
      normalizedBasename
    );
    if (url.pathname === manifestUrl) {
      try {
        let res = await handleManifestRequest(_build, routes, url);
        return res;
      } catch (e) {
        handleError(e);
        return new Response("Unknown Server Error", { status: 500 });
      }
    }
    let matches = matchServerRoutes(routes, normalizedPath, _build.basename);
    if (matches && matches.length > 0) {
      Object.assign(params, matches[0].params);
    }
    let response;
    if (url.pathname.endsWith(".data")) {
      let handlerUrl = new URL(request.url);
      handlerUrl.pathname = normalizedPath;
      let singleFetchMatches = matchServerRoutes(
        routes,
        handlerUrl.pathname,
        _build.basename
      );
      response = await handleSingleFetchRequest(
        serverMode,
        _build,
        staticHandler,
        request,
        handlerUrl,
        loadContext,
        handleError
      );
      if (_build.entry.module.handleDataRequest) {
        response = await _build.entry.module.handleDataRequest(response, {
          context: loadContext,
          params: singleFetchMatches ? singleFetchMatches[0].params : {},
          request
        });
        if (_chunk7OQROU2Djs.isRedirectResponse.call(void 0, response)) {
          let result = getSingleFetchRedirect(
            response.status,
            response.headers,
            _build.basename
          );
          if (request.method === "GET") {
            result = {
              [_chunk7OQROU2Djs.SingleFetchRedirectSymbol]: result
            };
          }
          let headers = new Headers(response.headers);
          headers.set("Content-Type", "text/x-script");
          return new Response(
            encodeViaTurboStream(
              result,
              request.signal,
              _build.entry.module.streamTimeout,
              serverMode
            ),
            {
              status: _chunk7OQROU2Djs.SINGLE_FETCH_REDIRECT_STATUS,
              headers
            }
          );
        }
      }
    } else if (!isSpaMode && matches && matches[matches.length - 1].route.module.default == null && matches[matches.length - 1].route.module.ErrorBoundary == null) {
      response = await handleResourceRequest(
        serverMode,
        _build,
        staticHandler,
        matches.slice(-1)[0].route.id,
        request,
        loadContext,
        handleError
      );
    } else {
      let { pathname } = url;
      let criticalCss = void 0;
      if (_build.unstable_getCriticalCss) {
        criticalCss = await _build.unstable_getCriticalCss({ pathname });
      } else if (mode === "development" /* Development */ && _optionalChain([getDevServerHooks, 'call', _10 => _10(), 'optionalAccess', _11 => _11.getCriticalCss])) {
        criticalCss = await _optionalChain([getDevServerHooks, 'call', _12 => _12(), 'optionalAccess', _13 => _13.getCriticalCss, 'optionalCall', _14 => _14(pathname)]);
      }
      response = await handleDocumentRequest(
        serverMode,
        _build,
        staticHandler,
        request,
        loadContext,
        handleError,
        isSpaMode,
        criticalCss
      );
    }
    if (request.method === "HEAD") {
      return new Response(null, {
        headers: response.headers,
        status: response.status,
        statusText: response.statusText
      });
    }
    return response;
  };
};
async function handleManifestRequest(build, routes, url) {
  if (build.assets.version !== url.searchParams.get("version")) {
    return new Response(null, {
      status: 204,
      headers: {
        "X-Remix-Reload-Document": "true"
      }
    });
  }
  let patches = {};
  if (url.searchParams.has("p")) {
    let paths = /* @__PURE__ */ new Set();
    url.searchParams.getAll("p").forEach((path) => {
      if (!path.startsWith("/")) {
        path = `/${path}`;
      }
      let segments = path.split("/").slice(1);
      segments.forEach((_, i) => {
        let partialPath = segments.slice(0, i + 1).join("/");
        paths.add(`/${partialPath}`);
      });
    });
    for (let path of paths) {
      let matches = matchServerRoutes(routes, path, build.basename);
      if (matches) {
        for (let match of matches) {
          let routeId = match.route.id;
          let route = build.assets.routes[routeId];
          if (route) {
            patches[routeId] = route;
          }
        }
      }
    }
    return Response.json(patches, {
      headers: {
        "Cache-Control": "public, max-age=31536000, immutable"
      }
    });
  }
  return new Response("Invalid Request", { status: 400 });
}
async function handleSingleFetchRequest(serverMode, build, staticHandler, request, handlerUrl, loadContext, handleError) {
  let response = request.method !== "GET" ? await singleFetchAction(
    build,
    serverMode,
    staticHandler,
    request,
    handlerUrl,
    loadContext,
    handleError
  ) : await singleFetchLoaders(
    build,
    serverMode,
    staticHandler,
    request,
    handlerUrl,
    loadContext,
    handleError
  );
  return response;
}
async function handleDocumentRequest(serverMode, build, staticHandler, request, loadContext, handleError, isSpaMode, criticalCss) {
  try {
    let response = await staticHandler.query(request, {
      requestContext: loadContext,
      unstable_respond: build.future.unstable_middleware ? (ctx) => renderHtml(ctx, isSpaMode) : void 0
    });
    return _chunk7OQROU2Djs.isResponse.call(void 0, response) ? response : renderHtml(response, isSpaMode);
  } catch (error) {
    handleError(error);
    return new Response(null, { status: 500 });
  }
  async function renderHtml(context, isSpaMode2) {
    if (_chunk7OQROU2Djs.isResponse.call(void 0, context)) {
      return context;
    }
    let headers = getDocumentHeaders(context, build);
    if (SERVER_NO_BODY_STATUS_CODES.has(context.statusCode)) {
      return new Response(null, { status: context.statusCode, headers });
    }
    if (context.errors) {
      Object.values(context.errors).forEach((err) => {
        if (!_chunk7OQROU2Djs.isRouteErrorResponse.call(void 0, err) || err.error) {
          handleError(err);
        }
      });
      context.errors = sanitizeErrors(context.errors, serverMode);
    }
    let state = {
      loaderData: context.loaderData,
      actionData: context.actionData,
      errors: serializeErrors(context.errors, serverMode)
    };
    let baseServerHandoff = {
      basename: build.basename,
      future: build.future,
      routeDiscovery: build.routeDiscovery,
      ssr: build.ssr,
      isSpaMode: isSpaMode2
    };
    let entryContext = {
      manifest: build.assets,
      routeModules: createEntryRouteModules(build.routes),
      staticHandlerContext: context,
      criticalCss,
      serverHandoffString: createServerHandoffString({
        ...baseServerHandoff,
        criticalCss
      }),
      serverHandoffStream: encodeViaTurboStream(
        state,
        request.signal,
        build.entry.module.streamTimeout,
        serverMode
      ),
      renderMeta: {},
      future: build.future,
      ssr: build.ssr,
      routeDiscovery: build.routeDiscovery,
      isSpaMode: isSpaMode2,
      serializeError: (err) => serializeError(err, serverMode)
    };
    let handleDocumentRequestFunction = build.entry.module.default;
    try {
      return await handleDocumentRequestFunction(
        request,
        context.statusCode,
        headers,
        entryContext,
        loadContext
      );
    } catch (error) {
      handleError(error);
      let errorForSecondRender = error;
      if (_chunk7OQROU2Djs.isResponse.call(void 0, error)) {
        try {
          let data2 = await unwrapResponse(error);
          errorForSecondRender = new (0, _chunk7OQROU2Djs.ErrorResponseImpl)(
            error.status,
            error.statusText,
            data2
          );
        } catch (e) {
        }
      }
      context = _chunk7OQROU2Djs.getStaticContextFromError.call(void 0, 
        staticHandler.dataRoutes,
        context,
        errorForSecondRender
      );
      if (context.errors) {
        context.errors = sanitizeErrors(context.errors, serverMode);
      }
      let state2 = {
        loaderData: context.loaderData,
        actionData: context.actionData,
        errors: serializeErrors(context.errors, serverMode)
      };
      entryContext = {
        ...entryContext,
        staticHandlerContext: context,
        serverHandoffString: createServerHandoffString(baseServerHandoff),
        serverHandoffStream: encodeViaTurboStream(
          state2,
          request.signal,
          build.entry.module.streamTimeout,
          serverMode
        ),
        renderMeta: {}
      };
      try {
        return await handleDocumentRequestFunction(
          request,
          context.statusCode,
          headers,
          entryContext,
          loadContext
        );
      } catch (error2) {
        handleError(error2);
        return returnLastResortErrorResponse(error2, serverMode);
      }
    }
  }
}
async function handleResourceRequest(serverMode, build, staticHandler, routeId, request, loadContext, handleError) {
  try {
    let response = await staticHandler.queryRoute(request, {
      routeId,
      requestContext: loadContext,
      unstable_respond: build.future.unstable_middleware ? (ctx) => ctx : void 0
    });
    if (_chunk7OQROU2Djs.isResponse.call(void 0, response)) {
      return response;
    }
    if (typeof response === "string") {
      return new Response(response);
    }
    return Response.json(response);
  } catch (error) {
    if (_chunk7OQROU2Djs.isResponse.call(void 0, error)) {
      error.headers.set("X-Remix-Catch", "yes");
      return error;
    }
    if (_chunk7OQROU2Djs.isRouteErrorResponse.call(void 0, error)) {
      if (error) {
        handleError(error);
      }
      return errorResponseToJson(error, serverMode);
    }
    if (error instanceof Error && error.message === "Expected a response from queryRoute") {
      let newError = new Error(
        "Expected a Response to be returned from resource route handler"
      );
      handleError(newError);
      return returnLastResortErrorResponse(newError, serverMode);
    }
    handleError(error);
    return returnLastResortErrorResponse(error, serverMode);
  }
}
function errorResponseToJson(errorResponse, serverMode) {
  return Response.json(
    serializeError(
      // @ts-expect-error This is "private" from users but intended for internal use
      errorResponse.error || new Error("Unexpected Server Error"),
      serverMode
    ),
    {
      status: errorResponse.status,
      statusText: errorResponse.statusText,
      headers: {
        "X-Remix-Error": "yes"
      }
    }
  );
}
function returnLastResortErrorResponse(error, serverMode) {
  let message = "Unexpected Server Error";
  if (serverMode !== "production" /* Production */) {
    message += `

${String(error)}`;
  }
  return new Response(message, {
    status: 500,
    headers: {
      "Content-Type": "text/plain"
    }
  });
}
function unwrapResponse(response) {
  let contentType = response.headers.get("Content-Type");
  return contentType && /\bapplication\/json\b/.test(contentType) ? response.body == null ? null : response.json() : response.text();
}

// lib/server-runtime/sessions.ts
function flash(name) {
  return `__flash_${name}__`;
}
var createSession = (initialData = {}, id = "") => {
  let map = new Map(Object.entries(initialData));
  return {
    get id() {
      return id;
    },
    get data() {
      return Object.fromEntries(map);
    },
    has(name) {
      return map.has(name) || map.has(flash(name));
    },
    get(name) {
      if (map.has(name)) return map.get(name);
      let flashName = flash(name);
      if (map.has(flashName)) {
        let value = map.get(flashName);
        map.delete(flashName);
        return value;
      }
      return void 0;
    },
    set(name, value) {
      map.set(name, value);
    },
    flash(name, value) {
      map.set(flash(name), value);
    },
    unset(name) {
      map.delete(name);
    }
  };
};
var isSession = (object) => {
  return object != null && typeof object.id === "string" && typeof object.data !== "undefined" && typeof object.has === "function" && typeof object.get === "function" && typeof object.set === "function" && typeof object.flash === "function" && typeof object.unset === "function";
};
function createSessionStorage({
  cookie: cookieArg,
  createData,
  readData,
  updateData,
  deleteData
}) {
  let cookie = isCookie(cookieArg) ? cookieArg : createCookie(_optionalChain([cookieArg, 'optionalAccess', _15 => _15.name]) || "__session", cookieArg);
  warnOnceAboutSigningSessionCookie(cookie);
  return {
    async getSession(cookieHeader, options) {
      let id = cookieHeader && await cookie.parse(cookieHeader, options);
      let data2 = id && await readData(id);
      return createSession(data2 || {}, id || "");
    },
    async commitSession(session, options) {
      let { id, data: data2 } = session;
      let expires = _optionalChain([options, 'optionalAccess', _16 => _16.maxAge]) != null ? new Date(Date.now() + options.maxAge * 1e3) : _optionalChain([options, 'optionalAccess', _17 => _17.expires]) != null ? options.expires : cookie.expires;
      if (id) {
        await updateData(id, data2, expires);
      } else {
        id = await createData(data2, expires);
      }
      return cookie.serialize(id, options);
    },
    async destroySession(session, options) {
      await deleteData(session.id);
      return cookie.serialize("", {
        ...options,
        maxAge: void 0,
        expires: /* @__PURE__ */ new Date(0)
      });
    }
  };
}
function warnOnceAboutSigningSessionCookie(cookie) {
  _chunk7OQROU2Djs.warnOnce.call(void 0, 
    cookie.isSigned,
    `The "${cookie.name}" cookie is not signed, but session cookies should be signed to prevent tampering on the client before they are sent back to the server. See https://reactrouter.com/explanation/sessions-and-cookies#signing-cookies for more information.`
  );
}

// lib/server-runtime/sessions/cookieStorage.ts
function createCookieSessionStorage({ cookie: cookieArg } = {}) {
  let cookie = isCookie(cookieArg) ? cookieArg : createCookie(_optionalChain([cookieArg, 'optionalAccess', _18 => _18.name]) || "__session", cookieArg);
  warnOnceAboutSigningSessionCookie(cookie);
  return {
    async getSession(cookieHeader, options) {
      return createSession(
        cookieHeader && await cookie.parse(cookieHeader, options) || {}
      );
    },
    async commitSession(session, options) {
      let serializedCookie = await cookie.serialize(session.data, options);
      if (serializedCookie.length > 4096) {
        throw new Error(
          "Cookie length will exceed browser maximum. Length: " + serializedCookie.length
        );
      }
      return serializedCookie;
    },
    async destroySession(_session, options) {
      return cookie.serialize("", {
        ...options,
        maxAge: void 0,
        expires: /* @__PURE__ */ new Date(0)
      });
    }
  };
}

// lib/server-runtime/sessions/memoryStorage.ts
function createMemorySessionStorage({ cookie } = {}) {
  let map = /* @__PURE__ */ new Map();
  return createSessionStorage({
    cookie,
    async createData(data2, expires) {
      let id = Math.random().toString(36).substring(2, 10);
      map.set(id, { data: data2, expires });
      return id;
    },
    async readData(id) {
      if (map.has(id)) {
        let { data: data2, expires } = map.get(id);
        if (!expires || expires > /* @__PURE__ */ new Date()) {
          return data2;
        }
        if (expires) map.delete(id);
      }
      return null;
    },
    async updateData(id, data2, expires) {
      map.set(id, { data: data2, expires });
    },
    async deleteData(id) {
      map.delete(id);
    }
  });
}

// lib/href.ts
function href(path, ...args) {
  let params = args[0];
  return path.split("/").map((segment) => {
    if (segment === "*") {
      return params ? params["*"] : void 0;
    }
    const match = segment.match(/^:([\w-]+)(\?)?/);
    if (!match) return segment;
    const param = match[1];
    const value = params ? params[param] : void 0;
    const isRequired = match[2] === void 0;
    if (isRequired && value === void 0) {
      throw Error(
        `Path '${path}' requires param '${param}' but it was not provided`
      );
    }
    return value;
  }).filter((segment) => segment !== void 0).join("/");
}

// lib/rsc/browser.tsx

var _reactdom = require('react-dom'); var ReactDOM = _interopRequireWildcard(_reactdom);

// lib/dom/ssr/hydration.tsx
function getHydrationData(state, routes, getRouteInfo, location2, basename, isSpaMode) {
  let hydrationData = {
    ...state,
    loaderData: { ...state.loaderData }
  };
  let initialMatches = _chunk7OQROU2Djs.matchRoutes.call(void 0, routes, location2, basename);
  if (initialMatches) {
    for (let match of initialMatches) {
      let routeId = match.route.id;
      let routeInfo = getRouteInfo(routeId);
      if (_chunk7OQROU2Djs.shouldHydrateRouteLoader.call(void 0, 
        routeId,
        routeInfo.clientLoader,
        routeInfo.hasLoader,
        isSpaMode
      ) && (routeInfo.hasHydrateFallback || !routeInfo.hasLoader)) {
        delete hydrationData.loaderData[routeId];
      } else if (!routeInfo.hasLoader) {
        hydrationData.loaderData[routeId] = null;
      }
    }
  }
  return hydrationData;
}

// lib/rsc/errorBoundaries.tsx

var RSCRouterGlobalErrorBoundary = class extends React.default.Component {
  constructor(props) {
    super(props);
    this.state = { error: null, location: props.location };
  }
  static getDerivedStateFromError(error) {
    return { error };
  }
  static getDerivedStateFromProps(props, state) {
    if (state.location !== props.location) {
      return { error: null, location: props.location };
    }
    return { error: state.error, location: state.location };
  }
  render() {
    if (this.state.error) {
      return /* @__PURE__ */ React.default.createElement(
        RSCDefaultRootErrorBoundaryImpl,
        {
          error: this.state.error,
          renderAppShell: true
        }
      );
    } else {
      return this.props.children;
    }
  }
};
function ErrorWrapper({
  renderAppShell,
  title,
  children
}) {
  if (!renderAppShell) {
    return children;
  }
  return /* @__PURE__ */ React.default.createElement("html", { lang: "en" }, /* @__PURE__ */ React.default.createElement("head", null, /* @__PURE__ */ React.default.createElement("meta", { charSet: "utf-8" }), /* @__PURE__ */ React.default.createElement(
    "meta",
    {
      name: "viewport",
      content: "width=device-width,initial-scale=1,viewport-fit=cover"
    }
  ), /* @__PURE__ */ React.default.createElement("title", null, title)), /* @__PURE__ */ React.default.createElement("body", null, /* @__PURE__ */ React.default.createElement("main", { style: { fontFamily: "system-ui, sans-serif", padding: "2rem" } }, children)));
}
function RSCDefaultRootErrorBoundaryImpl({
  error,
  renderAppShell
}) {
  console.error(error);
  let heyDeveloper = /* @__PURE__ */ React.default.createElement(
    "script",
    {
      dangerouslySetInnerHTML: {
        __html: `
        console.log(
          "\u{1F4BF} Hey developer \u{1F44B}. You can provide a way better UX than this when your app throws errors. Check out https://reactrouter.com/how-to/error-boundary for more information."
        );
      `
      }
    }
  );
  if (_chunk7OQROU2Djs.isRouteErrorResponse.call(void 0, error)) {
    return /* @__PURE__ */ React.default.createElement(
      ErrorWrapper,
      {
        renderAppShell,
        title: "Unhandled Thrown Response!"
      },
      /* @__PURE__ */ React.default.createElement("h1", { style: { fontSize: "24px" } }, error.status, " ", error.statusText),
      _chunk7OQROU2Djs.ENABLE_DEV_WARNINGS ? heyDeveloper : null
    );
  }
  let errorInstance;
  if (error instanceof Error) {
    errorInstance = error;
  } else {
    let errorString = error == null ? "Unknown Error" : typeof error === "object" && "toString" in error ? error.toString() : JSON.stringify(error);
    errorInstance = new Error(errorString);
  }
  return /* @__PURE__ */ React.default.createElement(ErrorWrapper, { renderAppShell, title: "Application Error!" }, /* @__PURE__ */ React.default.createElement("h1", { style: { fontSize: "24px" } }, "Application Error"), /* @__PURE__ */ React.default.createElement(
    "pre",
    {
      style: {
        padding: "2rem",
        background: "hsla(10, 50%, 50%, 0.1)",
        color: "red",
        overflow: "auto"
      }
    },
    errorInstance.stack
  ), heyDeveloper);
}
function RSCDefaultRootErrorBoundary({
  hasRootLayout
}) {
  let error = _chunk7OQROU2Djs.useRouteError.call(void 0, );
  if (hasRootLayout === void 0) {
    throw new Error("Missing 'hasRootLayout' prop");
  }
  return /* @__PURE__ */ React.default.createElement(
    RSCDefaultRootErrorBoundaryImpl,
    {
      renderAppShell: !hasRootLayout,
      error
    }
  );
}

// lib/rsc/browser.tsx
function createCallServer({
  createFromReadableStream,
  createTemporaryReferenceSet,
  encodeReply,
  fetch: fetchImplementation = fetch
}) {
  const globalVar = window;
  let landedActionId = 0;
  return async (id, args) => {
    let actionId = globalVar.__routerActionID = (_nullishCoalesce(globalVar.__routerActionID, () => ( (globalVar.__routerActionID = 0)))) + 1;
    const temporaryReferences = createTemporaryReferenceSet();
    const response = await fetchImplementation(
      new Request(location.href, {
        body: await encodeReply(args, { temporaryReferences }),
        method: "POST",
        headers: {
          Accept: "text/x-component",
          "rsc-action-id": id
        }
      })
    );
    if (!response.body) {
      throw new Error("No response body");
    }
    const payload = await createFromReadableStream(response.body, {
      temporaryReferences
    });
    if (payload.type === "redirect") {
      if (payload.reload) {
        window.location.href = payload.location;
        return;
      }
      globalVar.__router.navigate(payload.location, {
        replace: payload.replace
      });
      return payload.actionResult;
    }
    if (payload.type !== "action") {
      throw new Error("Unexpected payload type");
    }
    if (payload.rerender) {
      React4.startTransition(
        // @ts-expect-error - We have old react types that don't know this can be async
        async () => {
          const rerender = await payload.rerender;
          if (!rerender) return;
          if (landedActionId < actionId && globalVar.__routerActionID <= actionId) {
            landedActionId = actionId;
            if (rerender.type === "redirect") {
              if (rerender.reload) {
                window.location.href = rerender.location;
                return;
              }
              globalVar.__router.navigate(rerender.location, {
                replace: rerender.replace
              });
              return;
            }
            let lastMatch;
            for (const match of rerender.matches) {
              globalVar.__router.patchRoutes(
                _nullishCoalesce(_optionalChain([lastMatch, 'optionalAccess', _19 => _19.id]), () => ( null)),
                [createRouteFromServerManifest(match)],
                true
              );
              lastMatch = match;
            }
            window.__router._internalSetStateDoNotUseOrYouWillBreakYourApp({});
            React4.startTransition(() => {
              window.__router._internalSetStateDoNotUseOrYouWillBreakYourApp({
                loaderData: Object.assign(
                  {},
                  globalVar.__router.state.loaderData,
                  rerender.loaderData
                ),
                errors: rerender.errors ? Object.assign(
                  {},
                  globalVar.__router.state.errors,
                  rerender.errors
                ) : null
              });
            });
          }
        }
      );
    }
    return payload.actionResult;
  };
}
function createRouterFromPayload({
  fetchImplementation,
  createFromReadableStream,
  unstable_getContext,
  payload
}) {
  const globalVar = window;
  if (globalVar.__router) return globalVar.__router;
  if (payload.type !== "render") throw new Error("Invalid payload type");
  let patches = /* @__PURE__ */ new Map();
  _optionalChain([payload, 'access', _20 => _20.patches, 'optionalAccess', _21 => _21.forEach, 'call', _22 => _22((patch) => {
    _chunk7OQROU2Djs.invariant.call(void 0, patch.parentId, "Invalid patch parentId");
    if (!patches.has(patch.parentId)) {
      patches.set(patch.parentId, []);
    }
    _optionalChain([patches, 'access', _23 => _23.get, 'call', _24 => _24(patch.parentId), 'optionalAccess', _25 => _25.push, 'call', _26 => _26(patch)]);
  })]);
  let routes = payload.matches.reduceRight((previous, match) => {
    const route = createRouteFromServerManifest(
      match,
      payload
    );
    if (previous.length > 0) {
      route.children = previous;
      let childrenToPatch = patches.get(match.id);
      if (childrenToPatch) {
        route.children.push(
          ...childrenToPatch.map((r) => createRouteFromServerManifest(r))
        );
      }
    }
    return [route];
  }, []);
  globalVar.__router = _chunk7OQROU2Djs.createRouter.call(void 0, {
    routes,
    unstable_getContext,
    basename: payload.basename,
    history: _chunk7OQROU2Djs.createBrowserHistory.call(void 0, ),
    hydrationData: getHydrationData(
      {
        loaderData: payload.loaderData,
        actionData: payload.actionData,
        errors: payload.errors
      },
      routes,
      (routeId) => {
        let match = payload.matches.find((m) => m.id === routeId);
        _chunk7OQROU2Djs.invariant.call(void 0, match, "Route not found in payload");
        return {
          clientLoader: match.clientLoader,
          hasLoader: match.hasLoader,
          hasHydrateFallback: match.hydrateFallbackElement != null
        };
      },
      payload.location,
      void 0,
      false
    ),
    async patchRoutesOnNavigation({ path, signal }) {
      if (discoveredPaths.has(path)) {
        return;
      }
      await fetchAndApplyManifestPatches(
        [path],
        createFromReadableStream,
        fetchImplementation,
        signal
      );
    },
    // FIXME: Pass `build.ssr` into this function
    dataStrategy: getRSCSingleFetchDataStrategy(
      () => globalVar.__router,
      true,
      payload.basename,
      createFromReadableStream,
      fetchImplementation
    )
  });
  if (globalVar.__router.state.initialized) {
    globalVar.__routerInitialized = true;
    globalVar.__router.initialize();
  } else {
    globalVar.__routerInitialized = false;
  }
  let lastLoaderData = void 0;
  globalVar.__router.subscribe(({ loaderData, actionData }) => {
    if (lastLoaderData !== loaderData) {
      globalVar.__routerActionID = (_nullishCoalesce(globalVar.__routerActionID, () => ( (globalVar.__routerActionID = 0)))) + 1;
    }
  });
  return globalVar.__router;
}
var renderedRoutesContext = _chunk7OQROU2Djs.unstable_createContext.call(void 0, );
function getRSCSingleFetchDataStrategy(getRouter, ssr, basename, createFromReadableStream, fetchImplementation) {
  let dataStrategy = _chunk7OQROU2Djs.getSingleFetchDataStrategyImpl.call(void 0, 
    getRouter,
    (match) => {
      let M = match;
      return {
        hasLoader: M.route.hasLoader,
        hasClientLoader: M.route.hasClientLoader,
        hasComponent: M.route.hasComponent,
        hasAction: M.route.hasAction,
        hasClientAction: M.route.hasClientAction,
        hasShouldRevalidate: M.route.hasShouldRevalidate
      };
    },
    // pass map into fetchAndDecode so it can add payloads
    getFetchAndDecodeViaRSC(createFromReadableStream, fetchImplementation),
    ssr,
    basename,
    // If the route has a component but we don't have an element, we need to hit
    // the server loader flow regardless of whether the client loader calls
    // `serverLoader` or not, otherwise we'll have nothing to render.
    (match) => {
      let M = match;
      return M.route.hasComponent && !M.route.element;
    }
  );
  return async (args) => args.unstable_runClientMiddleware(async () => {
    let context = args.context;
    context.set(renderedRoutesContext, []);
    let results = await dataStrategy(args);
    const renderedRoutesById = /* @__PURE__ */ new Map();
    for (const route of context.get(renderedRoutesContext)) {
      if (!renderedRoutesById.has(route.id)) {
        renderedRoutesById.set(route.id, []);
      }
      renderedRoutesById.get(route.id).push(route);
    }
    for (const match of args.matches) {
      const renderedRoutes = renderedRoutesById.get(match.route.id);
      if (renderedRoutes) {
        for (const rendered of renderedRoutes) {
          window.__router.patchRoutes(
            _nullishCoalesce(rendered.parentId, () => ( null)),
            [createRouteFromServerManifest(rendered)],
            true
          );
        }
      }
    }
    return results;
  });
}
function getFetchAndDecodeViaRSC(createFromReadableStream, fetchImplementation) {
  return async (args, basename, targetRoutes) => {
    let { request, context } = args;
    let url = _chunk7OQROU2Djs.singleFetchUrl.call(void 0, request.url, basename, "rsc");
    if (request.method === "GET") {
      url = _chunk7OQROU2Djs.stripIndexParam.call(void 0, url);
      if (targetRoutes) {
        url.searchParams.set("_routes", targetRoutes.join(","));
      }
    }
    let res = await fetchImplementation(
      new Request(url, await _chunk7OQROU2Djs.createRequestInit.call(void 0, request))
    );
    if (res.status === 404 && !res.headers.has("X-Remix-Response")) {
      throw new (0, _chunk7OQROU2Djs.ErrorResponseImpl)(404, "Not Found", true);
    }
    _chunk7OQROU2Djs.invariant.call(void 0, res.body, "No response body to decode");
    try {
      const payload = await createFromReadableStream(res.body, {
        temporaryReferences: void 0
      });
      if (payload.type === "redirect") {
        return {
          status: res.status,
          data: {
            redirect: {
              redirect: payload.location,
              reload: payload.reload,
              replace: payload.replace,
              revalidate: false,
              status: payload.status
            }
          }
        };
      }
      if (payload.type !== "render") {
        throw new Error("Unexpected payload type");
      }
      context.get(renderedRoutesContext).push(...payload.matches);
      let results = { routes: {} };
      const dataKey = _chunk7OQROU2Djs.isMutationMethod.call(void 0, request.method) ? "actionData" : "loaderData";
      for (let [routeId, data2] of Object.entries(payload[dataKey] || {})) {
        results.routes[routeId] = { data: data2 };
      }
      if (payload.errors) {
        for (let [routeId, error] of Object.entries(payload.errors)) {
          results.routes[routeId] = { error };
        }
      }
      return { status: res.status, data: results };
    } catch (e) {
      throw new Error("Unable to decode RSC response");
    }
  };
}
function RSCHydratedRouter({
  createFromReadableStream,
  fetch: fetchImplementation = fetch,
  payload,
  routeDiscovery = "eager",
  unstable_getContext
}) {
  if (payload.type !== "render") throw new Error("Invalid payload type");
  let router = React4.useMemo(
    () => createRouterFromPayload({
      payload,
      fetchImplementation,
      unstable_getContext,
      createFromReadableStream
    }),
    [
      createFromReadableStream,
      payload,
      fetchImplementation,
      unstable_getContext
    ]
  );
  React4.useLayoutEffect(() => {
    const globalVar = window;
    if (!globalVar.__routerInitialized) {
      globalVar.__routerInitialized = true;
      globalVar.__router.initialize();
    }
  }, []);
  let [location2, setLocation] = React4.useState(router.state.location);
  React4.useLayoutEffect(
    () => router.subscribe((newState) => {
      if (newState.location !== location2) {
        setLocation(newState.location);
      }
    }),
    [router, location2]
  );
  React4.useEffect(() => {
    if (routeDiscovery === "lazy" || // @ts-expect-error - TS doesn't know about this yet
    _optionalChain([window, 'access', _27 => _27.navigator, 'optionalAccess', _28 => _28.connection, 'optionalAccess', _29 => _29.saveData]) === true) {
      return;
    }
    function registerElement(el) {
      let path = el.tagName === "FORM" ? el.getAttribute("action") : el.getAttribute("href");
      if (!path) {
        return;
      }
      let pathname = el.tagName === "A" ? el.pathname : new URL(path, window.location.origin).pathname;
      if (!discoveredPaths.has(pathname)) {
        nextPaths.add(pathname);
      }
    }
    async function fetchPatches() {
      document.querySelectorAll("a[data-discover], form[data-discover]").forEach(registerElement);
      let paths = Array.from(nextPaths.keys()).filter((path) => {
        if (discoveredPaths.has(path)) {
          nextPaths.delete(path);
          return false;
        }
        return true;
      });
      if (paths.length === 0) {
        return;
      }
      try {
        await fetchAndApplyManifestPatches(
          paths,
          createFromReadableStream,
          fetchImplementation
        );
      } catch (e) {
        console.error("Failed to fetch manifest patches", e);
      }
    }
    let debouncedFetchPatches = debounce(fetchPatches, 100);
    fetchPatches();
    let observer = new MutationObserver(() => debouncedFetchPatches());
    observer.observe(document.documentElement, {
      subtree: true,
      childList: true,
      attributes: true,
      attributeFilter: ["data-discover", "href", "action"]
    });
  }, [routeDiscovery, createFromReadableStream, fetchImplementation]);
  const frameworkContext = {
    future: {
      // These flags have no runtime impact so can always be false.  If we add
      // flags that drive runtime behavior they'll need to be proxied through.
      unstable_middleware: false,
      unstable_subResourceIntegrity: false
    },
    isSpaMode: true,
    ssr: true,
    criticalCss: "",
    manifest: {
      routes: {},
      version: "1",
      url: "",
      entry: {
        module: "",
        imports: []
      }
    },
    routeDiscovery: { mode: "lazy", manifestPath: "/__manifest" },
    routeModules: {}
  };
  return /* @__PURE__ */ React4.createElement(_chunk7OQROU2Djs.RSCRouterContext.Provider, { value: true }, /* @__PURE__ */ React4.createElement(RSCRouterGlobalErrorBoundary, { location: location2 }, /* @__PURE__ */ React4.createElement(_chunk7OQROU2Djs.FrameworkContext.Provider, { value: frameworkContext }, /* @__PURE__ */ React4.createElement(_chunk4DGLNKXFjs.RouterProvider, { router, flushSync: ReactDOM.flushSync }))));
}
function createRouteFromServerManifest(match, payload) {
  let hasInitialData = payload && match.id in payload.loaderData;
  let initialData = _optionalChain([payload, 'optionalAccess', _30 => _30.loaderData, 'access', _31 => _31[match.id]]);
  let hasInitialError = _optionalChain([payload, 'optionalAccess', _32 => _32.errors]) && match.id in payload.errors;
  let initialError = _optionalChain([payload, 'optionalAccess', _33 => _33.errors, 'optionalAccess', _34 => _34[match.id]]);
  let isHydrationRequest = _optionalChain([match, 'access', _35 => _35.clientLoader, 'optionalAccess', _36 => _36.hydrate]) === true || !match.hasLoader || // If the route has a component but we don't have an element, we need to hit
  // the server loader flow regardless of whether the client loader calls
  // `serverLoader` or not, otherwise we'll have nothing to render.
  match.hasComponent && !match.element;
  let dataRoute = {
    id: match.id,
    element: match.element,
    errorElement: match.errorElement,
    handle: match.handle,
    hasErrorBoundary: match.hasErrorBoundary,
    hydrateFallbackElement: match.hydrateFallbackElement,
    index: match.index,
    loader: match.clientLoader ? async (args, singleFetch) => {
      try {
        let result = await match.clientLoader({
          ...args,
          serverLoader: () => {
            preventInvalidServerHandlerCall(
              "loader",
              match.id,
              match.hasLoader
            );
            if (isHydrationRequest) {
              if (hasInitialData) {
                return initialData;
              }
              if (hasInitialError) {
                throw initialError;
              }
            }
            return callSingleFetch(singleFetch);
          }
        });
        return result;
      } finally {
        isHydrationRequest = false;
      }
    } : (
      // We always make the call in this RSC world since even if we don't
      // have a `loader` we may need to get the `element` implementation
      (_, singleFetch) => callSingleFetch(singleFetch)
    ),
    action: match.clientAction ? (args, singleFetch) => match.clientAction({
      ...args,
      serverAction: async () => {
        preventInvalidServerHandlerCall(
          "action",
          match.id,
          match.hasLoader
        );
        return await callSingleFetch(singleFetch);
      }
    }) : match.hasAction ? (_, singleFetch) => callSingleFetch(singleFetch) : () => {
      throw _chunk7OQROU2Djs.noActionDefinedError.call(void 0, "action", match.id);
    },
    path: match.path,
    shouldRevalidate: match.shouldRevalidate,
    // We always have a "loader" in this RSC world since even if we don't
    // have a `loader` we may need to get the `element` implementation
    hasLoader: true,
    hasClientLoader: match.clientLoader != null,
    hasAction: match.hasAction,
    hasClientAction: match.clientAction != null,
    hasShouldRevalidate: match.shouldRevalidate != null
  };
  if (typeof dataRoute.loader === "function") {
    dataRoute.loader.hydrate = _chunk7OQROU2Djs.shouldHydrateRouteLoader.call(void 0, 
      match.id,
      match.clientLoader,
      match.hasLoader,
      false
    );
  }
  return dataRoute;
}
function callSingleFetch(singleFetch) {
  _chunk7OQROU2Djs.invariant.call(void 0, typeof singleFetch === "function", "Invalid singleFetch parameter");
  return singleFetch();
}
function preventInvalidServerHandlerCall(type, routeId, hasHandler) {
  if (!hasHandler) {
    let fn = type === "action" ? "serverAction()" : "serverLoader()";
    let msg = `You are trying to call ${fn} on a route that does not have a server ${type} (routeId: "${routeId}")`;
    console.error(msg);
    throw new (0, _chunk7OQROU2Djs.ErrorResponseImpl)(400, "Bad Request", new Error(msg), true);
  }
}
var nextPaths = /* @__PURE__ */ new Set();
var discoveredPathsMaxSize = 1e3;
var discoveredPaths = /* @__PURE__ */ new Set();
var URL_LIMIT = 7680;
function getManifestUrl(paths) {
  if (paths.length === 0) {
    return null;
  }
  if (paths.length === 1) {
    return new URL(`${paths[0]}.manifest`, window.location.origin);
  }
  const globalVar = window;
  let basename = (_nullishCoalesce(globalVar.__router.basename, () => ( ""))).replace(/^\/|\/$/g, "");
  let url = new URL(`${basename}/.manifest`, window.location.origin);
  paths.sort().forEach((path) => url.searchParams.append("p", path));
  return url;
}
async function fetchAndApplyManifestPatches(paths, createFromReadableStream, fetchImplementation, signal) {
  let url = getManifestUrl(paths);
  if (url == null) {
    return;
  }
  if (url.toString().length > URL_LIMIT) {
    nextPaths.clear();
    return;
  }
  let response = await fetchImplementation(new Request(url, { signal }));
  if (!response.body || response.status < 200 || response.status >= 300) {
    throw new Error("Unable to fetch new route matches from the server");
  }
  let payload = await createFromReadableStream(response.body, {
    temporaryReferences: void 0
  });
  if (payload.type !== "manifest") {
    throw new Error("Failed to patch routes");
  }
  paths.forEach((p) => addToFifoQueue(p, discoveredPaths));
  payload.patches.forEach((p) => {
    window.__router.patchRoutes(
      _nullishCoalesce(p.parentId, () => ( null)),
      [createRouteFromServerManifest(p)]
    );
  });
}
function addToFifoQueue(path, queue) {
  if (queue.size >= discoveredPathsMaxSize) {
    let first = queue.values().next().value;
    queue.delete(first);
  }
  queue.add(path);
}
function debounce(callback, wait) {
  let timeoutId;
  return (...args) => {
    window.clearTimeout(timeoutId);
    timeoutId = window.setTimeout(() => callback(...args), wait);
  };
}

// lib/rsc/server.ssr.tsx


// lib/rsc/html-stream/server.ts
var encoder2 = new TextEncoder();
var trailer = "</body></html>";
function injectRSCPayload(rscStream) {
  let decoder = new TextDecoder();
  let resolveFlightDataPromise;
  let flightDataPromise = new Promise(
    (resolve) => resolveFlightDataPromise = resolve
  );
  let startedRSC = false;
  let buffered = [];
  let timeout = null;
  function flushBufferedChunks(controller) {
    for (let chunk of buffered) {
      let buf = decoder.decode(chunk, { stream: true });
      if (buf.endsWith(trailer)) {
        buf = buf.slice(0, -trailer.length);
      }
      controller.enqueue(encoder2.encode(buf));
    }
    buffered.length = 0;
    timeout = null;
  }
  return new TransformStream({
    transform(chunk, controller) {
      buffered.push(chunk);
      if (timeout) {
        return;
      }
      timeout = setTimeout(async () => {
        flushBufferedChunks(controller);
        if (!startedRSC) {
          startedRSC = true;
          writeRSCStream(rscStream, controller).catch((err) => controller.error(err)).then(resolveFlightDataPromise);
        }
      }, 0);
    },
    async flush(controller) {
      await flightDataPromise;
      if (timeout) {
        clearTimeout(timeout);
        flushBufferedChunks(controller);
      }
      controller.enqueue(encoder2.encode("</body></html>"));
    }
  });
}
async function writeRSCStream(rscStream, controller) {
  let decoder = new TextDecoder("utf-8", { fatal: true });
  const reader = rscStream.getReader();
  try {
    let read;
    while ((read = await reader.read()) && !read.done) {
      const chunk = read.value;
      try {
        writeChunk(
          JSON.stringify(decoder.decode(chunk, { stream: true })),
          controller
        );
      } catch (err) {
        let base64 = JSON.stringify(btoa(String.fromCodePoint(...chunk)));
        writeChunk(
          `Uint8Array.from(atob(${base64}), m => m.codePointAt(0))`,
          controller
        );
      }
    }
  } finally {
    reader.releaseLock();
  }
  let remaining = decoder.decode();
  if (remaining.length) {
    writeChunk(JSON.stringify(remaining), controller);
  }
}
function writeChunk(chunk, controller) {
  controller.enqueue(
    encoder2.encode(
      `<script>${escapeScript(
        `(self.__FLIGHT_DATA||=[]).push(${chunk})`
      )}</script>`
    )
  );
}
function escapeScript(script) {
  return script.replace(/<!--/g, "<\\!--").replace(/<\/(script)/gi, "</\\$1");
}

// lib/rsc/server.ssr.tsx
async function routeRSCServerRequest({
  request,
  fetchServer,
  createFromReadableStream,
  renderHTML,
  hydrate = true
}) {
  const url = new URL(request.url);
  const isDataRequest = isReactServerRequest(url);
  const respondWithRSCPayload = isDataRequest || isManifestRequest(url) || request.headers.has("rsc-action-id");
  const serverResponse = await fetchServer(request);
  if (respondWithRSCPayload || serverResponse.headers.get("React-Router-Resource") === "true") {
    return serverResponse;
  }
  if (!serverResponse.body) {
    throw new Error("Missing body in server response");
  }
  let serverResponseB = null;
  if (hydrate) {
    serverResponseB = serverResponse.clone();
  }
  const body = serverResponse.body;
  let payloadPromise;
  const getPayload = async () => {
    if (payloadPromise) return payloadPromise;
    payloadPromise = createFromReadableStream(body);
    return payloadPromise;
  };
  try {
    const html = await renderHTML(getPayload);
    const headers = new Headers(serverResponse.headers);
    headers.set("Content-Type", "text/html");
    if (!hydrate) {
      return new Response(html, {
        status: serverResponse.status,
        headers
      });
    }
    if (!_optionalChain([serverResponseB, 'optionalAccess', _37 => _37.body])) {
      throw new Error("Failed to clone server response");
    }
    const body2 = html.pipeThrough(injectRSCPayload(serverResponseB.body));
    return new Response(body2, {
      status: serverResponse.status,
      headers
    });
  } catch (reason) {
    if (reason instanceof Response) {
      return reason;
    }
    throw reason;
  }
}
function RSCStaticRouter({ getPayload }) {
  const payload = React5.use(getPayload());
  if (payload.type === "redirect") {
    throw new Response(null, {
      status: payload.status,
      headers: {
        Location: payload.location
      }
    });
  }
  if (payload.type !== "render") return null;
  let patchedLoaderData = { ...payload.loaderData };
  for (const match of payload.matches) {
    if (_chunk7OQROU2Djs.shouldHydrateRouteLoader.call(void 0, 
      match.id,
      match.clientLoader,
      match.hasLoader,
      false
    ) && (match.hydrateFallbackElement || !match.hasLoader)) {
      delete patchedLoaderData[match.id];
    }
  }
  const context = {
    actionData: payload.actionData,
    actionHeaders: {},
    basename: payload.basename,
    errors: payload.errors,
    loaderData: patchedLoaderData,
    loaderHeaders: {},
    location: payload.location,
    statusCode: 200,
    matches: payload.matches.map((match) => ({
      params: match.params,
      pathname: match.pathname,
      pathnameBase: match.pathnameBase,
      route: {
        id: match.id,
        action: match.hasAction || !!match.clientAction,
        handle: match.handle,
        hasErrorBoundary: match.hasErrorBoundary,
        loader: match.hasLoader || !!match.clientLoader,
        index: match.index,
        path: match.path,
        shouldRevalidate: match.shouldRevalidate
      }
    }))
  };
  const router = _chunk4DGLNKXFjs.createStaticRouter.call(void 0, 
    payload.matches.reduceRight((previous, match) => {
      const route = {
        id: match.id,
        action: match.hasAction || !!match.clientAction,
        element: match.element,
        errorElement: match.errorElement,
        handle: match.handle,
        hasErrorBoundary: !!match.errorElement,
        hydrateFallbackElement: match.hydrateFallbackElement,
        index: match.index,
        loader: match.hasLoader || !!match.clientLoader,
        path: match.path,
        shouldRevalidate: match.shouldRevalidate
      };
      if (previous.length > 0) {
        route.children = previous;
      }
      return [route];
    }, []),
    context
  );
  const frameworkContext = {
    future: {
      // These flags have no runtime impact so can always be false.  If we add
      // flags that drive runtime behavior they'll need to be proxied through.
      unstable_middleware: false,
      unstable_subResourceIntegrity: false
    },
    isSpaMode: false,
    ssr: true,
    criticalCss: "",
    manifest: {
      routes: {},
      version: "1",
      url: "",
      entry: {
        module: "",
        imports: []
      }
    },
    routeDiscovery: { mode: "lazy", manifestPath: "/__manifest" },
    routeModules: {}
  };
  return /* @__PURE__ */ React5.createElement(_chunk7OQROU2Djs.RSCRouterContext.Provider, { value: true }, /* @__PURE__ */ React5.createElement(RSCRouterGlobalErrorBoundary, { location: payload.location }, /* @__PURE__ */ React5.createElement(_chunk7OQROU2Djs.FrameworkContext.Provider, { value: frameworkContext }, /* @__PURE__ */ React5.createElement(
    _chunk4DGLNKXFjs.StaticRouterProvider,
    {
      context,
      router,
      hydrate: false,
      nonce: payload.nonce
    }
  ))));
}
function isReactServerRequest(url) {
  return url.pathname.endsWith(".rsc");
}
function isManifestRequest(url) {
  return url.pathname.endsWith(".manifest");
}

// lib/rsc/html-stream/browser.ts
function getRSCStream() {
  let encoder3 = new TextEncoder();
  let streamController = null;
  let rscStream = new ReadableStream({
    start(controller) {
      if (typeof window === "undefined") {
        return;
      }
      let handleChunk = (chunk) => {
        if (typeof chunk === "string") {
          controller.enqueue(encoder3.encode(chunk));
        } else {
          controller.enqueue(chunk);
        }
      };
      window.__FLIGHT_DATA || (window.__FLIGHT_DATA = []);
      window.__FLIGHT_DATA.forEach(handleChunk);
      window.__FLIGHT_DATA.push = (chunk) => {
        handleChunk(chunk);
        return 0;
      };
      streamController = controller;
    }
  });
  if (typeof document !== "undefined" && document.readyState === "loading") {
    document.addEventListener("DOMContentLoaded", () => {
      _optionalChain([streamController, 'optionalAccess', _38 => _38.close, 'call', _39 => _39()]);
    });
  } else {
    _optionalChain([streamController, 'optionalAccess', _40 => _40.close, 'call', _41 => _41()]);
  }
  return rscStream;
}

// lib/dom/ssr/errors.ts
function deserializeErrors(errors) {
  if (!errors) return null;
  let entries = Object.entries(errors);
  let serialized = {};
  for (let [key, val] of entries) {
    if (val && val.__type === "RouteErrorResponse") {
      serialized[key] = new (0, _chunk7OQROU2Djs.ErrorResponseImpl)(
        val.status,
        val.statusText,
        val.data,
        val.internal === true
      );
    } else if (val && val.__type === "Error") {
      if (val.__subType) {
        let ErrorConstructor = window[val.__subType];
        if (typeof ErrorConstructor === "function") {
          try {
            let error = new ErrorConstructor(val.message);
            error.stack = val.stack;
            serialized[key] = error;
          } catch (e) {
          }
        }
      }
      if (serialized[key] == null) {
        let error = new Error(val.message);
        error.stack = val.stack;
        serialized[key] = error;
      }
    } else {
      serialized[key] = val;
    }
  }
  return serialized;
}

































































































































exports.Await = _chunk4DGLNKXFjs.Await; exports.BrowserRouter = _chunk4DGLNKXFjs.BrowserRouter; exports.Form = _chunk4DGLNKXFjs.Form; exports.HashRouter = _chunk4DGLNKXFjs.HashRouter; exports.IDLE_BLOCKER = _chunk7OQROU2Djs.IDLE_BLOCKER; exports.IDLE_FETCHER = _chunk7OQROU2Djs.IDLE_FETCHER; exports.IDLE_NAVIGATION = _chunk7OQROU2Djs.IDLE_NAVIGATION; exports.Link = _chunk4DGLNKXFjs.Link; exports.Links = _chunk7OQROU2Djs.Links; exports.MemoryRouter = _chunk4DGLNKXFjs.MemoryRouter; exports.Meta = _chunk7OQROU2Djs.Meta; exports.NavLink = _chunk4DGLNKXFjs.NavLink; exports.Navigate = _chunk4DGLNKXFjs.Navigate; exports.NavigationType = _chunk7OQROU2Djs.Action; exports.Outlet = _chunk4DGLNKXFjs.Outlet; exports.PrefetchPageLinks = _chunk7OQROU2Djs.PrefetchPageLinks; exports.Route = _chunk4DGLNKXFjs.Route; exports.Router = _chunk4DGLNKXFjs.Router; exports.RouterProvider = _chunk4DGLNKXFjs.RouterProvider; exports.Routes = _chunk4DGLNKXFjs.Routes; exports.Scripts = _chunk7OQROU2Djs.Scripts; exports.ScrollRestoration = _chunk4DGLNKXFjs.ScrollRestoration; exports.ServerRouter = ServerRouter; exports.StaticRouter = _chunk4DGLNKXFjs.StaticRouter; exports.StaticRouterProvider = _chunk4DGLNKXFjs.StaticRouterProvider; exports.UNSAFE_DataRouterContext = _chunk7OQROU2Djs.DataRouterContext; exports.UNSAFE_DataRouterStateContext = _chunk7OQROU2Djs.DataRouterStateContext; exports.UNSAFE_ErrorResponseImpl = _chunk7OQROU2Djs.ErrorResponseImpl; exports.UNSAFE_FetchersContext = _chunk7OQROU2Djs.FetchersContext; exports.UNSAFE_FrameworkContext = _chunk7OQROU2Djs.FrameworkContext; exports.UNSAFE_LocationContext = _chunk7OQROU2Djs.LocationContext; exports.UNSAFE_NavigationContext = _chunk7OQROU2Djs.NavigationContext; exports.UNSAFE_RSCDefaultRootErrorBoundary = RSCDefaultRootErrorBoundary; exports.UNSAFE_RemixErrorBoundary = _chunk7OQROU2Djs.RemixErrorBoundary; exports.UNSAFE_RouteContext = _chunk7OQROU2Djs.RouteContext; exports.UNSAFE_ServerMode = ServerMode; exports.UNSAFE_SingleFetchRedirectSymbol = _chunk7OQROU2Djs.SingleFetchRedirectSymbol; exports.UNSAFE_ViewTransitionContext = _chunk7OQROU2Djs.ViewTransitionContext; exports.UNSAFE_WithComponentProps = _chunk4DGLNKXFjs.WithComponentProps; exports.UNSAFE_WithErrorBoundaryProps = _chunk4DGLNKXFjs.WithErrorBoundaryProps; exports.UNSAFE_WithHydrateFallbackProps = _chunk4DGLNKXFjs.WithHydrateFallbackProps; exports.UNSAFE_createBrowserHistory = _chunk7OQROU2Djs.createBrowserHistory; exports.UNSAFE_createClientRoutes = _chunk7OQROU2Djs.createClientRoutes; exports.UNSAFE_createClientRoutesWithHMRRevalidationOptOut = _chunk7OQROU2Djs.createClientRoutesWithHMRRevalidationOptOut; exports.UNSAFE_createRouter = _chunk7OQROU2Djs.createRouter; exports.UNSAFE_decodeViaTurboStream = _chunk7OQROU2Djs.decodeViaTurboStream; exports.UNSAFE_deserializeErrors = deserializeErrors; exports.UNSAFE_getHydrationData = getHydrationData; exports.UNSAFE_getPatchRoutesOnNavigationFunction = _chunk7OQROU2Djs.getPatchRoutesOnNavigationFunction; exports.UNSAFE_getTurboStreamSingleFetchDataStrategy = _chunk7OQROU2Djs.getTurboStreamSingleFetchDataStrategy; exports.UNSAFE_hydrationRouteProperties = _chunk4DGLNKXFjs.hydrationRouteProperties; exports.UNSAFE_invariant = _chunk7OQROU2Djs.invariant; exports.UNSAFE_mapRouteProperties = _chunk4DGLNKXFjs.mapRouteProperties; exports.UNSAFE_shouldHydrateRouteLoader = _chunk7OQROU2Djs.shouldHydrateRouteLoader; exports.UNSAFE_useFogOFWarDiscovery = _chunk7OQROU2Djs.useFogOFWarDiscovery; exports.UNSAFE_useScrollRestoration = _chunk4DGLNKXFjs.useScrollRestoration; exports.UNSAFE_withComponentProps = _chunk4DGLNKXFjs.withComponentProps; exports.UNSAFE_withErrorBoundaryProps = _chunk4DGLNKXFjs.withErrorBoundaryProps; exports.UNSAFE_withHydrateFallbackProps = _chunk4DGLNKXFjs.withHydrateFallbackProps; exports.createBrowserRouter = _chunk4DGLNKXFjs.createBrowserRouter; exports.createCookie = createCookie; exports.createCookieSessionStorage = createCookieSessionStorage; exports.createHashRouter = _chunk4DGLNKXFjs.createHashRouter; exports.createMemoryRouter = _chunk4DGLNKXFjs.createMemoryRouter; exports.createMemorySessionStorage = createMemorySessionStorage; exports.createPath = _chunk7OQROU2Djs.createPath; exports.createRequestHandler = createRequestHandler; exports.createRoutesFromChildren = _chunk4DGLNKXFjs.createRoutesFromChildren; exports.createRoutesFromElements = _chunk4DGLNKXFjs.createRoutesFromElements; exports.createRoutesStub = createRoutesStub; exports.createSearchParams = _chunk4DGLNKXFjs.createSearchParams; exports.createSession = createSession; exports.createSessionStorage = createSessionStorage; exports.createStaticHandler = _chunk4DGLNKXFjs.createStaticHandler; exports.createStaticRouter = _chunk4DGLNKXFjs.createStaticRouter; exports.data = _chunk7OQROU2Djs.data; exports.generatePath = _chunk7OQROU2Djs.generatePath; exports.href = href; exports.isCookie = isCookie; exports.isRouteErrorResponse = _chunk7OQROU2Djs.isRouteErrorResponse; exports.isSession = isSession; exports.matchPath = _chunk7OQROU2Djs.matchPath; exports.matchRoutes = _chunk7OQROU2Djs.matchRoutes; exports.parsePath = _chunk7OQROU2Djs.parsePath; exports.redirect = _chunk7OQROU2Djs.redirect; exports.redirectDocument = _chunk7OQROU2Djs.redirectDocument; exports.renderMatches = _chunk4DGLNKXFjs.renderMatches; exports.replace = _chunk7OQROU2Djs.replace; exports.resolvePath = _chunk7OQROU2Djs.resolvePath; exports.unstable_HistoryRouter = _chunk4DGLNKXFjs.HistoryRouter; exports.unstable_RSCHydratedRouter = RSCHydratedRouter; exports.unstable_RSCStaticRouter = RSCStaticRouter; exports.unstable_RouterContextProvider = _chunk7OQROU2Djs.unstable_RouterContextProvider; exports.unstable_createCallServer = createCallServer; exports.unstable_createContext = _chunk7OQROU2Djs.unstable_createContext; exports.unstable_getRSCStream = getRSCStream; exports.unstable_routeRSCServerRequest = routeRSCServerRequest; exports.unstable_setDevServerHooks = setDevServerHooks; exports.unstable_usePrompt = _chunk4DGLNKXFjs.usePrompt; exports.useActionData = _chunk7OQROU2Djs.useActionData; exports.useAsyncError = _chunk7OQROU2Djs.useAsyncError; exports.useAsyncValue = _chunk7OQROU2Djs.useAsyncValue; exports.useBeforeUnload = _chunk4DGLNKXFjs.useBeforeUnload; exports.useBlocker = _chunk7OQROU2Djs.useBlocker; exports.useFetcher = _chunk4DGLNKXFjs.useFetcher; exports.useFetchers = _chunk4DGLNKXFjs.useFetchers; exports.useFormAction = _chunk4DGLNKXFjs.useFormAction; exports.useHref = _chunk7OQROU2Djs.useHref; exports.useInRouterContext = _chunk7OQROU2Djs.useInRouterContext; exports.useLinkClickHandler = _chunk4DGLNKXFjs.useLinkClickHandler; exports.useLoaderData = _chunk7OQROU2Djs.useLoaderData; exports.useLocation = _chunk7OQROU2Djs.useLocation; exports.useMatch = _chunk7OQROU2Djs.useMatch; exports.useMatches = _chunk7OQROU2Djs.useMatches; exports.useNavigate = _chunk7OQROU2Djs.useNavigate; exports.useNavigation = _chunk7OQROU2Djs.useNavigation; exports.useNavigationType = _chunk7OQROU2Djs.useNavigationType; exports.useOutlet = _chunk7OQROU2Djs.useOutlet; exports.useOutletContext = _chunk7OQROU2Djs.useOutletContext; exports.useParams = _chunk7OQROU2Djs.useParams; exports.useResolvedPath = _chunk7OQROU2Djs.useResolvedPath; exports.useRevalidator = _chunk7OQROU2Djs.useRevalidator; exports.useRouteError = _chunk7OQROU2Djs.useRouteError; exports.useRouteLoaderData = _chunk7OQROU2Djs.useRouteLoaderData; exports.useRoutes = _chunk7OQROU2Djs.useRoutes; exports.useSearchParams = _chunk4DGLNKXFjs.useSearchParams; exports.useSubmit = _chunk4DGLNKXFjs.useSubmit; exports.useViewTransitionState = _chunk4DGLNKXFjs.useViewTransitionState;
