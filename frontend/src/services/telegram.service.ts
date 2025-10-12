import { useEffect, useState } from 'react';
import axios from 'axios';
import WebApp from '@twa-dev/sdk';

interface WebAppUser {
  id: number;
  first_name: string;
  last_name?: string;
}

interface RegisterData {
  telegram_id: number;
  role: 'pet_owner' | 'donor';
  full_name: string;
  phone: string;
  email: string;
  consent_pd: boolean;
}

interface RegisterResponse {
  message: string;
  user?: {
    telegram_id: number;
    role: string;
    full_name: string;
    phone: string;
    email: string;
    consent_pd: boolean;
  };
}

export function useTelegramAuth(): { user: WebAppUser | undefined; initData: string; isRegistered: boolean } {
  const [user, setUser] = useState<WebAppUser | undefined>(undefined);
  const [initData, setInitData] = useState<string>('');
  const [isRegistered, setIsRegistered] = useState<boolean>(false);

  useEffect(() => {
    const initialize = async () => {
      if (window.Telegram?.WebApp) {
        const tg = window.Telegram.WebApp;
        tg.ready();
        setUser(tg.initDataUnsafe?.user as WebAppUser | undefined);
        setInitData(tg.initData);
        if (tg.initData) {
          try {
            await axios.get('/api/users/profile', {
              headers: { 'X-Telegram-Init-Data': tg.initData },
            });
            setIsRegistered(true);
          } catch {
            setIsRegistered(false);
          }
        }
      } else {
        console.warn('Telegram WebApp SDK не найден. Используется моковый режим.');
        const mockUser = { id: 314638947, first_name: 'Тест', last_name: 'Пользователь' };
        setUser(mockUser);
        setInitData('test_init_data');
        try {
          await axios.get('/api/users/profile', {
            headers: { 'X-Telegram-Init-Data': 'test_init_data' },
          });
          setIsRegistered(true);
        } catch {
          setIsRegistered(false);
        }
      }
    };
    initialize();
  }, []);

  return { user, initData, isRegistered };
}

export const registerUser = async (data: RegisterData): Promise<RegisterResponse> => {
  try {
    const response = await axios.post<RegisterResponse>('/api/users/register', data);
    if (!response.data.message) {
      throw new Error('Некорректный ответ от сервера');
    }
    return response.data;
  } catch (error: any) {
    throw new Error(error.response?.data?.message || 'Ошибка регистрации');
  }
};