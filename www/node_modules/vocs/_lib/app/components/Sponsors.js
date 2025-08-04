import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { assignInlineVars } from '@vanilla-extract/dynamic';
import clsx from 'clsx';
import { Fragment } from 'react';
import { useConfig } from '../hooks/useConfig.js';
import { Link } from './Link.js';
import * as styles from './Sponsors.css.js';
export function Sponsors() {
    const { sponsors } = useConfig();
    return (_jsx("div", { className: styles.root, children: sponsors?.map((sponsorSet, i) => (_jsxs(Fragment, { children: [_jsx("div", { className: styles.title, children: sponsorSet.name }), sponsorSet.items.map((sponsorRow, i) => (_jsx("div", { className: styles.row, style: assignInlineVars({
                        [styles.columnsVar]: sponsorRow.length.toString(),
                        [styles.heightVar]: `${sponsorSet.height?.toString() ?? '40'}px`,
                    }), children: sponsorRow.map((sponsor, i) => (_jsx(Link, { className: clsx(styles.column, sponsor ? styles.sponsor : undefined), href: sponsor?.link, hideExternalIcon: true, variant: "styleless", children: sponsor?.image && (_jsx("img", { className: styles.image, src: sponsor?.image, alt: sponsor?.name })) }, i))) }, i)))] }, i))) }));
}
//# sourceMappingURL=Sponsors.js.map