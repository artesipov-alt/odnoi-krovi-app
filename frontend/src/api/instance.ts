import axios from 'axios';

export const instance = axios.create({
    timeout: 20000,
    // baseURL: 'https://1krovi.app/api',
    baseURL: 'http://localhost:3001/api',
    // headers: { 'Content-Type': 'application/json', 'X-Telegram-WebApp': 'true' },
});
