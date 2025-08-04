import type { Element } from 'hast';
import { type ShikiTransformer, type ShikiTransformerContext } from 'shiki';
export interface TransformerNotationMapOptions {
    classMap?: Record<string, string | string[]>;
    /**
     * Class added to the <pre> element when the current code has diff
     */
    classActivePre?: string;
}
export declare function transformerNotationMap(options?: TransformerNotationMapOptions, name?: string): ShikiTransformer;
export declare function createCommentNotationTransformer(name: string, regex: RegExp, onMatch: (this: ShikiTransformerContext, match: string[], line: Element, commentNode: Element, lines: Element[], index: number) => boolean, removeEmptyLines?: boolean): ShikiTransformer;
export declare function highlightWordInLine(line: Element, ignoredElement: Element | null, word: string, className: string): void;
//# sourceMappingURL=utils.d.js.map