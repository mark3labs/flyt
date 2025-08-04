import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import * as Dialog from '@radix-ui/react-dialog';
import { MagnifyingGlassIcon } from '@radix-ui/react-icons';
import { useState } from 'react';
import * as styles from './MobileSearch.css.js';
import { SearchDialog } from './SearchDialog.js';
export function MobileSearch() {
    const [open, setOpen] = useState(false);
    return (_jsxs(Dialog.Root, { open: open, onOpenChange: setOpen, children: [_jsx(Dialog.Trigger, { asChild: true, children: _jsx("button", { className: styles.searchButton, type: "button", "aria-label": "Search", children: _jsx(MagnifyingGlassIcon, { height: 21, width: 21 }) }) }), _jsx(SearchDialog, { open: open, onClose: () => setOpen(false) })] }));
}
//# sourceMappingURL=MobileSearch.js.map