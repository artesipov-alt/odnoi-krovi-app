import axios from 'axios';

// Путь /api - является отностельным, обращается к localhost

// const localURL = 'https://dev.1krovi.app/api';

export const instance = axios.create({
    timeout: 20000,
    baseURL: '/api',
    // headers: { 'Content-Type': 'application/json', 'X-Telegram-WebApp': 'true' },
});
