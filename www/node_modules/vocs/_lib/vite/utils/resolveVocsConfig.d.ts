import { type ConfigEnv } from 'vite';
import { type ParsedConfig } from '../../config.js';
type ResolveVocsConfigParameters = {
    command?: ConfigEnv['command'];
    configPath?: string;
    mode?: ConfigEnv['mode'];
};
export declare function resolveVocsConfig(parameters?: ResolveVocsConfigParameters): Promise<{
    config: ParsedConfig;
    configPath: string | undefined;
}>;
export {};
//# sourceMappingURL=resolveVocsConfig.d.js.map