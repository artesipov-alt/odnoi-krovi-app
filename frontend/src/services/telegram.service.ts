import { useState, useEffect } from 'react';
import axios from 'axios';

interface TelegramUser {
  id: number;
  first_name: string;
  last_name?: string;
  username?: string;
}

interface TelegramAuth {
  user: TelegramUser | null;
  initData: string | null;
  isRegistered: boolean;
  recheckRegistered: () => Promise<void>;
}

export const useTelegramAuth = (): TelegramAuth => {
  const [user, setUser] = useState<TelegramUser | null>(null);
  const [initData, setInitData] = useState<string | null>(null);
  const [isRegistered, setIsRegistered] = useState(false);

  const initialize = async () => {
    try {
      if (window.Telegram?.WebApp?.initDataUnsafe?.user) {
        // Реальный режим Telegram
        const tgUser = window.Telegram.WebApp.initDataUnsafe.user;
        setUser({
          id: tgUser.id,
          first_name: tgUser.first_name,
          last_name: tgUser.last_name,
          username: tgUser.username,
        });
        setInitData(window.Telegram.WebApp.initData);
      } else {
        // Моковый режим
        console.log('Telegram WebApp SDK не найден. Используется моковый режим.');
        setUser({
          id: 314638947,
          first_name: 'Тест',
          last_name: 'Пользователь',
          username: 'test_user',
        });
        setInitData('test_init_data');
      }

      await recheckRegistered();
    } catch (error) {
      console.error('Initialization error:', error);
    }
  };

  const recheckRegistered = async () => {
    try {
      const response = await axios.get('/api/users/profile', {
        headers: { 'X-Telegram-Init-Data': initData || 'test_init_data' },
      });
      setIsRegistered(!!response.data);
    } catch (error) {
      console.error('Error checking profile:', error);
      setIsRegistered(false);
    }
  };

  useEffect(() => {
    initialize();
  }, []);

  return { user, initData, isRegistered, recheckRegistered };
};

export const registerUser = async (data: {
  telegram_id: number;
  role: 'pet_owner' | 'donor' | 'clinic_admin';
  full_name: string;
  phone: string;
  email: string;
  consent_pd: boolean;
}, initData: string | null) => {
  const response = await axios.post('/api/users/register', data, {
    headers: { 'X-Telegram-Init-Data': initData || 'test_init_data' },
  });
  return response.data;
};

export const checkProfile = async (initData: string | null) => {
  try {
    const response = await axios.get('/api/users/profile', {
      headers: { 'X-Telegram-Init-Data': initData || 'test_init_data' },
    });
    return !!response.data;
  } catch (error) {
    console.error('Error checking profile:', error);
    return false;
  }
};