import React, {} from 'react';
export function deserializeElement(element, key) {
    if (typeof element !== 'object')
        return element;
    if (element === null)
        return element;
    if (Array.isArray(element))
        return element.map((el, i) => deserializeElement(el, i));
    const props = element.props.children
        ? { ...element.props, children: deserializeElement(element.props.children) }
        : element.props;
    return React.createElement(element.type, { ...props, key });
}
//# sourceMappingURL=deserializeElement.js.map