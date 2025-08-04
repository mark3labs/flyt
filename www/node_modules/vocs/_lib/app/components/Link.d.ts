type LinkProps = {
    children: React.ReactNode;
    className?: string;
    hideExternalIcon?: boolean;
    onClick?: (e: React.MouseEvent<HTMLAnchorElement>) => void;
    href?: string;
    variant?: 'accent' | 'styleless';
};
export declare const Link: import("react").ForwardRefExoticComponent<LinkProps & import("react").RefAttributes<unknown>>;
export {};
//# sourceMappingURL=Link.d.js.map