type SetValue<type> = (newVal: type | ((prevVal: type) => type)) => void;
export declare function useLocalStorage<type>(key: string, defaultValue: type | undefined): [type | undefined, SetValue<type>];
export {};
//# sourceMappingURL=useLocalStorage.d.js.map