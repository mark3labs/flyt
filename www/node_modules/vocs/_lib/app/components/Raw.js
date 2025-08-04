import { jsx as _jsx } from "react/jsx-runtime";
import { MDXProvider } from '@mdx-js/react';
export function Raw({ children }) {
    return _jsx(MDXProvider, { disableParentContext: true, children: children });
}
//# sourceMappingURL=Raw.js.map