import { jsx as _jsx } from "react/jsx-runtime";
import { clsx } from 'clsx';
import { forwardRef } from 'react';
import { useLocation } from 'react-router';
import { ExternalLink } from './ExternalLink.js';
import * as styles from './Link.css.js';
import { RouterLink } from './RouterLink.js';
export const Link = forwardRef((props, ref) => {
    const { hideExternalIcon, href, variant = 'accent' } = props;
    const { pathname } = useLocation();
    // External links
    if (href?.match(/^(www|https?)/))
        return (_jsx(ExternalLink, { ...props, ref: ref, className: clsx(props.className, styles.root, variant === 'accent' && styles.accent, variant === 'styleless' && styles.styleless), hideExternalIcon: hideExternalIcon }));
    // Internal links
    const [before, after] = (href || '').split('#');
    const to = `${before ? before : pathname}${after ? `#${after}` : ''}`;
    return (_jsx(RouterLink, { ...props, ref: ref, className: clsx(props.className, styles.root, variant === 'accent' && styles.accent, variant === 'styleless' && styles.styleless), to: to }));
});
//# sourceMappingURL=Link.js.map