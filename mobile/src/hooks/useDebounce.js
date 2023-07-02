//https://github.com/sondnpt00343/tiktok-ui/blob/main/src/hooks/useDebounce.js
import { useState, useEffect } from 'react';

/**
 * A hook for handle text input change. It will wait a time before next action
 * @param {any} value : react state value 
 * @param {number | undefined} delay : wait time before next action
 * @returns 
 */
const useDebounce = (value, delay = 500) => {
    const [debouncedValue, setDebouncedValue] = useState(value);

    useEffect(() => {
        const handler = setTimeout(() => setDebouncedValue(value), delay);

        return () => clearTimeout(handler);
    }, [value]);

    return debouncedValue;
}

export default useDebounce;