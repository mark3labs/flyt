/**
 * @typedef {import('hast').Nodes} Nodes
 */

/**
 * @typedef {Record<string, boolean>} ClassMap
 *   Map of class names as keys, with whether they’re turned on or not as
 *   values.
 * @typedef {Array<Array<ConditionalMap | ConditionalPrimitive> | ConditionalMap | ConditionalPrimitive> | ConditionalMap | ConditionalPrimitive | null | undefined} Conditional
 *   Different ways to turn class names on or off.
 * @typedef {Record<string, boolean>} ConditionalMap
 *   Map of class names as keys, with whether they’re turned on or not as
 *   values.
 * @typedef {number | string} ConditionalPrimitive
 *   Basic class names.
 */

import {parse} from 'space-separated-tokens'

const own = {}.hasOwnProperty

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
export const classnames =
  // Note: JSDoc overloads don’t support optional templates yet.
  /**
   * @type {(
   *   (<T extends Nodes>(node: T, ...conditionals: Array<Conditional>) => T) &
   *   ((...conditionals: Array<Conditional>) => Array<string>)
   * )}
   */
  (
    /**
     * @param {Conditional | Nodes | null | undefined} [node]
     * @param {...Conditional} conditionals
     */
    function (node, ...conditionals) {
      let index = -1
      /** @type {Map<string, boolean>} */
      const map = new Map()
      /** @type {Array<string>} */
      const list = []

      if (isNode(node)) {
        if (node.type !== 'element') throw new Error('Expected element node')

        if (Array.isArray(node.properties.className)) {
          add(map, node.properties.className)
        }

        // We’ll add to `list` later.
        node.properties.className = list
      } else {
        conditionals.unshift(node)
      }

      while (++index < conditionals.length) {
        add(map, conditionals[index])
      }

      for (const [key, value] of map) {
        if (value) list.push(key)
      }

      return isNode(node) ? node : list
    }
  )

/**
 * @param {Map<string, boolean>} result
 * @param {Conditional} conditional
 */
function add(result, conditional) {
  let index = -1
  /** @type {string} */
  let key
  /** @type {Array<string>} */
  let list

  if (typeof conditional === 'number') {
    result.set(String(conditional), true)
  } else if (typeof conditional === 'string') {
    list = parse(conditional)

    while (++index < list.length) {
      result.set(list[index], true)
    }
  } else if (conditional && typeof conditional === 'object') {
    if (Array.isArray(conditional)) {
      while (++index < conditional.length) {
        add(result, conditional[index])
      }
    } else {
      for (key in conditional) {
        if (own.call(conditional, key)) {
          result.set(key, conditional[key])
        }
      }
    }
  }
}

/**
 * @param {Conditional | Nodes | null | undefined} value
 * @returns {value is Nodes}
 */
function isNode(value) {
  return Boolean(value && typeof value === 'object' && 'type' in value)
}
