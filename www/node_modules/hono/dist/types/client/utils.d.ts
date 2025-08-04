export declare const mergePath: (base: string, path: string) => string;
export declare const replaceUrlParam: (urlString: string, params: Record<string, string | undefined>) => string;
export declare const buildSearchParams: (query: Record<string, string | string[]>) => URLSearchParams;
export declare const replaceUrlProtocol: (urlString: string, protocol: "ws" | "http") => string;
export declare const removeIndexString: (urlSting: string) => string;
export declare function deepMerge<T>(target: T, source: Record<string, unknown>): T;
