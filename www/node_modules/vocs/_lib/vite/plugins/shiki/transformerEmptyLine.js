export const transformerEmptyLine = () => ({
    name: 'empty-line',
    line(hast) {
        const child = hast.children[0];
        if (child)
            return;
        hast.properties['data-empty-line'] = true;
        hast.children = [
            {
                type: 'text',
                value: ' ',
            },
        ];
    },
});
//# sourceMappingURL=transformerEmptyLine.js.map