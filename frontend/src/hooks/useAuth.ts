import { useCallback, useEffect, useState } from 'react';

import { signinExtServ } from 'api/apiServices/signinExtServ';
import { signinMax } from 'api/apiServices/signinMax';
import { signinTg } from 'api/apiServices/signinTg';
import { SigninResponse } from 'api/auth';

type UserAuth = {
    userId: string;
    initialize: () => Promise<void>;
};

export const useAuth = (): UserAuth => {
    const [userId, setUserId] = useState<string>('');

    const initialize = useCallback(async () => {
        let isWebAppNotFind = false;
        let signinData: null | { data: SigninResponse } = null;

        if (!window.Telegram?.WebApp?.initData) {
            isWebAppNotFind = true;
        } else {
            signinData = await signinTg({ appInitData: window.Telegram.WebApp.initData });
        }

        if (!window.WebApp?.initData) {
            isWebAppNotFind = true;
        } else {
            signinData = await signinMax({ appInitData: window.WebApp.initData });
        }

        signinData = await signinExtServ({ providerId: '248185030', providerName: 'telegram_bot' });
        // signinData = await signinExtServ({ providerId: '995757392', providerName: 'telegram_bot' });

        if (!signinData && isWebAppNotFind) {
            throw new Error('WebApp SDK не найден');
        }

        if (signinData) {
            setUserId(signinData.data.userId);
        }
    }, []);

    useEffect(() => {
        initialize();
    }, [initialize]);

    return { userId, initialize };
};
