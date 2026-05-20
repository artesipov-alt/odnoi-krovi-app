import axios from 'axios';

// Путь /api - является относительным, обращается к localhost

// const localURL = 'https://dev.1krovi.app/api';

export const instance = axios.create({
    timeout: 20000,
    baseURL: '/api',
    // baseURL: localURL,
    // headers: { 'Content-Type': 'application/json', 'X-Telegram-WebApp': 'true' },
});
