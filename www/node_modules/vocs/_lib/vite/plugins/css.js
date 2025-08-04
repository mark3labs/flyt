import { default as autoprefixer } from 'autoprefixer';
export function css() {
    return {
        name: 'css',
        async config() {
            return {
                css: {
                    postcss: {
                        plugins: [autoprefixer()].filter(Boolean),
                    },
                },
            };
        },
    };
}
//# sourceMappingURL=css.js.map