import { type ClassValue } from 'clsx';
import type { ReactNode } from 'react';
export type StepProps = {
    children: ReactNode;
    className?: ClassValue;
    title: ReactNode | string;
    titleLevel?: 2 | 3 | 4 | 5 | 6;
};
export declare function Step({ children, className, title, titleLevel }: StepProps): import("react/jsx-runtime").JSX.Element;
//# sourceMappingURL=Step.d.js.map