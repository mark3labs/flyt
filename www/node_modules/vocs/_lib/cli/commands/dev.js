import pc from 'picocolors';
import { createLogger } from 'vite';
import { version } from '../version.js';
export async function dev(_, { clean, clearScreen = true, host, port } = {}) {
    const { createDevServer } = await import('../../vite/devServer.js');
    const server = await createDevServer({ clean, host, port });
    await server.listen();
    const logger = createLogger();
    if (clearScreen)
        logger.clearScreen('info');
    logger.info('');
    logger.info(`  ${pc.green('[running]')} ${pc.bold('vocs')}@${pc.dim(`v${version}`)}`);
    logger.info('');
    server.printUrls();
}
//# sourceMappingURL=dev.js.map