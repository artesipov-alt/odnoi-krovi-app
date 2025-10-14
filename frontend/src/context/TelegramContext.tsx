import React, { createContext, useContext, ReactNode } from 'react';
import { useTelegramAuth } from '@/services/telegram.service';

const TelegramContext = createContext<ReturnType<typeof useTelegramAuth>>({
  user: null,
  initData: null,
  isRegistered: false,
  recheckRegistered: async () => {},
});

export const TelegramProvider = ({ children }: { children: ReactNode }) => {
  const auth = useTelegramAuth();
  return <TelegramContext.Provider value={auth}>{children}</TelegramContext.Provider>;
};

export const useTelegram = () => useContext(TelegramContext);