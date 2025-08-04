import { jsx as _jsx } from "react/jsx-runtime";
import { cloneElement } from 'react';
import * as stepStyles from '../Step.css.js';
import { Step } from '../Step.js';
import { Steps as Steps_ } from '../Steps.js';
export function Steps({ children }) {
    if (!Array.isArray(children))
        return null;
    return (_jsx(Steps_, { children: children.map(({ props }, i) => {
            const [title, ...children] = Array.isArray(props.children)
                ? props.children
                : [props.children];
            return (_jsx(Step, { title: cloneElement(title, { className: stepStyles.title }), children: children }, i));
        }) }));
}
//# sourceMappingURL=Steps.js.map