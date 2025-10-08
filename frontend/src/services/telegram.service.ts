import WebApp from '@twa-dev/sdk';

interface WebAppUser {
  id: number;
  first_name: string;
  last_name?: string;
}

export function useTelegramAuth(): { user: WebAppUser | undefined; initData: string } {
  const initData = WebApp.initData;
  return { user: WebApp.initDataUnsafe.user as WebAppUser | undefined, initData };
}