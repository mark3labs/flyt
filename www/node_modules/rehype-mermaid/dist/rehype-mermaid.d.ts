import { type Element, type ElementContent, type Root } from 'hast';
import { type CreateMermaidRendererOptions, type RenderOptions } from 'mermaid-isomorphic';
import { type Plugin } from 'unified';
import { type VFile } from 'vfile';
/**
 * Allowed output strategies.
 */
type Strategy = 'img-png' | 'img-svg' | 'inline-svg' | 'pre-mermaid';
type ColorScheme = 'dark' | 'light';
export interface RehypeMermaidOptions extends CreateMermaidRendererOptions, Omit<RenderOptions, 'screenshot'> {
    /**
     * If specified, add responsive dark mode using a `<picture>` element.
     *
     * This option is only supported by the `img-png` and `img-svg` strategies.
     */
    dark?: RenderOptions['mermaidConfig'] | true;
    /**
     * The default color scheme.
     *
     * If not specified, `rehype-mermaid` will determine the color scheme based on the `color-scheme`
     * meta tag. If this doesn’t exist, the default color scheme is `light`.
     */
    colorScheme?: ColorScheme;
    /**
     * Create a fallback node if processing of a mermaid diagram fails.
     *
     * @param element
     *   The hast element that could not be rendered.
     * @param diagram
     *   The Mermaid diagram that could not be rendered.
     * @param error
     *   The error that was thrown.
     * @param file
     *   The file on which the error occurred.
     * @returns
     *   A fallback node to render instead of the invalid diagram. If nothing is returned, the code
     *   block is removed
     */
    errorFallback?: (element: Element, diagram: string, error: unknown, file: VFile) => ElementContent | null | undefined | void;
    /**
     * How to insert the rendered diagram into the document.
     *
     * - `'img-png'`: An `<img>` tag with the diagram as a base64 PNG data URL.
     * - `'img-svg'`: An `<img>` tag with the diagram as an SVG data URL.
     * - `'inline-svg'`: The SVG image as an inline `<svg>` element.
     * - `'pre-mermaid'`: The raw mermaid diagram as a child of a `<pre class="mermaid">` element.
     *
     * @default 'inline-svg'
     */
    strategy?: Strategy;
}
/**
 * A [rehype](https://rehype.js.org) plugin to render [mermaid](https://mermaid-js.github.io)
 * diagrams.
 *
 * @param options
 *   Options that may be used to tweak the output.
 */
declare const rehypeMermaid: Plugin<[RehypeMermaidOptions?], Root>;
export default rehypeMermaid;
//# sourceMappingURL=rehype-mermaid.d.ts.map