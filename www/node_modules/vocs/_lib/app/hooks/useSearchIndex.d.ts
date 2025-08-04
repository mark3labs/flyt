import MiniSearch from 'minisearch';
export type Result = {
    href: string;
    html: string;
    isPage: boolean;
    text?: string;
    title: string;
    titles: string[];
};
export declare function useSearchIndex(): MiniSearch<Result> | undefined;
//# sourceMappingURL=useSearchIndex.d.js.map