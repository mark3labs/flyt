import type { ReactNode } from 'react';
import { type ButtonProps } from './Button.js';
export type HomePageProps = {
    description?: ReactNode;
    tagline?: ReactNode;
};
export declare function Root({ children, className }: {
    children: ReactNode;
    className?: string;
}): import("react/jsx-runtime").JSX.Element;
export declare function Logo({ className }: {
    className?: string;
}): import("react/jsx-runtime").JSX.Element;
export declare function Tagline({ children, className }: {
    children: ReactNode;
    className?: string;
}): import("react/jsx-runtime").JSX.Element;
export declare function Description({ children, className }: {
    children: ReactNode;
    className?: string;
}): import("react/jsx-runtime").JSX.Element;
export declare function Buttons({ children, className }: {
    children: ReactNode;
    className?: string;
}): import("react/jsx-runtime").JSX.Element;
export declare function Button(props: ButtonProps): import("react/jsx-runtime").JSX.Element;
export declare function InstallPackage({ name, type, }: {
    children: ReactNode;
    className?: string;
    name: string;
    type?: 'install' | 'init';
}): import("react/jsx-runtime").JSX.Element;
//# sourceMappingURL=HomePage.d.js.map