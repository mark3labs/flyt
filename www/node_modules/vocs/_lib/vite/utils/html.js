import { jsx as _jsx, jsxs as _jsxs, Fragment as _Fragment } from "react/jsx-runtime";
import { renderToStaticMarkup, renderToString } from 'react-dom/server';
export async function toMarkup(parameters) {
    const { body, config, location, template } = parameters;
    const { theme } = config;
    const preHtml = renderToStaticMarkup(_jsxs("html", { lang: "en", children: [_jsx("head", {}), _jsx("body", { children: body })] }));
    const preHead = preHtml.match(/<head>([\s\S]*?)<\/head>/)?.[1];
    const configHead = await (async () => {
        if (typeof config.head === 'function')
            return await config.head({ path: location });
        if (typeof config.head === 'object') {
            const entry = Object.entries(config.head)
                .reverse()
                .find(([key]) => location.startsWith(key));
            return entry?.[1];
        }
        return config.head;
    })();
    const head = renderToString(_jsxs(_Fragment, { children: [parameters.head, configHead] }));
    let html = template.replace('<!--head-->', `${head}${preHead}`);
    html = html.replace('<!--body-->', renderToString(body).replace(preHead ?? '', ''));
    const match = html.match(/property="og:image" content="(.*)"/);
    if (match?.[1]) {
        html = html.replace(/property="og:image" content="(.*)"/, `property="og:image" content="${match[1].replace(/&amp;/g, '&')}"`);
    }
    if (theme?.colorScheme && theme?.colorScheme !== 'system')
        html = html.replace('lang="en"', `lang="en" class="${theme.colorScheme}"`);
    return html;
}
//# sourceMappingURL=html.js.map