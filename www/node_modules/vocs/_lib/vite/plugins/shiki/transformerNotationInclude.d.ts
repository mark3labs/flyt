import type { ShikiTransformer } from 'shiki';
export type TransformerNotationIncludeOptions = {
    rootDir: string;
};
export declare const transformerNotationInclude: ({ rootDir, }: TransformerNotationIncludeOptions) => ShikiTransformer;
export declare function processIncludes({ code, getSource, }: {
    code: string;
    getSource: (fileName: string) => string | undefined;
}): string;
//# sourceMappingURL=transformerNotationInclude.d.js.map