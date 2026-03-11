import { useEffect } from 'react';

/**
 * Хук для блокировки прокрутки body
 * @param isLocked - если true, прокрутка блокируется
 */
const useBodyScrollLock = (isLocked: boolean) => {
    useEffect(() => {
        if (isLocked) {
            document.body.style.overflow = 'hidden';
            document.body.style.height = '100%';
        } else {
            document.body.style.overflow = '';
            document.body.style.height = '';
        }

        // Очистка при размонтировании
        return () => {
            document.body.style.overflow = '';
            document.body.style.height = '';
        };
    }, [isLocked]);
};

export default useBodyScrollLock;
