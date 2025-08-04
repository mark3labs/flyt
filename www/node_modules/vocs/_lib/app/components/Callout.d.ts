import { type ClassValue } from 'clsx';
import type { ReactNode } from 'react';
export type CalloutProps = {
    className: ClassValue;
    children: ReactNode;
    type: 'note' | 'info' | 'warning' | 'danger' | 'tip' | 'success';
};
export declare function Callout({ className, children, type }: CalloutProps): import("react/jsx-runtime").JSX.Element;
//# sourceMappingURL=Callout.d.js.map