const pool = require('../config/db');
const winston = require('winston');

const logger = winston.createLogger({
  transports: [
    new winston.transports.File({ filename: 'logs/app.log' })
  ]
});

const register = async (req, res) => {
  try {
    const { telegram_id, role, full_name, phone, email, consent_pd } = req.body;
    // Вставка в users
    const userResult = await pool.query(
      'INSERT INTO odna_krov.users (telegram_id, role, full_name, phone, email, consent_pd) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (telegram_id) DO NOTHING RETURNING *',
      [telegram_id, role, full_name, phone, email, consent_pd]
    );
    if (!userResult.rows.length) {
      logger.warn(`User with telegram_id ${telegram_id} already exists`);
      return res.status(409).json({ message: 'User already exists' });
    }
    // Вставка в logs
    try {
      await pool.query(
        'INSERT INTO odna_krov.logs (user_id, action, timestamp, details) VALUES ($1, $2, $3, $4::jsonb)',
        [telegram_id, 'user_registered', new Date().toISOString(), '{}']
      );
      logger.info(`Log created for user_id: ${telegram_id}`);
    } catch (logError) {
      logger.error(`Failed to insert into logs: ${logError.message}`);
      // Продолжаем, не прерывая регистрацию
    }
    res.status(201).json({ message: 'User registered', user: userResult.rows[0] });
  } catch (error) {
    logger.error(`Registration error: ${error.message}`);
    res.status(500).json({ message: 'Registration failed', error: error.message });
  }
};

const getProfile = async (req, res) => {
  try {
    const { telegram_id } = req.user; // Предполагается, что middleware добавляет req.user
    const { rows } = await pool.query('SELECT * FROM odna_krov.users WHERE telegram_id = $1', [telegram_id]);
    if (!rows.length) {
      return res.status(404).json({ message: 'User not found' });
    }
    res.json(rows[0]);
  } catch (error) {
    logger.error(`Get profile error: ${error.message}`);
    res.status(500).json({ message: 'Failed to get profile', error: error.message });
  }
};

const updateProfile = async (req, res) => {
  try {
    const { telegram_id } = req.user;
    const { role, full_name, phone, email, consent_pd } = req.body;
    const { rows } = await pool.query(
      'UPDATE odna_krov.users SET role = $1, full_name = $2, phone = $3, email = $4, consent_pd = $5 WHERE telegram_id = $6 RETURNING *',
      [role, full_name, phone, email, consent_pd, telegram_id]
    );
    if (!rows.length) {
      return res.status(404).json({ message: 'User not found' });
    }
    res.json(rows[0]);
  } catch (error) {
    logger.error(`Update profile error: ${error.message}`);
    res.status(500).json({ message: 'Failed to update profile', error: error.message });
  }
};

const deleteProfile = async (req, res) => {
  try {
    const { telegram_id } = req.user;
    const { rows } = await pool.query('DELETE FROM odna_krov.users WHERE telegram_id = $1 RETURNING *', [telegram_id]);
    if (!rows.length) {
      return res.status(404).json({ message: 'User not found' });
    }
    res.json({ message: 'User deleted' });
  } catch (error) {
    logger.error(`Delete profile error: ${error.message}`);
    res.status(500).json({ message: 'Failed to delete profile', error: error.message });
  }
};

module.exports = { register, getProfile, updateProfile, deleteProfile };