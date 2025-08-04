import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { clsx } from 'clsx';
import { useMemo, } from 'react';
import { useCopyCode } from '../../hooks/useCopyCode.js';
import { CopyButton } from '../CopyButton.js';
import { CodeBlock } from './CodeBlock.js';
import { CodeTitle } from './CodeTitle.js';
import * as styles from './Pre.css.js';
export function Pre({ children, className, ...props }) {
    const { copied, copy, ref } = useCopyCode();
    function recurseChildren(children) {
        if (!children)
            return children;
        if (typeof children !== 'object')
            return children;
        if ('props' in children)
            return {
                ...children,
                props: {
                    ...children.props,
                    children: Array.isArray(children.props.children)
                        ? children.props.children.map(recurseChildren)
                        : recurseChildren(children.props.children),
                },
            };
        return children;
    }
    // biome-ignore lint/correctness/useExhaustiveDependencies:
    const children_ = useMemo(() => recurseChildren(children), [children]);
    const wrap = (children) => {
        if (className?.includes('shiki'))
            return (_jsxs(CodeBlock, { children: [props['data-title'] && (_jsx(CodeTitle, { language: props['data-lang'], children: props['data-title'] })), children] }));
        return children;
    };
    return wrap(_jsx("div", { className: clsx(styles.wrapper), children: _jsxs("pre", { ref: ref, ...props, className: clsx(className, styles.root), children: [_jsx(CopyButton, { copied: copied, copy: copy }), children_] }) }));
}
//# sourceMappingURL=Pre.js.map