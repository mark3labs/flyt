import { existsSync, readFileSync } from 'node:fs';
import { resolve } from 'node:path';
import toml from 'toml';
import { loadConfigFromFile } from 'vite';
import { defineConfig, getDefaultConfig } from '../../config.js';
const moduleExtensions = ['js', 'jsx', 'ts', 'tsx', 'mjs', 'mts'];
const staticExtensions = ['toml', 'json'];
const extensions = [...moduleExtensions, ...staticExtensions];
const defaultConfigPaths = ['.vocs/config', 'vocs.config', 'Vocs'];
export async function resolveVocsConfig(parameters = {}) {
    const { command = 'serve', mode = 'development' } = parameters;
    const [configPath, ext] = (() => {
        for (const ext of extensions) {
            if (parameters.configPath)
                return parameters.configPath;
            for (const filePath of defaultConfigPaths)
                if (existsSync(resolve(process.cwd(), `${filePath}.${ext}`)))
                    return [`${filePath}.${ext}`, ext];
        }
        return [undefined, undefined];
    })();
    const result = await (async () => {
        if (!ext)
            return;
        if (moduleExtensions.includes(ext))
            return await loadConfigFromFile({ command, mode }, configPath);
        if (staticExtensions.includes(ext)) {
            const file = readFileSync(configPath, 'utf8');
            const rawConfig = (() => {
                if (ext === 'toml')
                    return camelCaseKeys(toml.parse(file));
                if (ext === 'json')
                    return JSON.parse(file);
                return;
            })();
            const config = await defineConfig(rawConfig);
            return config ? { config } : undefined;
        }
        return;
    })();
    const config = (result ? result.config : await getDefaultConfig());
    return {
        config,
        configPath,
    };
}
function camelCaseKeys(obj) {
    if (typeof obj !== 'object')
        return obj;
    if (Array.isArray(obj))
        return obj.map(camelCaseKeys);
    return Object.fromEntries(Object.entries(obj).map(([key, value]) => [
        key.replace(/[-_](.)/g, (_, c) => c.toUpperCase()),
        camelCaseKeys(value),
    ]));
}
//# sourceMappingURL=resolveVocsConfig.js.map