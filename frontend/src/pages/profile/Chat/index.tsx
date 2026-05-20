import { FC, useEffect, useRef, useState } from 'react';

import Loading from 'components/Loading';

import styles from './Chat.module.less';

type Props = {
    onClose: () => void;
};

const Chat: FC<Props> = ({ onClose }) => {
    const [isLoading, setIsLoading] = useState(true);

    const hostRef = useRef<HTMLDivElement | null>(null);

    useEffect(() => {
        const scriptId = 'twin-widget-script';

        const initWidget = async () => {
            if (!window.appChatClient || !hostRef.current) return;

            hostRef.current.innerHTML = '';

            const api = await window.appChatClient(
                {
                    chatId: 'd0cbe3ec-06ab-495c-98f4-f686858b8f39',
                },
                {
                    host: hostRef.current,
                    injectStyles: `
      [data-id=chat-host] {
      // visibility: hidden;
        // bottom: 8px;
        // right: 8px;
        // align-items: end;
      }`,
                },
            );

            api.api.open();

            setTimeout(() => {
                const closeButton =
                    hostRef.current?.shadowRoot?.querySelector<HTMLElement>('[data-test-id="ButtonClose"]');

                closeButton?.addEventListener('click', () => {
                    onClose();
                });
            }, 0);

            setTimeout(() => {
                setIsLoading(false);
            }, 1000);
        };

        const existing = document.getElementById(scriptId) as HTMLScriptElement | null;

        if (existing) {
            initWidget();

            return;
        }

        const script = document.createElement('script');
        script.id = scriptId;
        script.src = 'https://twin24.ai/app/chat-client/widget.js';
        script.charset = 'utf-8';
        script.onload = initWidget;
        document.body.appendChild(script);

        return () => {
            const scriptToRemove = document.getElementById(scriptId);
            scriptToRemove?.remove();
        };
    }, [onClose]);

    return (
        <>
            <div ref={hostRef} />
            {isLoading && (
                <div className={styles.loading}>
                    <Loading size={90} thickness={4} />
                </div>
            )}
        </>
    );
};

export default Chat;
