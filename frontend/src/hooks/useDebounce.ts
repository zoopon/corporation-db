import { useState, useEffect } from 'react';

/**
 * デバウンスフック - 入力値の変更を遅延させる
 * @param value - デバウンスしたい値
 * @param delay - 遅延時間（ミリ秒）
 * @returns デバウンスされた値
 */
export const useDebounce = <T>(value: T, delay: number): T => {
  const [debouncedValue, setDebouncedValue] = useState<T>(value);

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => {
      clearTimeout(handler);
    };
  }, [value, delay]);

  return debouncedValue;
};
