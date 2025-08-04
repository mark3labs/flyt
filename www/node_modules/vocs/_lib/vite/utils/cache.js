import { resolve } from 'node:path';
import { default as fs } from 'fs-extra';
export const search = create('search');
export const twoslash = create('twoslash');
export function create(key) {
    return ({ cacheDir = resolve(import.meta.dirname, '../.vocs/cache'), } = {}) => {
        const pathname = (k) => resolve(cacheDir, `${key}${k ? `.${k}` : ''}.json`);
        return {
            get(k) {
                let data = fs.readJSONSync(pathname(k), { throws: false });
                data = JSON.parse(data ?? '{}');
                return data.value;
            },
            set(k, value) {
                fs.ensureDirSync(cacheDir);
                fs.writeJSONSync(pathname(k), JSON.stringify({ value }));
            },
        };
    };
}
export function clear({ cacheDir = resolve(import.meta.dirname, '../.vocs/cache'), } = {}) {
    if (!fs.existsSync(cacheDir))
        return;
    fs.rmSync(cacheDir, { recursive: true });
}
//# sourceMappingURL=cache.js.map