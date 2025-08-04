import { fromHtmlIsomorphic } from 'hast-util-from-html-isomorphic';
import { toText } from 'hast-util-to-text';
import { createMermaidRenderer } from 'mermaid-isomorphic';
import svgToDataURI from 'mini-svg-data-uri';
import { parse } from 'space-separated-tokens';
import { visitParents } from 'unist-util-visit-parents';
/**
 * A regular expression to test for non-whitespace characters.
 */
const nonWhitespacePattern = /\w/;
const strategies = ['img-png', 'img-svg', 'inline-svg', 'pre-mermaid'];
/**
 * Validate the strategy option is valid.
 *
 * @param strategy
 *   The user provided strategy.
 * @returns
 *   The strategy if valid.
 */
function validateStrategy(strategy = 'inline-svg') {
    if (strategies.includes(strategy)) {
        return strategy;
    }
    throw new Error(`Expected strategy to be one of ${strategies.join(', ')}, got: ${strategy}`);
}
/**
 * Check if a hast element has the `language-mermaid` class name.
 *
 * @param element
 *   The hast element to check.
 * @param strategy
 *   The mermaid strategy to use.
 * @returns
 *   Whether or not the element has the `language-mermaid` class name.
 */
function isMermaidElement(element, strategy) {
    let mermaidClassName;
    if (element.tagName === 'pre') {
        if (strategy === 'pre-mermaid') {
            return false;
        }
        mermaidClassName = 'mermaid';
    }
    else if (element.tagName === 'code') {
        mermaidClassName = 'language-mermaid';
    }
    else {
        return false;
    }
    let className = element.properties?.className;
    if (typeof className === 'string') {
        className = parse(className);
    }
    if (!Array.isArray(className)) {
        return false;
    }
    return className.includes(mermaidClassName);
}
/**
 * Convert a render result to a data URI.
 *
 * @param result
 *   The render result to turn into a data URI.
 * @param isSrcset
 *   Whether the result is for a `srcset` or a `src` attribute.
 * @returns
 *   The data URI.
 */
function toDataURI(result, isSrcset) {
    if (result.screenshot) {
        return `data:image/png;base64,${result.screenshot.toString('base64')}`;
    }
    return isSrcset ? svgToDataURI.toSrcset(result.svg) : svgToDataURI(result.svg);
}
/**
 * Invert the given color scheme.
 *
 * @param colorScheme
 *   The color scheme to invert.
 * @returns
 *   `light` if the color scheme is `dark`, otherwise `dark`.
 */
function invertColorScheme(colorScheme) {
    return colorScheme === 'dark' ? 'light' : 'dark';
}
/**
 * Convert a Mermaid render result to a hast element.
 *
 * @param light
 *   The light Mermaid render result.
 * @param dark
 *   The dark mermaid render result.
 * @param colorScheme
 *   The default color scheme.
 * @returns
 *   If a dark render result exists, a responsive `<picture>` element that favors light mode.
 *   Otherwise an `<img>` element containing only the light mode result.
 */
function toImageElement(light, dark, colorScheme) {
    let imgResult;
    let pictureResult;
    if (colorScheme === 'dark') {
        imgResult = dark || light;
        pictureResult = light;
    }
    else {
        imgResult = light;
        pictureResult = dark;
    }
    const img = {
        type: 'element',
        tagName: 'img',
        properties: {
            alt: imgResult.description || '',
            height: imgResult.height,
            id: imgResult.id,
            src: toDataURI(imgResult),
            title: imgResult.title,
            width: imgResult.width
        },
        children: []
    };
    if (!dark) {
        return img;
    }
    return {
        type: 'element',
        tagName: 'picture',
        properties: {},
        children: [
            {
                type: 'element',
                tagName: 'source',
                properties: {
                    height: pictureResult.height,
                    id: pictureResult.id,
                    media: `(prefers-color-scheme: ${invertColorScheme(colorScheme)})`,
                    srcset: toDataURI(pictureResult, true),
                    width: pictureResult.width
                },
                children: []
            },
            img
        ]
    };
}
/**
 * Handle an error.
 *
 * If the error fallback is defined, use its result. Otherwise an error is thrown.
 *
 * @param reason
 *   The reason the error occurred.
 * @param instance
 *   The diagram code instance.
 * @param file
 *   The file on which the error should be reported.
 * @param options
 *   The render options.
 * @returns
 *   The error fallback renderer.
 */
function handleError(reason, instance, file, options) {
    const { ancestors, diagram } = instance;
    if (options?.errorFallback) {
        return options.errorFallback(ancestors.at(-1), diagram, reason, file);
    }
    const message = file.message(reason, {
        ruleId: 'rehype-mermaid',
        source: 'rehype-mermaid',
        ancestors
    });
    message.fatal = true;
    message.url = 'https://github.com/remcohaszing/rehype-mermaid';
    throw message;
}
/**
 * Get the color scheme from a `color-scheme` meta element.
 *
 * @param element
 *   The meta element to get the color scheme from.
 * @returns
 *   The detected color scheme.
 */
function getColorScheme(element) {
    if (typeof element.properties.content !== 'string') {
        return;
    }
    const colorSchemes = parse(element.properties.content);
    for (const colorScheme of colorSchemes) {
        if (colorScheme === 'light' || colorScheme === 'dark') {
            return colorScheme;
        }
    }
}
/**
 * A [rehype](https://rehype.js.org) plugin to render [mermaid](https://mermaid-js.github.io)
 * diagrams.
 *
 * @param options
 *   Options that may be used to tweak the output.
 */
const rehypeMermaid = (options) => {
    const strategy = validateStrategy(options?.strategy);
    const renderDiagrams = createMermaidRenderer(options);
    let colorScheme = options?.colorScheme;
    return (ast, file) => {
        const instances = [];
        visitParents(ast, 'element', (node, ancestors) => {
            if (!colorScheme && node.tagName === 'meta' && node.properties.name === 'color-scheme') {
                colorScheme = getColorScheme(node);
            }
            if (!isMermaidElement(node, strategy)) {
                return;
            }
            const parent = ancestors.at(-1);
            let inclusiveAncestors = ancestors;
            // This is <code> wrapped in a <pre> element.
            if (parent.type === 'element' && parent.tagName === 'pre') {
                for (const child of parent.children) {
                    // We allow whitespace text siblings, but any other siblings mean we don’t process the
                    // diagram.
                    if (child.type === 'text') {
                        if (nonWhitespacePattern.test(child.value)) {
                            return;
                        }
                    }
                    else if (child !== node) {
                        return;
                    }
                }
            }
            else {
                inclusiveAncestors = [...inclusiveAncestors, node];
            }
            instances.push({
                diagram: toText(node, { whitespace: 'pre' }),
                ancestors: inclusiveAncestors
            });
        });
        // Nothing to do. No need to start a browser in this case.
        if (!instances.length) {
            return;
        }
        if (strategy === 'pre-mermaid') {
            for (const { ancestors, diagram } of instances) {
                const parent = ancestors.at(-2);
                const node = ancestors.at(-1);
                parent.children[parent.children.indexOf(node)] = {
                    type: 'element',
                    tagName: 'pre',
                    properties: {
                        className: ['mermaid']
                    },
                    children: [{ type: 'text', value: diagram }]
                };
            }
            return;
        }
        const promises = [
            renderDiagrams(instances.map((instance) => instance.diagram), { ...options, screenshot: strategy === 'img-png' })
        ];
        if (options?.dark) {
            promises.push(renderDiagrams(instances.map((instance) => instance.diagram), {
                ...options,
                screenshot: strategy === 'img-png',
                mermaidConfig: options.dark === true ? { theme: 'dark' } : options.dark,
                prefix: `${options.prefix || 'mermaid'}-dark`
            }));
        }
        return Promise.all(promises).then(([lightResults, darkResults]) => {
            for (const [index, instance] of instances.entries()) {
                const lightResult = lightResults[index];
                const darkResult = darkResults?.[index];
                let replacement;
                if (lightResult.status === 'rejected') {
                    replacement = handleError(lightResult.reason, instance, file, options);
                    /* c8 ignore start */
                }
                else if (darkResult?.status === 'rejected') {
                    replacement = handleError(darkResult.reason, instance, file, options);
                    /* c8 ignore stop */
                }
                else if (strategy === 'inline-svg') {
                    replacement = fromHtmlIsomorphic(lightResult.value.svg, { fragment: true })
                        .children[0];
                }
                else {
                    replacement = toImageElement(lightResult.value, darkResult?.value, colorScheme);
                }
                const { ancestors } = instance;
                const node = ancestors.at(-1);
                const parent = ancestors.at(-2);
                const nodeIndex = parent.children.indexOf(node);
                if (replacement) {
                    parent.children[nodeIndex] = replacement;
                }
                else {
                    parent.children.splice(nodeIndex, 1);
                }
            }
        });
    };
};
export default rehypeMermaid;
//# sourceMappingURL=rehype-mermaid.js.map