export const transformerLineNumbers = () => ({
    name: 'line-numbers',
    code(hast) {
        if (!this.options.meta?.__raw?.includes('showLineNumbers'))
            return;
        hast.properties['data-line-numbers'] = true;
    },
});
//# sourceMappingURL=transformerLineNumbers.js.map