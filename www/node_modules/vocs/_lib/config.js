export async function defineConfig({ aiCta = true, blogDir = './pages/blog', cacheDir, head, ogImageUrl, rootDir = 'docs', title = 'Docs', titleTemplate = `%s – ${title}`, ...config }) {
    const basePath = parseBasePath(config.basePath);
    return {
        aiCta,
        blogDir,
        cacheDir,
        head,
        ogImageUrl,
        rootDir,
        title,
        titleTemplate,
        ...config,
        basePath,
        banner: await parseBanner(config.banner ?? ''),
        font: parseFont(config.font ?? {}),
        iconUrl: parseImageUrl(config.iconUrl, {
            basePath,
        }),
        logoUrl: parseImageUrl(config.logoUrl, {
            basePath,
        }),
        markdown: parseMarkdown(config.markdown ?? {}),
        socials: parseSocials(config.socials ?? []),
        topNav: parseTopNav(config.topNav ?? []),
        theme: await parseTheme(config.theme ?? {}),
        vite: parseViteConfig(config.vite, {
            basePath,
        }),
    };
}
export const getDefaultConfig = async () => await defineConfig({});
//////////////////////////////////////////////////////
// Parsers
function parseBasePath(basePath_) {
    let basePath = basePath_;
    if (!basePath)
        return '';
    if (!basePath.startsWith('/'))
        basePath = `/${basePath}`;
    if (basePath.endsWith('/'))
        basePath = basePath.slice(0, -1);
    return basePath;
}
async function parseBanner(banner) {
    if (!banner)
        return undefined;
    const bannerContent = (() => {
        if (typeof banner === 'string')
            return banner;
        if (typeof banner === 'object' && 'content' in banner)
            return banner.content;
        return undefined;
    })();
    const content = await (async () => {
        if (typeof bannerContent !== 'string')
            return bannerContent;
        const { compile } = await import('@mdx-js/mdx');
        const remarkGfm = (await import('remark-gfm')).default;
        return String(await compile(bannerContent, {
            outputFormat: 'function-body',
            remarkPlugins: [remarkGfm],
        }));
    })();
    if (!content)
        return undefined;
    const textColor = await (async () => {
        if (typeof banner === 'string')
            return undefined;
        if (typeof banner === 'object') {
            if ('textColor' in banner)
                return banner.textColor;
            if ('backgroundColor' in banner && banner.backgroundColor) {
                const chroma = (await import('chroma-js')).default;
                return chroma.contrast(banner.backgroundColor, 'white') < 4.5 ? 'black' : 'white';
            }
        }
        return undefined;
    })();
    return {
        height: '32px',
        ...(typeof banner === 'object' ? banner : {}),
        content,
        textColor,
    };
}
function parseFont(font) {
    if ('google' in font)
        return { default: font };
    return font;
}
function parseImageUrl(imageUrl, { basePath }) {
    if (!imageUrl)
        return;
    if (process.env.NODE_ENV === 'development')
        return imageUrl;
    if (typeof imageUrl === 'string') {
        if (imageUrl.startsWith('http'))
            return imageUrl;
        return `${basePath}${imageUrl}`;
    }
    return {
        dark: imageUrl.dark.startsWith('http') ? imageUrl.dark : `${basePath}${imageUrl.dark}`,
        light: imageUrl.light.startsWith('http') ? imageUrl.light : `${basePath}${imageUrl.light}`,
    };
}
function parseMarkdown(markdown) {
    return {
        ...markdown,
        code: {
            themes: {
                dark: 'github-dark-dimmed',
                light: 'github-light',
            },
            ...markdown.code,
        },
    };
}
const socialsMeta = {
    discord: { label: 'Discord', type: 'discord' },
    github: { label: 'GitHub', type: 'github' },
    telegram: { label: 'Telegram', type: 'telegram' },
    warpcast: { label: 'Warpcast', type: 'warpcast' },
    x: { label: 'X (Twitter)', type: 'x' },
};
function parseSocials(socials) {
    return socials.map((social) => {
        return {
            icon: social.icon,
            link: social.link,
            ...socialsMeta[social.icon],
        };
    });
}
let id = 0;
function parseTopNav(topNav) {
    const parsedTopNav = [];
    for (const item of topNav) {
        parsedTopNav.push({
            ...item,
            id: id++,
            items: item.items ? parseTopNav(item.items) : [],
        });
    }
    return parsedTopNav;
}
async function parseTheme(theme) {
    const chroma = (await import('chroma-js')).default;
    const accentColor = (() => {
        if (!theme.accentColor)
            return theme.accentColor;
        if (typeof theme.accentColor === 'object' &&
            !Object.keys(theme.accentColor).includes('light') &&
            !Object.keys(theme.accentColor).includes('dark'))
            return theme.accentColor;
        const accentColor = theme.accentColor;
        const accentColorLight = typeof accentColor === 'object' ? accentColor.light : accentColor;
        const accentColorDark = typeof accentColor === 'object' ? accentColor.dark : accentColor;
        return {
            backgroundAccent: {
                dark: accentColorDark,
                light: accentColorLight,
            },
            backgroundAccentHover: {
                dark: chroma(accentColorDark).darken(0.25).hex(),
                light: chroma(accentColorLight).darken(0.25).hex(),
            },
            backgroundAccentText: {
                dark: chroma.contrast(accentColorDark, 'white') < 4.5 ? 'black' : 'white',
                light: chroma.contrast(accentColorLight, 'white') < 4.5 ? 'black' : 'white',
            },
            borderAccent: {
                dark: chroma(accentColorDark).brighten(0.5).hex(),
                light: chroma(accentColorLight).darken(0.25).hex(),
            },
            textAccent: {
                dark: accentColorDark,
                light: accentColorLight,
            },
            textAccentHover: {
                dark: chroma(accentColorDark).darken(0.5).hex(),
                light: chroma(accentColorLight).darken(0.5).hex(),
            },
        };
    })();
    return {
        ...theme,
        accentColor,
    };
}
export function parseViteConfig(viteConfig, { basePath }) {
    return {
        ...viteConfig,
        ...(basePath ? { base: basePath } : {}),
    };
}
//////////////////////////////////////////////////////
// Utilities
export function serializeConfig(config) {
    return JSON.stringify(serializeFunctions(config));
}
export function deserializeConfig(config) {
    return deserializeFunctions(JSON.parse(config));
}
export function serializeFunctions(value, key) {
    if (Array.isArray(value)) {
        return value.map((v) => serializeFunctions(v));
    }
    if (typeof value === 'object' && value !== null) {
        return Object.keys(value).reduce((acc, key) => {
            if (key[0] === '_')
                return acc;
            acc[key] = serializeFunctions(value[key], key);
            return acc;
        }, {});
    }
    if (typeof value === 'function') {
        let serialized = value.toString();
        if (key && (serialized.startsWith(key) || serialized.startsWith(`async ${key}`))) {
            serialized = serialized.replace(key, 'function');
        }
        return `_vocs-fn_${serialized}`;
    }
    return value;
}
export function deserializeFunctions(value) {
    if (Array.isArray(value)) {
        return value.map(deserializeFunctions);
    }
    if (typeof value === 'object' && value !== null) {
        return Object.keys(value).reduce((acc, key) => {
            acc[key] = deserializeFunctions(value[key]);
            return acc;
        }, {});
    }
    if (typeof value === 'string' && value.includes('_vocs-fn_')) {
        return new Function(`return ${value.slice(9)}`)();
    }
    return value;
}
export const deserializeFunctionsStringified = `
  function deserializeFunctions(value) {
    if (Array.isArray(value)) {
      return value.map(deserializeFunctions)
    } else if (typeof value === 'object' && value !== null) {
      return Object.keys(value).reduce((acc, key) => {
        acc[key] = deserializeFunctions(value[key])
        return acc
      }, {})
    } else if (typeof value === 'string' && value.includes('_vocs-fn_')) {
      return new Function(\`return \${value.slice(9)}\`)()
    } else {
      return value
    }
  }
`;
//# sourceMappingURL=config.js.map