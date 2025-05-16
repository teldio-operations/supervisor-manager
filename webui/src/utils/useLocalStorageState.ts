import { useCallback, useState } from "react";

export const useLocalStorageState = <T>(defaultValue: T, key: string) => {
  const [get, _set] = useState<T>(() => {
    const localData = localStorage.getItem(key);
    return localData ? JSON.parse(localData) : defaultValue;
  });

  const set = useCallback((value: T | React.SetStateAction<T>) => {
    localStorage.setItem(key, JSON.stringify(value));
    return _set(value);
  }, [key]);

  return [get, set] as const;
};