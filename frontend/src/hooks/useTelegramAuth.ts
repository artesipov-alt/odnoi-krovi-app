import { useCallback, useEffect, useState } from 'react';
import { toast } from 'react-toastify';

import { getUserByTelegramId } from 'api/apiServices/getUserByTelegramId';

import { TelegramUser } from '../types';

type TelegramAuth = {
    isRegistered: boolean;
    user: TelegramUser | null;
};

export const useTelegramAuth = (): TelegramAuth => {
    const [isRegistered, setIsRegistered] = useState(false);
    const [user, setUser] = useState<TelegramUser | null>(null);

    const initialize = useCallback(async () => {
        if (!window.Telegram?.WebApp?.initDataUnsafe?.user) {
            throw new Error('Telegram WebApp SDK не найден');
        }

        // const id = 995757392;
        // const id = 248185030;
        const { id } = window.Telegram.WebApp.initDataUnsafe.user;

        const { data, error } = await getUserByTelegramId(id);

        if (error) {
            toast.warn(error);

            return;
        }

        if (!data?.phone) {
            setUser({
                telegramId: id,
                id: data?.id!,
                fullName: data?.fullName!,
            });

            return;
        }

        setUser(data);
        setIsRegistered(true);
    }, []);

    useEffect(() => {
        initialize();
    }, [initialize]);

    return { user, isRegistered };
};
