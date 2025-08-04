import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import * as Dialog from '@radix-ui/react-dialog';
import { MagnifyingGlassIcon } from '@radix-ui/react-icons';
import { useEffect, useState } from 'react';
import { useSearchIndex } from '../hooks/useSearchIndex.js';
import * as styles from './DesktopSearch.css.js';
import { SearchDialog } from './SearchDialog.js';
export function DesktopSearch() {
    useSearchIndex();
    const [open, setOpen] = useState(false);
    useEffect(() => {
        function keyDownHandler(event) {
            const isInput = document.activeElement instanceof HTMLElement &&
                (['input', 'select', 'textarea'].includes(document.activeElement.tagName.toLowerCase()) ||
                    document.activeElement.isContentEditable);
            if (event.key === '/' && !open && !isInput) {
                event.preventDefault();
                setOpen(true);
            }
            else if (event.metaKey === true && event.key === 'k') {
                event.preventDefault();
                setOpen((x) => !x);
            }
        }
        window.addEventListener('keydown', keyDownHandler);
        return () => {
            window.removeEventListener('keydown', keyDownHandler);
        };
    }, [open]);
    return (_jsxs(Dialog.Root, { open: open, onOpenChange: setOpen, children: [_jsx(Dialog.Trigger, { asChild: true, children: _jsxs("button", { className: styles.search, type: "button", children: [_jsx(MagnifyingGlassIcon, { style: { marginTop: 2 } }), "Search...", _jsx("div", { className: styles.searchCommand, children: _jsx("div", { style: {
                                    background: 'currentColor',
                                    transform: 'rotate(45deg)',
                                    width: 1.5,
                                    borderRadius: 2,
                                    height: '100%',
                                } }) })] }) }), _jsx(SearchDialog, { open: open, onClose: () => setOpen(false) })] }));
}
//# sourceMappingURL=DesktopSearch.js.map