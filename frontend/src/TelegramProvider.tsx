// import { useAuth } from 'hooks/useAuth';
// import { createContext, ReactNode, useContext } from 'react';
//
// const TelegramContext = createContext<ReturnType<typeof useAuth>>({
//     user: null,
//     isRegistered: false,
// });
//
// export const TelegramProvider = ({ children }: { children: ReactNode }) => {
//     const auth = useAuth();
//
//     return <TelegramContext.Provider value={auth}>{children}</TelegramContext.Provider>;
// };
//
// export const useTelegram = () => useContext(TelegramContext);
