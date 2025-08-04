import { jsx as _jsx, jsxs as _jsxs } from "react/jsx-runtime";
import { FloatingArrow, arrow, offset, safePolygon, shift, useFloating, useHover, useInteractions, } from '@floating-ui/react';
import { useRef, useState } from 'react';
import { primitiveColorVars } from '../../styles/vars.css.js';
export function TwoslashPopover({ children, ...props }) {
    const [popover, target] = children;
    const arrowRef = useRef(null);
    const [isOpen, setIsOpen] = useState(false);
    const { context, refs, floatingStyles } = useFloating({
        middleware: [
            arrow({
                element: arrowRef,
            }),
            offset(8),
            shift(),
        ],
        open: isOpen,
        onOpenChange: setIsOpen,
        placement: 'bottom-start',
    });
    const hover = useHover(context, { handleClose: safePolygon() });
    const { getReferenceProps, getFloatingProps } = useInteractions([hover]);
    const targetChildren = target.props.children;
    const popoverChildren = popover.props.children;
    return (_jsxs("span", { ...props, children: [_jsx("span", { className: "twoslash-target", ref: refs.setReference, ...getReferenceProps(), children: targetChildren }), isOpen && (_jsxs("div", { className: "twoslash-popup-info-hover", ref: refs.setFloating, style: floatingStyles, ...getFloatingProps(), children: [_jsx(FloatingArrow, { ref: arrowRef, context: context, fill: primitiveColorVars.background5, height: 3, stroke: primitiveColorVars.border2, strokeWidth: 1, width: 7 }), _jsx("div", { className: "twoslash-popup-scroll-container", children: popoverChildren })] }))] }));
}
//# sourceMappingURL=TwoslashPopover.js.map