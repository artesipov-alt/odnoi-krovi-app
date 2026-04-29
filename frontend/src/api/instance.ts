import axios from 'axios';

const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3001/api';

export const instance = axios.create({
    timeout: 20000,
    baseURL,
    // headers: { 'Content-Type': 'application/json', 'X-Telegram-WebApp': 'true' },
});
