import type { RehypeShikiOptions } from '@shikijs/rehype';
import type { SearchOptions } from 'minisearch';
import type { ReactElement } from 'react';
import type { TwoslashOptions } from 'twoslash';
import type { PluggableList } from 'unified';
import type { UserConfig } from 'vite';
import type { borderRadiusVars, contentVars, fontFamilyVars, fontSizeVars, fontWeightVars, lineHeightVars, outlineVars, primitiveColorVars, semanticColorVars, sidebarVars, spaceVars, topNavVars, viewportVars, zIndexVars } from './app/styles/vars.css.js';
type RequiredBy<T, K extends keyof T> = Omit<T, K> & Required<Pick<T, K>>;
type RequiredProperties = 'blogDir' | 'markdown' | 'rootDir' | 'title' | 'titleTemplate';
export type Config<parsed extends boolean = false, colorScheme extends ColorScheme = ColorScheme> = RequiredBy<{
    /**
     * Whether or not to show the AI call-to-action dropdown on docs pages (ie. "Open in ChatGPT"),
     * as well as any configuration.
     */
    aiCta?: boolean | {
        /**
         * Query for the LLM.
         */
        query: (p: {
            location: string;
        }) => string;
    };
    /**
     * Configuration for the banner fixed to the top of the page.
     *
     * Can be a Markdown string, a React Element, or an object with the following properties:
     * - `dismissable`: Whether or not the banner can be dismissed.
     * - `backgroundColor`: The background color of the banner.
     * - `content`: The content of the banner.
     * - `height`: The height of the banner.
     * - `textColor`: The text color of the banner.
     */
    banner?: Banner<parsed>;
    /**
     * The base path the documentation will be deployed at. All assets and pages
     * will be prefixed with this path. This is useful for deploying to GitHub Pages.
     * For example, if you are deploying to `https://example.github.io/foo`, then the
     * basePath should be set to `/foo`.
     *
     * @example
     * /docs
     */
    basePath?: string;
    /**
     * The base URL for your documentation. This is used to populate the `<base>` tag in the
     * `<head>` of the page, and is used to form the `%logo` variable for dynamic OG images.
     *
     * @example
     * https://vocs.dev
     */
    baseUrl?: string;
    /**
     * Path to blog pages relative to project root.
     * Used to extract posts from the filesystem.
     *
     * @default "./pages/blog"
     */
    blogDir?: string;
    /**
     * Directory to store cache files.
     *
     * @default "node_modules/vocs/_lib/vite/.vocs/cache"
     */
    cacheDir?: string;
    /**
     * General description for the documentation.
     */
    description?: string;
    /**
     * Edit location for the documentation.
     */
    editLink?: Normalize<EditLink>;
    /**
     * Base font face.
     *
     * @default { google: "Inter" }
     */
    font?: Normalize<Font<parsed>>;
    /**
     * Additional tags to include in the `<head>` tag of the page HTML.
     */
    head?: ReactElement | {
        [path: string]: ReactElement;
    } | ((params: {
        path: string;
    }) => ReactElement | Promise<ReactElement>);
    /**
     * Icon URL.
     */
    iconUrl?: Normalize<IconUrl>;
    /**
     * Logo URL.
     */
    logoUrl?: Normalize<LogoUrl>;
    /**
     * OG Image URL. `null` to disable.
     *
     * Template variables: `%logo`, `%title`, `%description`
     *
     * @default "https://vocs.dev/api/og?logo=%logo&title=%title&description=%description"
     */
    ogImageUrl?: string | {
        [path: string]: string;
    };
    /**
     * Outline footer.
     */
    outlineFooter?: ReactElement;
    /**
     * Markdown configuration.
     */
    markdown?: Normalize<Markdown<parsed>>;
    /**
     * Documentation root directory. Can be an absolute path, or a path relative from
     * the location of the config file itself.
     *
     * @default "docs"
     */
    rootDir?: string;
    /**
     * Configuration for docs search.
     */
    search?: Normalize<Search>;
    /**
     * Navigation displayed on the sidebar.
     */
    sidebar?: Normalize<Sidebar>;
    /**
     * Social links displayed in the top navigation.
     */
    socials?: Normalize<Socials<parsed>>;
    /**
     * Set of sponsors to display on MDX directives and (optionally) the sidebar.
     */
    sponsors?: SponsorSet[];
    /**
     * Theme configuration.
     */
    theme?: Normalize<Theme<parsed, colorScheme>>;
    /**
     * Documentation title.
     *
     * @default "Docs"
     */
    title?: string;
    /**
     * Template for the page title.
     *
     * @default `%s – ${title}`
     */
    titleTemplate?: string;
    /**
     * Navigation displayed on the top.
     */
    topNav?: Normalize<TopNav<parsed>>;
    /**
     * TwoSlash configuration.
     */
    twoslash?: Normalize<TwoslashOptions>;
    /**
     * Vite configuration.
     */
    vite?: UserConfig;
}, parsed extends true ? RequiredProperties : never>;
export type ParsedConfig = Config<true>;
export declare function defineConfig<colorScheme extends ColorScheme = undefined>({ aiCta, blogDir, cacheDir, head, ogImageUrl, rootDir, title, titleTemplate, ...config }: Config<false, colorScheme>): Promise<ParsedConfig>;
export declare const getDefaultConfig: () => Promise<ParsedConfig>;
export declare function parseViteConfig(viteConfig: UserConfig | undefined, { basePath }: {
    basePath?: string;
}): UserConfig;
type Normalize<T> = {
    [K in keyof T]: T[K];
} & {};
export type Banner<parsed extends boolean = false> = Exclude<string | ReactElement | {
    /** Whether or not the banner can be dismissed. */
    dismissable?: boolean;
    /** The background color of the banner. */
    backgroundColor?: string;
    /** The content of the banner. */
    content: string | ReactElement;
    /** The height of the banner. */
    height?: string;
    /** The text color of the banner. */
    textColor?: string;
} | undefined, parsed extends true ? string | ReactElement : never>;
export type ColorScheme = 'light' | 'dark' | 'system' | undefined;
export type EditLink = {
    /**
     * Link pattern
     */
    pattern: string | (() => string);
    /**
     * Link text
     *
     * @default "Edit page"
     */
    text?: string;
};
type FontSource = Normalize<{
    /** Name of the Google Font to use. */
    google?: string;
}>;
type ParsedFont = {
    default?: FontSource;
    mono?: FontSource;
};
export type Font<parsed extends boolean = false> = parsed extends true ? ParsedFont : FontSource | ParsedFont;
export type ImageUrl = string | {
    light: string;
    dark: string;
};
export type IconUrl = ImageUrl;
export type LogoUrl = ImageUrl;
export type Markdown<parsed extends boolean = false> = RequiredBy<{
    code?: Normalize<RehypeShikiOptions>;
    remarkPlugins?: PluggableList;
    rehypePlugins?: PluggableList;
}, parsed extends true ? 'code' : never>;
export type Search = SearchOptions;
export type SidebarItem = {
    /** Whether or not to collapse the sidebar item by default. */
    collapsed?: boolean;
    /** Whether or not to disable the sidebar item. */
    disabled?: boolean;
    /** Text to display on the sidebar. */
    text: string;
    /** Optional pathname to the target documentation page. */
    link?: string;
    /** Optional children to nest under this item. */
    items?: SidebarItem[];
};
export type Sidebar = SidebarItem[] | {
    [path: string]: SidebarItem[] | {
        backLink?: boolean;
        items: SidebarItem[];
    };
};
export type SocialType = 'discord' | 'github' | 'telegram' | 'warpcast' | 'x';
export type SocialItem = {
    /** Social icon to display. */
    icon: SocialType;
    /** Label for the social. */
    label?: string;
    /** Link to the social. */
    link: string;
};
export type ParsedSocialItem = Required<SocialItem> & {
    /** The type of social item. */
    type: SocialType;
};
export type Socials<parsed extends boolean = false> = parsed extends true ? ParsedSocialItem[] : SocialItem[];
export type Sponsor = {
    /** The name of the sponsor. */
    name: string;
    /** The link to the sponsor's website. */
    link: string;
    /** The image to display for the sponsor. */
    image: string;
};
export type SponsorSet = {
    /** The list of sponsors to display. */
    items: (Sponsor | null)[][];
    /** The name of the sponsor set (e.g. "Gold Sponsors", "Collaborators", etc). */
    name: string;
    /** The height of the sponsor images. */
    height?: number;
};
export type ThemeVariables<variables extends Record<string, unknown>, value> = {
    [key in keyof variables]?: value;
};
export type Theme<parsed extends boolean = false, colorScheme extends ColorScheme = ColorScheme, colorValue = colorScheme extends 'system' | undefined ? {
    light: string;
    dark: string;
} : string> = {
    accentColor?: Exclude<string | (colorScheme extends 'system' | undefined ? {
        light: string;
        dark: string;
    } : never) | Required<ThemeVariables<Pick<typeof primitiveColorVars, 'backgroundAccent' | 'backgroundAccentHover' | 'backgroundAccentText' | 'borderAccent' | 'textAccent' | 'textAccentHover'>, colorValue>>, parsed extends true ? string | {
        light: string;
        dark: string;
    } : never>;
    colorScheme?: colorScheme;
    variables?: {
        borderRadius?: ThemeVariables<typeof borderRadiusVars, string>;
        color?: ThemeVariables<typeof primitiveColorVars, colorValue> & ThemeVariables<typeof semanticColorVars, colorValue>;
        content?: ThemeVariables<typeof contentVars, string>;
        fontFamily?: ThemeVariables<typeof fontFamilyVars, string>;
        fontSize?: ThemeVariables<typeof fontSizeVars, string>;
        fontWeight?: ThemeVariables<typeof fontWeightVars, string>;
        lineHeight?: ThemeVariables<typeof lineHeightVars, string>;
        outline?: ThemeVariables<typeof outlineVars, string>;
        sidebar?: ThemeVariables<typeof sidebarVars, string>;
        space?: ThemeVariables<typeof spaceVars, string>;
        topNav?: ThemeVariables<typeof topNavVars, string>;
        viewport?: ThemeVariables<typeof viewportVars, string>;
        zIndex?: ThemeVariables<typeof zIndexVars, string>;
    };
};
export type TopNavItem<parsed extends boolean = false> = {
    element: ReactElement;
    link?: never;
    items?: never;
    match?: never;
    text?: never;
} | ({
    element?: never;
    match?: string;
    text: string;
} & ({
    link: string;
    items?: never;
} | {
    link?: string;
    items: parsed extends true ? ParsedTopNavItem[] : TopNavItem[];
}));
export type ParsedTopNavItem = TopNavItem<true> & {
    id: number;
};
export type TopNav<parsed extends boolean = false> = parsed extends true ? ParsedTopNavItem[] : TopNavItem[];
export declare function serializeConfig(config: Config): string;
export declare function deserializeConfig(config: string): any;
export declare function serializeFunctions(value: any, key?: string): any;
export declare function deserializeFunctions(value: any): any;
export declare const deserializeFunctionsStringified = "\n  function deserializeFunctions(value) {\n    if (Array.isArray(value)) {\n      return value.map(deserializeFunctions)\n    } else if (typeof value === 'object' && value !== null) {\n      return Object.keys(value).reduce((acc, key) => {\n        acc[key] = deserializeFunctions(value[key])\n        return acc\n      }, {})\n    } else if (typeof value === 'string' && value.includes('_vocs-fn_')) {\n      return new Function(`return ${value.slice(9)}`)()\n    } else {\n      return value\n    }\n  }\n";
export {};
//# sourceMappingURL=config.d.js.map