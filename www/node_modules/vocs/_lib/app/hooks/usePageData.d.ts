import type { Module } from '../types.js';
export declare function usePageData(): {
    content?: string;
    filePath?: string;
    frontmatter: Module["frontmatter"];
    lastUpdatedAt?: number;
    previousPath?: string;
};
export declare const PageDataContext: import("react").Context<{
    content?: string;
    filePath?: string;
    frontmatter: Module["frontmatter"];
    lastUpdatedAt?: number;
    previousPath?: string;
} | undefined>;
//# sourceMappingURL=usePageData.d.js.map