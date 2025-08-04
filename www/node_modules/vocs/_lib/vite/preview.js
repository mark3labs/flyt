import { resolve } from 'node:path';
import { serve } from '@hono/node-server';
import { Hono } from 'hono';
import { compress } from 'hono/compress';
import { resolveVocsConfig } from './utils/resolveVocsConfig.js';
import { serveStatic } from './utils/serveStatic.js';
export async function preview({ outDir = 'dist' } = {}) {
    const { config } = await resolveVocsConfig();
    const { basePath, rootDir } = config;
    const app = new Hono();
    app.use('*', compress());
    app.use('/*', serveStatic({
        root: resolve(rootDir, outDir),
        rewriteRequestPath(path) {
            return basePath ? path.replace(basePath, '') : path;
        },
    }));
    return new Promise((res) => {
        async function createServer(port = 4173) {
            process.on('uncaughtException', (err) => {
                if (err.code !== 'EADDRINUSE')
                    throw err;
                process.removeAllListeners();
                createServer(port + 1);
            });
            const server = serve({
                fetch: app.fetch,
                port,
            }).on('listening', () => {
                res(Object.assign(server, { port }));
            });
        }
        createServer();
    });
}
//# sourceMappingURL=preview.js.map