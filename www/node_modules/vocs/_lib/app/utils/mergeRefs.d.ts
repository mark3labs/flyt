import type { MutableRefObject, RefCallback } from 'react';
type MutableRefList<T> = Array<RefCallback<T> | MutableRefObject<T> | undefined | null>;
export declare function mergeRefs<T>(...refs: MutableRefList<T>): RefCallback<T>;
export declare function setRef<T>(val: T, ...refs: MutableRefList<T>): void;
export {};
//# sourceMappingURL=mergeRefs.d.js.map