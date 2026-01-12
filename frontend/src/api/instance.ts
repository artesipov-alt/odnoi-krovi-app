import axios from 'axios';

export const instance = axios.create({
  timeout: 20000,
  baseURL: 'http://1krovi.app/api',
  // headers: { 'Content-Type': 'application/json', 'X-Telegram-WebApp': 'true' },
});
