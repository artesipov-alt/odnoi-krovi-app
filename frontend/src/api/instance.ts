import axios from 'axios';

export const instance = axios.create({
    timeout: 20000,
    baseURL: 'https://1krovi.app/api/v1/',
    // headers: { 'Content-Type': 'application/json', 'X-Telegram-WebApp': 'true' },
});
