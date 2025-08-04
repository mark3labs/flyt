/**
 * @module
 * CSRF Protection Middleware for Hono.
 */
import type { Context } from '../../context';
import type { MiddlewareHandler } from '../../types';
type IsAllowedOriginHandler = (origin: string, context: Context) => boolean;
interface CSRFOptions {
    origin?: string | string[] | IsAllowedOriginHandler;
}
/**
 * CSRF Protection Middleware for Hono.
 *
 * @see {@link https://hono.dev/docs/middleware/builtin/csrf}
 *
 * @param {CSRFOptions} [options] - The options for the CSRF protection middleware.
 * @param {string|string[]|(origin: string, context: Context) => boolean} [options.origin] - Specify origins.
 * @returns {MiddlewareHandler} The middleware handler function.
 *
 * @example
 * ```ts
 * const app = new Hono()
 *
 * app.use(csrf())
 *
 * // Specifying origins with using `origin` option
 * // string
 * app.use(csrf({ origin: 'myapp.example.com' }))
 *
 * // string[]
 * app.use(
 *   csrf({
 *     origin: ['myapp.example.com', 'development.myapp.example.com'],
 *   })
 * )
 *
 * // Function
 * // It is strongly recommended that the protocol be verified to ensure a match to `$`.
 * // You should *never* do a forward match.
 * app.use(
 *   '*',
 *   csrf({
 *     origin: (origin) => /https:\/\/(\w+\.)?myapp\.example\.com$/.test(origin),
 *   })
 * )
 * ```
 */
export declare const csrf: (options?: CSRFOptions) => MiddlewareHandler;
export {};
