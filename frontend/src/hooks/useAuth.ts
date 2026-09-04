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
        let signinData: null | { data: SigninResponse } = null;

        // Для локальной разработки используем signinExtServ
        if (window.location.hostname === 'localhost') {
            const arturID = '11111111';
            const ruslanID = '2222222';
            signinData = await signinExtServ({ providerId: arturID, providerName: 'service' }); // Или другой тестовый ID
        } else if (window.WebApp?.initData) {
            signinData = await signinMax({ appInitData: window.WebApp.initData });

            localStorage.setItem('environment', 'max');
        } else {
            try {
                const healthResponse = await fetch('https://bridge.1krovi.app/health', {
                    method: 'HEAD',
                    signal: AbortSignal.timeout(5000),
                });

                if (healthResponse.ok) {
                    const s = document.createElement('script');

                    s.src = '/tg-js/js/telegram-web-app.js';
                    document.head.appendChild(s);

                    await new Promise((resolve) => setTimeout(resolve, 500));
                }
            } catch (e) {
                console.warn('Bridge недоступен, telegram-web-app.js не загружен:', e);
            }

            if (window.Telegram?.WebApp?.initData) {
                signinData = await signinTg({ appInitData: window.Telegram.WebApp.initData });

                localStorage.setItem('environment', 'tg');
            }
        }

        if (!signinData) {
            throw new Error('WebApp SDK не найден');
        }

        setUserId(signinData.data.userId);
    }, []);

    useEffect(() => {
        initialize();
    }, [initialize]);

    return { userId, initialize };
};
