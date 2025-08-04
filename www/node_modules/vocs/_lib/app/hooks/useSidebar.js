import { useMemo } from 'react';
import { useLocation } from 'react-router';
import { useConfig } from './useConfig.js';
export function useSidebar() {
    const { pathname } = useLocation();
    const config = useConfig();
    const { sidebar } = config;
    if (!sidebar)
        return { items: [] };
    if (Array.isArray(sidebar))
        return { items: sidebar };
    const sidebarKey = useMemo(() => {
        const keys = Object.keys(sidebar).filter((key) => pathname.startsWith(key));
        return keys[keys.length - 1];
    }, [sidebar, pathname]);
    if (!sidebarKey)
        return { items: [] };
    if (Array.isArray(sidebar[sidebarKey]))
        return { key: sidebarKey, items: sidebar[sidebarKey] };
    return { ...sidebar[sidebarKey], key: sidebarKey };
}
//# sourceMappingURL=useSidebar.js.map