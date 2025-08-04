export declare const search: ({ cacheDir, }?: {
    cacheDir?: string;
}) => {
    get(k: string): any;
    set<v>(k: string, value: v): void;
};
export declare const twoslash: ({ cacheDir, }?: {
    cacheDir?: string;
}) => {
    get(k: string): any;
    set<v>(k: string, value: v): void;
};
export declare function create(key: string): ({ cacheDir, }?: {
    cacheDir?: string;
}) => {
    get(k: string): any;
    set<v>(k: string, value: v): void;
};
export declare function clear({ cacheDir, }?: {
    cacheDir?: string;
}): void;
//# sourceMappingURL=cache.d.js.map