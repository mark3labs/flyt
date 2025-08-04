import { jsx as _jsx } from "react/jsx-runtime";
import { config as virtualConfig } from 'virtual:config';
import { sha256 } from '@noble/hashes/sha256';
import { bytesToHex } from '@noble/hashes/utils';
import { createContext, useContext, useEffect, useState } from 'react';
import { deserializeConfig, serializeConfig } from '../../config.js';
const ConfigContext = createContext(virtualConfig);
export const configHash = import.meta.env.DEV
    ? bytesToHex(sha256(serializeConfig(virtualConfig))).slice(0, 8)
    : '';
export function getConfig() {
    if (typeof window !== 'undefined' && import.meta.env.DEV) {
        const storedConfig = window.localStorage.getItem(`vocs.config.${configHash}`);
        if (storedConfig)
            return deserializeConfig(storedConfig);
    }
    return virtualConfig;
}
export function ConfigProvider({ children, config: initialConfig, }) {
    const [config, setConfig] = useState(() => {
        if (initialConfig)
            return initialConfig;
        return getConfig();
    });
    useEffect(() => {
        if (import.meta.hot)
            import.meta.hot.on('vocs:config', setConfig);
    }, []);
    useEffect(() => {
        if (typeof window !== 'undefined' && import.meta.env.DEV)
            window.localStorage.setItem(`vocs.config.${configHash}`, serializeConfig(config));
    }, [config]);
    return _jsx(ConfigContext.Provider, { value: config, children: children });
}
export function useConfig() {
    return useContext(ConfigContext);
}
//# sourceMappingURL=useConfig.js.map