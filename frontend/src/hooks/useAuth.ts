import { useCallback, useEffect, useState } from 'react';
import { toast } from 'react-toastify';

import { getUserByTelegramId } from 'api/apiServices/getUserByTelegramId';

import { TelegramUser } from '../types';
import { signinTg } from '../api/apiServices/signinTg';
import { useGetUserById } from './useGetUserById';
import { SigninResponse } from '../api/auth';
import { signinMax } from '../api/apiServices/signinMax';
import { signinExtServ } from '../api/apiServices/signinExtServ';

type TelegramAuth = {
    isRegistered: boolean;
    user: TelegramUser | null;
};

type UserAuth = {
    userId: string;
};

export const useAuth = (): UserAuth => {
    // const [isRegistered, setIsRegistered] = useState(false);
    // const [user, setUser] = useState<TelegramUser | null>(null);
    const [userId, setUserId] = useState<string>('');

    const initialize = useCallback(async () => {
        let isWebAppNotFind = false;
        let signinData: null | { data: SigninResponse } = null;

        if (!window.Telegram?.WebApp?.initDataUnsafe?.user) {
            isWebAppNotFind = true;
        } else {
            signinData = await signinTg({ appInitData: window.Telegram.WebApp.initData });
        }

        if (!window.WebApp.initDate) {
            isWebAppNotFind = true;
        } else {
            signinData = await signinMax({ appInitData: window.WebApp.initDate });
        }

        // signinData = await signinExtServ({ providerId: '248185030' });

        if (!signinData && isWebAppNotFind) {
            throw new Error('WebApp SDK не найден');
        }

        if (signinData) {
            setUserId(signinData.data.userId);
        }

        // const id = 995757392;
        // const id = 248185030;
        // const { id } = window.Telegram.WebApp.initDataUnsafe.user;

        // const { data, error } = await getUserByTelegramId(id);
        //
        // if (error) {
        //     toast.warn(error);
        //
        //     return;
        // }

        // if (!data?.phone) {
        //     setUser({
        //         telegramId: id,
        //         id: data?.id!,
        //         fullName: data?.fullName!,
        //     });
        //
        //     return;
        // }
        //
        // setUser(data);
        // setIsRegistered(true);
    }, []);

    useEffect(() => {
        initialize();
    }, [initialize]);

    return { userId };
};
