import { addClassToHast } from 'shiki';
function escapeRegExp(str) {
    return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}
export function transformerNotationMap(options = {}, name = '@shikijs/transformers:notation-map') {
    const { classMap = {}, classActivePre = undefined } = options;
    return createCommentNotationTransformer(name, new RegExp(`\\s*(?://|/\\*|<!--|#)\\s+\\[!code (${Object.keys(classMap).map(escapeRegExp).join('|')})(:\\d+)?\\]\\s*(?:\\*/|-->)?`), function ([_, match, range = ':1'], _line, _comment, lines, index) {
        const lineNum = Number.parseInt(range.slice(1), 10);
        // biome-ignore lint/complexity/noForEach:
        lines.slice(index, index + lineNum).forEach((line) => {
            addClassToHast(line, classMap[match]);
        });
        if (classActivePre)
            addClassToHast(this.pre, classActivePre);
        return true;
    });
}
export function createCommentNotationTransformer(name, regex, onMatch, removeEmptyLines = false) {
    return {
        name,
        code(code) {
            const lines = code.children.filter((i) => i.type === 'element');
            const linesToRemove = [];
            lines.forEach((line, idx) => {
                let nodeToRemove;
                for (const child of line.children) {
                    if (child.type !== 'element')
                        continue;
                    const text = child.children[0];
                    if (text.type !== 'text')
                        continue;
                    let replaced = false;
                    text.value = text.value.replace(regex, (...match) => {
                        if (onMatch.call(this, match, line, child, lines, idx)) {
                            replaced = true;
                            return '';
                        }
                        return match[0];
                    });
                    if (replaced && !text.value.trim())
                        nodeToRemove = child;
                }
                if (nodeToRemove) {
                    line.children.splice(line.children.indexOf(nodeToRemove), 1);
                    // Remove if empty
                    if (line.children.length === 0) {
                        linesToRemove.push(line);
                        if (removeEmptyLines) {
                            const next = code.children[code.children.indexOf(line) + 1];
                            if (next && next.type === 'text' && next.value === '\n')
                                linesToRemove.push(next);
                        }
                    }
                }
            });
            for (const line of linesToRemove)
                code.children.splice(code.children.indexOf(line), 1);
        },
    };
}
export function highlightWordInLine(line, ignoredElement, word, className) {
    line.children = line.children.flatMap((span) => {
        if (span.type !== 'element' || span.tagName !== 'span' || span === ignoredElement)
            return span;
        const textNode = span.children[0];
        if (textNode.type !== 'text')
            return span;
        return replaceSpan(span, textNode.value, word, className) ?? span;
    });
}
function inheritElement(original, overrides) {
    return {
        ...original,
        properties: {
            ...original.properties,
        },
        ...overrides,
    };
}
function replaceSpan(span, text, word, className) {
    const index = text.indexOf(word);
    if (index === -1)
        return;
    const createNode = (value) => inheritElement(span, {
        children: [
            {
                type: 'text',
                value,
            },
        ],
    });
    const nodes = [];
    if (index > 0)
        nodes.push(createNode(text.slice(0, index)));
    const highlightedNode = createNode(word);
    addClassToHast(highlightedNode, className);
    nodes.push(highlightedNode);
    if (index + word.length < text.length)
        nodes.push(createNode(text.slice(index + word.length)));
    return nodes;
}
//# sourceMappingURL=utils.js.map