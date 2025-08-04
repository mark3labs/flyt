import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { runSync } from '@mdx-js/mdx';
import { Cross1Icon } from '@radix-ui/react-icons';
import { assignInlineVars } from '@vanilla-extract/dynamic';
import clsx from 'clsx';
import { Fragment, useMemo } from 'react';
import * as runtime from 'react/jsx-runtime';
import { useConfig } from '../hooks/useConfig.js';
import { deserializeElement } from '../utils/deserializeElement.js';
import * as styles from './Banner.css.js';
export function Banner({ hide }) {
    const { banner } = useConfig();
    const ConsumerBanner = useMemo(() => {
        const content = banner?.content ?? '';
        if (!content)
            return null;
        if (typeof content !== 'string')
            return () => deserializeElement(content);
        const { default: MDXBanner } = runSync(content, { ...runtime, Fragment });
        return MDXBanner;
    }, [banner]);
    if (!ConsumerBanner)
        return null;
    return (_jsx("div", { className: clsx(styles.root), style: assignInlineVars({
            [styles.bannerBackgroundColor]: banner?.backgroundColor,
            [styles.bannerTextColor]: banner?.textColor,
        }), children: _jsxs("div", { className: clsx(styles.inner), children: [_jsx("div", { className: clsx(styles.content), children: _jsx(ConsumerBanner, {}) }), banner?.dismissable !== 'false' && (_jsx("button", { className: clsx(styles.closeButton), onClick: hide, type: "button", children: _jsx(Cross1Icon, { width: 14, height: 14 }) }))] }) }));
}
//# sourceMappingURL=Banner.js.map