import { selectAll } from "hast-util-select";
import { classnames } from 'hast-util-classnames';
const rehypeClassNames = (additions) => {
    return (tree) => {
        if (additions)
            Object.entries(additions).map(([selector, cName]) => {
                return (nodes) => selectAll(selector, nodes).forEach((elem) => {
                    classnames(elem, cName);
                });
            }).forEach((a) => a(tree));
    };
};
export default rehypeClassNames;
