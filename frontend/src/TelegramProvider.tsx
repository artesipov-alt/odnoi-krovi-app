import { useTelegramAuth } from 'hooks/useTelegramAuth';
import { createContext, ReactNode, useContext } from 'react';

const TelegramContext = createContext<ReturnType<typeof useTelegramAuth>>({
    user: null,
    isRegistered: false,
});

export const TelegramProvider = ({ children }: { children: ReactNode }) => {
    const auth = useTelegramAuth();

    return <TelegramContext.Provider value={auth}>{children}</TelegramContext.Provider>;
};

export const useTelegram = () => useContext(TelegramContext);
