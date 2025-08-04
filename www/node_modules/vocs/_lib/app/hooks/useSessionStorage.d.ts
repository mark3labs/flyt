type SetValue<type> = (newVal: type | ((prevVal: type) => type)) => void;
export declare function useSessionStorage<type>(key: string, defaultValue: type | undefined): [type | undefined, SetValue<type>];
export {};
//# sourceMappingURL=useSessionStorage.d.js.map