import * as NavigationMenu from '@radix-ui/react-navigation-menu';
import type { ReactNode } from 'react';
export declare const Root: (props: NavigationMenu.NavigationMenuProps) => import("react/jsx-runtime").JSX.Element;
export declare const List: (props: NavigationMenu.NavigationMenuListProps) => import("react/jsx-runtime").JSX.Element;
export declare const Link: ({ active, children, className, href, }: {
    active?: boolean;
    children: ReactNode;
    className?: string;
    href?: string;
}) => import("react/jsx-runtime").JSX.Element;
export declare const Item: (props: NavigationMenu.NavigationMenuItemProps) => import("react/jsx-runtime").JSX.Element;
export declare const Trigger: ({ active, className, ...props }: NavigationMenu.NavigationMenuTriggerProps & {
    active?: boolean;
}) => import("react/jsx-runtime").JSX.Element;
export declare const Content: (props: NavigationMenu.NavigationMenuContentProps) => import("react/jsx-runtime").JSX.Element;
//# sourceMappingURL=NavigationMenu.d.js.map