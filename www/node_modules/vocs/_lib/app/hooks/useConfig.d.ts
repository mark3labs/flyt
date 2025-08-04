import { type ReactNode } from 'react';
import { type ParsedConfig } from '../../config.js';
export declare const configHash: string;
export declare function getConfig(): ParsedConfig;
export declare function ConfigProvider({ children, config: initialConfig, }: {
    children: ReactNode;
    config?: ParsedConfig;
}): import("react/jsx-runtime").JSX.Element;
export declare function useConfig(): ParsedConfig;
//# sourceMappingURL=useConfig.d.js.map