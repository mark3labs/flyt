import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { clsx } from 'clsx';
import * as styles from './Step.css.js';
import { H2 } from './mdx/H2.js';
import { H3 } from './mdx/H3.js';
import { H4 } from './mdx/H4.js';
import { H5 } from './mdx/H5.js';
import { H6 } from './mdx/H6.js';
export function Step({ children, className, title, titleLevel = 2 }) {
    const Element = (() => {
        if (titleLevel === 2)
            return H2;
        if (titleLevel === 3)
            return H3;
        if (titleLevel === 4)
            return H4;
        if (titleLevel === 5)
            return H5;
        if (titleLevel === 6)
            return H6;
        throw new Error('Invalid.');
    })();
    return (_jsxs("div", { className: clsx(className, styles.root), children: [typeof title === 'string' ? _jsx(Element, { className: styles.title, children: title }) : title, _jsx("div", { className: styles.content, children: children })] }));
}
//# sourceMappingURL=Step.js.map