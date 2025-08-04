import { jsx as _jsx } from "react/jsx-runtime";
import { clsx } from 'clsx';
import { assignInlineVars } from '@vanilla-extract/dynamic';
import { forwardRef } from 'react';
import { useConfig } from '../hooks/useConfig.js';
import * as styles from './ExternalLink.css.js';
export const ExternalLink = forwardRef(({ className, children, hideExternalIcon, href, ...props }, ref) => {
    const { basePath } = useConfig();
    const assetBasePath = import.meta.env.PROD ? basePath : '';
    return (_jsx("a", { ref: ref, className: clsx(className, hideExternalIcon || typeof children !== 'string' ? undefined : styles.root), href: href, target: "_blank", rel: "noopener noreferrer", style: assignInlineVars({
            [styles.iconUrl]: `url(${assetBasePath}/.vocs/icons/arrow-diagonal.svg)`,
        }), ...props, children: children }));
});
//# sourceMappingURL=ExternalLink.js.map