import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import clsx from 'clsx';
import { DropdownMenu } from 'radix-ui';
import { useCallback, useEffect, useMemo, useState } from 'react';
import { useLocation } from 'react-router';
import { useConfig } from '../hooks/useConfig.js';
import { usePageData } from '../hooks/usePageData.js';
import * as styles from './AiCtaDropdown.css.js';
import * as buttonStyles from './Button.css.js';
import { Link } from './Link.js';
import { CheckCircle } from './icons/CheckCircle.js';
import { ChevronDown } from './icons/ChevronDown.js';
import { Copy } from './icons/Copy.js';
import { OpenAi } from './icons/OpenAi.js';
export function AiCtaDropdown() {
    const { content } = usePageData();
    const { aiCta } = useConfig();
    const location = useLocation();
    const [copied, setCopied] = useState(false);
    useEffect(() => {
        if (!copied)
            return;
        const timeout = setTimeout(() => setCopied(false), 1000);
        return () => clearTimeout(timeout);
    }, [copied]);
    const copy = useCallback(() => {
        setCopied(true);
        navigator.clipboard.writeText(content ?? '');
    }, [content]);
    const query = useMemo(() => {
        const href = window.location.origin + location.pathname;
        if (typeof aiCta === 'object')
            return aiCta.query({ location: href });
        return `Please research and analyze this page: ${href} so I can ask you questions about it. Once you have read it, prompt me with any questions I have. Do not post content from the page in your response. Any of my follow up questions must reference the site I gave you.`;
    }, [aiCta, location.pathname]);
    return (_jsxs("div", { className: styles.root, children: [copied ? (_jsxs("div", { className: clsx(buttonStyles.button, styles.buttonLeft), children: [_jsx("div", { style: { width: '14px', height: '14px' }, children: _jsx(CheckCircle, {}) }), "Copied"] })) : (_jsxs(Link, { className: clsx(buttonStyles.button, styles.buttonLeft), href: `https://chatgpt.com?hints=search&q=${encodeURIComponent(query)}`, variant: "styleless", children: [_jsx("div", { style: { width: '14px', height: '14px' }, children: _jsx(OpenAi, {}) }), "Ask in ChatGPT"] })), _jsxs(DropdownMenu.Root, { children: [_jsx(DropdownMenu.Trigger, { asChild: true, children: _jsx("button", { className: clsx(buttonStyles.button, styles.buttonRight), type: "button", children: _jsx("div", { style: { width: '14px', height: '14px' }, children: _jsx(ChevronDown, {}) }) }) }), _jsx(DropdownMenu.Portal, { children: _jsx(DropdownMenu.Content, { align: "end", className: styles.dropdownMenuContent, sideOffset: 4, children: _jsxs(DropdownMenu.Item, { className: styles.dropdownMenuItem, onClick: copy, children: [_jsx("div", { style: { width: '14px', height: '14px' }, children: _jsx(Copy, {}) }), "Copy page for LLMs"] }) }) })] })] }));
}
//# sourceMappingURL=AiCtaDropdown.js.map