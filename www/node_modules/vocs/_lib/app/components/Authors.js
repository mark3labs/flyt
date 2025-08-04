import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { Fragment, useMemo } from 'react';
import { usePageData } from '../hooks/usePageData.js';
import * as styles from './Authors.css.js';
export function Authors(props) {
    const { frontmatter } = usePageData();
    const { authors: authors_ = frontmatter?.authors, date = frontmatter?.date } = props;
    const authors = useMemo(() => {
        if (!authors_)
            return undefined;
        if (Array.isArray(authors_))
            return authors_;
        return authors_.split(',').map((author) => author.trim());
    }, [authors_]);
    const formattedDate = useMemo(() => {
        if (!date)
            return null;
        const dateObject = new Date(date);
        return dateObject.toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
        });
    }, [date]);
    return (_jsxs("div", { className: styles.root, children: [formattedDate, authors && (formattedDate ? ' by ' : 'By '), _jsx("span", { className: styles.authors, children: authors?.map((author, index) => {
                    const { text, url } = parseAuthor(author);
                    return (_jsxs(Fragment, { children: [url ? (_jsx("a", { className: styles.link, href: url, target: "_blank", rel: "noopener noreferrer", children: text })) : (text), index < authors.length - 2 && _jsx("span", { className: styles.separator, children: ", " }), index < authors.length - 1 && _jsx("span", { className: styles.separator, children: " & " })] }, index));
                }) })] }));
}
function parseAuthor(author) {
    const match = author.match(/\[(.+)\]\((.+)\)/);
    if (!match)
        return { text: author, url: undefined };
    return {
        text: match[1],
        url: match[2],
    };
}
//# sourceMappingURL=Authors.js.map