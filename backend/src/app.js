require('dotenv').config({ path: '../.env' });

const express = require('express');
const cors = require('cors');
const pool = require('./config/db');
const winston = require('winston');

const logger = winston.createLogger({
  transports: [
    new winston.transports.File({ filename: 'logs/app.log' })
  ]
});

const app = express();
app.use(cors()); // Разрешаем CORS для frontend
app.use(express.json());

// Корневой маршрут
app.get('/', (req, res) => {
  res.json({ status: 'Welcome to Odnoi Krovi API' });
});

// Health check маршрут
app.get('/health', (req, res) => {
  res.json({ status: 'Server is running' });
});

// Middleware для Telegram-авторизации (моковый режим)
app.use('/api', (req, res, next) => {
  const initData = req.headers['x-telegram-init-data'];
  logger.info(`API request: ${req.method} ${req.url}, Headers: ${JSON.stringify(req.headers)}, req.user: ${JSON.stringify(req.user)}`);
  if (initData === 'test_init_data') {
    req.user = {
      id: 314638947,
      telegram_id: 314638947,
      full_name: 'Test User'
    };
    logger.info(`Mock user set: ${JSON.stringify(req.user)}`);
  } else {
    logger.warn(`Invalid or missing X-Telegram-Init-Data: ${initData}`);
    return res.status(401).json({ error: 'Unauthorized: Invalid Telegram init data' });
  }
  next();
});

// API роуты
const apiRoutes = require('./routes/api');
app.use('/api', apiRoutes);

// Обработка 404
app.use((req, res) => {
  logger.warn(`404: ${req.method} ${req.url}`);
  res.status(404).json({ error: 'Route not found' });
});

// Обработка ошибок
app.use((err, req, res, next) => {
  logger.error(`Server error: ${err.message}`);
  res.status(500).json({ error: 'Internal Server Error' });
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
  logger.info(`Backend started on port ${PORT}`);
  console.log(`Backend on ${PORT}`);
});