// src/client/utils.ts
var mergePath = (base, path) => {
  base = base.replace(/\/+$/, "");
  base = base + "/";
  path = path.replace(/^\/+/, "");
  return base + path;
};
var replaceUrlParam = (urlString, params) => {
  for (const [k, v] of Object.entries(params)) {
    const reg = new RegExp("/:" + k + "(?:{[^/]+})?\\??");
    urlString = urlString.replace(reg, v ? `/${v}` : "");
  }
  return urlString;
};
var buildSearchParams = (query) => {
  const searchParams = new URLSearchParams();
  for (const [k, v] of Object.entries(query)) {
    if (v === void 0) {
      continue;
    }
    if (Array.isArray(v)) {
      for (const v2 of v) {
        searchParams.append(k, v2);
      }
    } else {
      searchParams.set(k, v);
    }
  }
  return searchParams;
};
var replaceUrlProtocol = (urlString, protocol) => {
  switch (protocol) {
    case "ws":
      return urlString.replace(/^http/, "ws");
    case "http":
      return urlString.replace(/^ws/, "http");
  }
};
var removeIndexString = (urlSting) => {
  if (/^https?:\/\/[^\/]+?\/index$/.test(urlSting)) {
    return urlSting.replace(/\/index$/, "/");
  }
  return urlSting.replace(/\/index$/, "");
};
function isObject(item) {
  return typeof item === "object" && item !== null && !Array.isArray(item);
}
function deepMerge(target, source) {
  if (!isObject(target) && !isObject(source)) {
    return source;
  }
  const merged = { ...target };
  for (const key in source) {
    const value = source[key];
    if (isObject(merged[key]) && isObject(value)) {
      merged[key] = deepMerge(merged[key], value);
    } else {
      merged[key] = value;
    }
  }
  return merged;
}
export {
  buildSearchParams,
  deepMerge,
  mergePath,
  removeIndexString,
  replaceUrlParam,
  replaceUrlProtocol
};
