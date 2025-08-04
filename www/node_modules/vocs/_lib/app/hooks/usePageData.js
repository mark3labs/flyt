import { createContext, useContext } from 'react';
export function usePageData() {
    const pageData = useContext(PageDataContext);
    if (!pageData)
        throw new Error('`usePageData` must be used within `PageDataContext.Provider`.');
    return pageData;
}
export const PageDataContext = createContext(undefined);
//# sourceMappingURL=usePageData.js.map