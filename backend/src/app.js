require('dotenv').config({ path: '../.env' });

const express = require('express');
const passport = require('passport');
const VKStrategy = require('passport-vk').Strategy;
const YandexStrategy = require('passport-yandex').Strategy;
const { validateInitData } = require('./services/auth.service');
const pool = require('./config/db');
const winston = require('winston');

const logger = winston.createLogger({
  transports: [
    new winston.transports.File({ filename: 'logs/app.log' })
  ]
});

const app = express();
app.use(express.json());
app.use(passport.initialize());

// Корневой маршрут
app.get('/', (req, res) => {
  res.json({ status: 'Welcome to Odnoi Krovi API' });
});

// Health check маршрут
app.get('/health', (req, res) => {
  res.json({ status: 'Server is running' });
});

// OAuth стратегии
passport.use(new VKStrategy({
  clientID: process.env.VK_CLIENT_ID,
  clientSecret: process.env.VK_CLIENT_SECRET,
  callbackURL: '/auth/vk/callback'
}, async (accessToken, refreshToken, params, profile, done) => {
  try {
    const { rows } = await pool.query('SELECT * FROM odna_krov.users WHERE email = $1', [profile.emails[0].value]);
    if (rows.length) return done(null, rows[0]);
    done(null, profile);
  } catch (err) {
    done(err);
  }
}));

passport.use(new YandexStrategy({
  clientID: process.env.YANDEX_CLIENT_ID,
  clientSecret: process.env.YANDEX_CLIENT_SECRET,
  callbackURL: '/auth/yandex/callback'
}, async (accessToken, refreshToken, profile, done) => {
  try {
    const { rows } = await pool.query('SELECT * FROM odna_krov.users WHERE email = $1', [profile.emails[0].value]);
    if (rows.length) return done(null, rows[0]);
    done(null, profile);
  } catch (err) {
    done(err);
  }
}));

// OAuth маршруты
app.get('/auth/vk', passport.authenticate('vk'));
app.get('/auth/vk/callback', passport.authenticate('vk'), (req, res) => res.redirect('https://1krovi.ru'));

app.get('/auth/yandex', passport.authenticate('yandex'));
app.get('/auth/yandex/callback', passport.authenticate('yandex'), (req, res) => res.redirect('https://1krovi.ru'));

// Middleware для Telegram auth (только для /api)
app.use('/api', (req, res, next) => {
  const initData = req.headers['x-telegram-init-data'];
  if (!initData || !validateInitData(initData, process.env.BOT_TOKEN)) {
    return res.status(401).json({ error: 'Unauthorized' });
  }
  req.user = JSON.parse(new URLSearchParams(initData).get('user'));
  next();
});

// API роуты
const apiRoutes = require('./routes/api');
app.use('/api', apiRoutes);

// Обработка 404
app.use((req, res) => {
  res.status(404).json({ error: 'Route not found' });
});

app.listen(3000, () => {
  logger.info('Backend started on port 3000');
  console.log('Backend on 3000');
});