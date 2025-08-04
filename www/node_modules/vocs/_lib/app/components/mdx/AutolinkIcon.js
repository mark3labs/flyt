import { jsx as _jsx } from "react/jsx-runtime";
import { assignInlineVars } from '@vanilla-extract/dynamic';
import { clsx } from 'clsx';
import { useConfig } from '../../hooks/useConfig.js';
import * as styles from './AutolinkIcon.css.js';
export function AutolinkIcon(props) {
    const { basePath } = useConfig();
    const assetBasePath = import.meta.env.PROD ? basePath : '';
    return (_jsx("div", { ...props, className: clsx(props.className, styles.root), style: assignInlineVars({
            [styles.iconUrl]: `url(${assetBasePath}/.vocs/icons/link.svg)`,
        }) }));
}
//# sourceMappingURL=AutolinkIcon.js.map