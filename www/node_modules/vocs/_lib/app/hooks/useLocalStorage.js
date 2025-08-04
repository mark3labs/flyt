import { useCallback, useEffect, useState } from 'react';
export function useLocalStorage(key, defaultValue) {
    const [value, setValue] = useState();
    useEffect(() => {
        const initialValue = getItem(key);
        if (typeof initialValue === 'undefined' || initialValue === null) {
            setValue(typeof defaultValue === 'function' ? defaultValue() : defaultValue);
        }
        else {
            setValue(initialValue);
        }
    }, [defaultValue, key]);
    const setter = useCallback((updater) => {
        setValue((old) => {
            let newVal;
            if (typeof updater === 'function')
                newVal = updater(old);
            else
                newVal = updater;
            try {
                localStorage.setItem(key, JSON.stringify(newVal));
            }
            catch { }
            return newVal;
        });
    }, [key]);
    return [value, setter];
}
function getItem(key) {
    try {
        const itemValue = localStorage.getItem(key);
        if (typeof itemValue === 'string') {
            return JSON.parse(itemValue);
        }
        return undefined;
    }
    catch {
        return undefined;
    }
}
//# sourceMappingURL=useLocalStorage.js.map