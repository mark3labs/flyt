/**
 * Merge classes.
 *
 * This function has two signatures, depending on whether a `node` was passed.
 *
 * @param node
 *   Optionally, node whose classes to append to (should be `Element`).
 * @param conditionals
 *   Class configuration to merge.
 * @returns
 *   The given node, if any, or a list of strings.
 */
export const classnames: (<T extends import("hast").Nodes>(node: T, ...conditionals: Array<Conditional>) => T) & ((...conditionals: Array<Conditional>) => Array<string>);
export type Nodes = import('hast').Nodes;
/**
 * Map of class names as keys, with whether they’re turned on or not as
 * values.
 */
export type ClassMap = Record<string, boolean>;
/**
 * Different ways to turn class names on or off.
 */
export type Conditional = Array<Array<ConditionalMap | ConditionalPrimitive> | ConditionalMap | ConditionalPrimitive> | ConditionalMap | ConditionalPrimitive | null | undefined;
/**
 * Map of class names as keys, with whether they’re turned on or not as
 * values.
 */
export type ConditionalMap = Record<string, boolean>;
/**
 * Basic class names.
 */
export type ConditionalPrimitive = number | string;
