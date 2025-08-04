import { jsx as _jsx } from "react/jsx-runtime";
import { clsx } from 'clsx';
import { useLocation } from 'react-router';
import { Link } from '../Link.js';
import * as styles from './Anchor.css.js';
import { Autolink } from './Autolink.js';
export function Anchor(props) {
    const { children, href } = props;
    const { pathname } = useLocation();
    // Heading slug links
    if (children &&
        typeof children === 'object' &&
        'props' in children &&
        children.props['data-autolink-icon'])
        return _jsx(Autolink, { className: clsx(props.className, styles.root), ...props });
    // ID links
    if (href?.match(/^#/))
        return (_jsx("a", { className: clsx(props.className, styles.root), ...props, href: `${pathname}${href}` }));
    return _jsx(Link, { className: clsx(props.className, styles.root), ...props, hideExternalIcon: true });
}
//# sourceMappingURL=Anchor.js.map