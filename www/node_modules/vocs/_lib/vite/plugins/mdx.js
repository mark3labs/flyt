import mdxPlugin from '@mdx-js/rollup';
import rehypeShiki, {} from '@shikijs/rehype';
import { transformerNotationDiff, transformerNotationFocus, transformerNotationHighlight, transformerNotationWordHighlight, } from '@shikijs/transformers';
import { defaultTwoslashOptions as defaultTwoslashOptions_, transformerTwoslash, } from '@shikijs/twoslash';
import { h } from 'hastscript';
import rehypeAutolinkHeadings from 'rehype-autolink-headings';
import rehypeMermaid from 'rehype-mermaid';
import rehypeSlug from 'rehype-slug';
import remarkDirective from 'remark-directive';
import remarkFrontmatter from 'remark-frontmatter';
import remarkGfm from 'remark-gfm';
import remarkMdxFrontmatter from 'remark-mdx-frontmatter';
import { resolveVocsConfig } from '../utils/resolveVocsConfig.js';
import { rehypeShikiDisplayNotation } from './rehype/display-shiki-notation.js';
import { rehypeInlineShiki } from './rehype/inline-shiki.js';
import { remarkAuthors } from './remark/authors.js';
import { remarkBlogPosts } from './remark/blog-posts.js';
import { remarkCallout } from './remark/callout.js';
import { remarkCodeGroup } from './remark/code-group.js';
import { remarkCode } from './remark/code.js';
import { remarkDetails } from './remark/details.js';
import { remarkFilename } from './remark/filename.js';
import { remarkInferFrontmatter } from './remark/inferred-frontmatter.js';
import { remarkLinks } from './remark/links.js';
import { remarkSponsors } from './remark/sponsors.js';
import { remarkSteps } from './remark/steps.js';
import { remarkStrongBlock } from './remark/strong-block.js';
import { remarkSubheading } from './remark/subheading.js';
import { remarkTwoslash } from './remark/twoslash.js';
import { transformerEmptyLine } from './shiki/transformerEmptyLine.js';
import { transformerLineNumbers } from './shiki/transformerLineNumbers.js';
import { transformerNotationInclude } from './shiki/transformerNotationInclude.js';
import { transformerSplitIdentifiers } from './shiki/transformerSplitIdentifiers.js';
import { transformerTagLine } from './shiki/transformerTagLine.js';
import { transformerTitle } from './shiki/transformerTitle.js';
import { twoslashRenderer } from './shiki/twoslashRenderer.js';
import { twoslasher } from './shiki/twoslasher.js';
const defaultTwoslashOptions = defaultTwoslashOptions_();
export const getRemarkPlugins = ({ markdown } = {}) => [
    remarkDirective,
    remarkInferFrontmatter,
    remarkFrontmatter,
    remarkMdxFrontmatter,
    remarkGfm,
    remarkLinks,
    remarkBlogPosts,
    remarkCallout,
    remarkCode,
    remarkCodeGroup,
    remarkDetails,
    remarkFilename,
    remarkSponsors,
    remarkSteps,
    remarkStrongBlock,
    remarkSubheading,
    remarkTwoslash,
    remarkAuthors,
    ...(markdown?.remarkPlugins || []),
];
const defaultThemes = {
    dark: 'github-dark-dimmed',
    light: 'github-light',
};
export const getRehypePlugins = ({ markdown, cacheDir, rootDir = '', twoslash = {}, } = {}) => [
    rehypeSlug,
    [rehypeMermaid, { dark: true, strategy: 'img-svg' }],
    [
        rehypeShiki,
        {
            transformers: [
                transformerLineNumbers(),
                transformerNotationDiff(),
                transformerNotationFocus(),
                transformerNotationHighlight(),
                transformerNotationWordHighlight(),
                transformerNotationInclude({ rootDir }),
                transformerEmptyLine(),
                transformerTagLine(),
                transformerTitle(),
                twoslash !== false
                    ? transformerTwoslash({
                        explicitTrigger: true,
                        renderer: twoslashRenderer(),
                        twoslasher: twoslasher({ cacheDir }),
                        twoslashOptions: {
                            ...twoslash,
                            customTags: [
                                'allowErrors',
                                ...(defaultTwoslashOptions.customTags ?? []),
                                ...(twoslash.customTags ?? []),
                            ],
                            compilerOptions: {
                                ...(twoslash.compilerOptions ?? {}),
                                ...defaultTwoslashOptions.compilerOptions,
                            },
                        },
                    })
                    : null,
                transformerSplitIdentifiers(),
            ].filter(Boolean),
            themes: defaultThemes,
            ...markdown?.code,
        },
    ],
    [
        rehypeInlineShiki,
        {
            themes: defaultThemes,
            ...markdown?.code,
        },
    ],
    rehypeShikiDisplayNotation,
    [
        rehypeAutolinkHeadings,
        {
            behavior: 'append',
            content() {
                return [h('div', { dataAutolinkIcon: true })];
            },
        },
    ],
    ...(markdown?.rehypePlugins || []),
];
export async function mdx() {
    const { config } = await resolveVocsConfig();
    const { cacheDir, markdown, rootDir, twoslash } = config;
    const remarkPlugins = getRemarkPlugins({ markdown });
    const rehypePlugins = getRehypePlugins({ cacheDir, markdown, rootDir, twoslash });
    return [
        mdxPlugin({
            providerImportSource: 'vocs/mdx-react',
            remarkPlugins,
            rehypePlugins,
        }),
    ];
}
//# sourceMappingURL=mdx.js.map