import { jsx as _jsx } from "react/jsx-runtime";
import { clsx } from 'clsx';
import { useLayout } from '../../hooks/useLayout.js';
import { Authors } from '../Authors.js';
import { BlogPosts } from '../BlogPosts.js';
import { Sponsors } from '../Sponsors.js';
import { AutolinkIcon } from './AutolinkIcon.js';
import { CodeGroup } from './CodeGroup.js';
import * as styles from './Div.css.js';
import { Steps } from './Steps.js';
import { Subtitle } from './Subtitle.js';
export function Div(props) {
    const { layout } = useLayout();
    const className = clsx(props.className, styles.root);
    if (props.className === 'code-group')
        return _jsx(CodeGroup, { ...props, className: className });
    if ('data-authors' in props)
        return _jsx(Authors, {});
    if ('data-blog-posts' in props)
        return _jsx(BlogPosts, {});
    if ('data-sponsors' in props)
        return _jsx(Sponsors, {});
    if ('data-autolink-icon' in props && layout === 'docs')
        return _jsx(AutolinkIcon, { ...props, className: className });
    if ('data-vocs-steps' in props)
        return _jsx(Steps, { ...props, className: className });
    if (props.role === 'doc-subtitle')
        return _jsx(Subtitle, { ...props });
    return _jsx("div", { ...props, className: className });
}
//# sourceMappingURL=Div.js.map