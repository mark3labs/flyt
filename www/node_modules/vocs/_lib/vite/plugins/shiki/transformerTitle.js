const titleRegex = /title="(.*)"|\[(.*)\]/;
export const transformerTitle = () => ({
    name: 'title',
    root(hast) {
        const titleMatch = this.options.meta?.__raw?.match(titleRegex);
        if (!titleMatch)
            return;
        const title = titleMatch[1] || titleMatch[2];
        const child = hast.children[0];
        hast.children = [
            {
                ...child,
                properties: {
                    ...child.properties,
                    'data-title': title,
                    'data-lang': this.options.lang,
                },
            },
        ];
    },
});
//# sourceMappingURL=transformerTitle.js.map