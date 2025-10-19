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
    if (!telegram_id || !role || !full_name || !phone || !email || consent_pd === undefined) {
      logger.warn(`Registration failed: Missing required fields for telegram_id: ${telegram_id}`);
      return res.status(400).json({ error: 'All fields are required' });
    }

    logger.info(`Attempting registration for telegram_id: ${telegram_id}`);

    const userResult = await pool.query(
      'INSERT INTO odna_krov.users (telegram_id, role, full_name, phone, email, consent_pd) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (telegram_id) DO NOTHING RETURNING telegram_id, role, full_name, phone, email, consent_pd',
      [telegram_id, role, full_name, phone, email, consent_pd]
    );

    if (!userResult.rows.length) {
      logger.warn(`User with telegram_id ${telegram_id} already exists`);
      return res.status(409).json({ error: 'User already exists' });
    }

    const newUser = userResult.rows[0];
    logger.info(`User registered: ${JSON.stringify(newUser)}`);

    try {
      await pool.query(
        'INSERT INTO odna_krov.logs (user_id, action, timestamp, details) VALUES ($1, $2, $3, $4::jsonb)',
        [telegram_id, 'user_registered', new Date().toISOString(), {}]
      );
      logger.info(`Log created for user_id: ${telegram_id}`);
    } catch (logError) {
      logger.error(`Failed to insert into logs for telegram_id ${telegram_id}: ${logError.message}`);
    }

    res.status(201).json({ message: 'Registration successful', user: newUser });
  } catch (error) {
    logger.error(`Registration error for telegram_id ${req.body.telegram_id || 'unknown'}: ${error.message}`);
    res.status(500).json({ error: 'Internal server error' });
  }
};

const getProfile = async (req, res) => {
  try {
    const telegram_id = req.user?.telegram_id;
    if (!telegram_id) {
      logger.error('No telegram_id in request');
      return res.status(401).json({ error: 'Unauthorized: No telegram_id provided' });
    }

    logger.info(`Fetching profile for telegram_id: ${telegram_id}`);

    const { rows } = await pool.query(
      'SELECT telegram_id, role, full_name, phone, email, consent_pd FROM odna_krov.users WHERE telegram_id = $1',
      [telegram_id]
    );

    if (!rows.length) {
      logger.warn(`User not found for telegram_id: ${telegram_id}`);
      return res.status(404).json({ error: 'User not found' });
    }

    logger.info(`Profile found: ${JSON.stringify(rows[0])}`);

    res.json(rows[0]);
  } catch (error) {
    logger.error(`Get profile error for telegram_id ${req.user?.telegram_id || 'unknown'}: ${error.message}`);
    res.status(500).json({ error: 'Internal server error' });
  }
};

const updateProfile = async (req, res) => {
  try {
    const telegram_id = req.user?.telegram_id;
    if (!telegram_id) {
      logger.error('No telegram_id in request');
      return res.status(401).json({ error: 'Unauthorized: No telegram_id provided' });
    }

    const { role, full_name, phone, email, consent_pd } = req.body;
    if (!role || !full_name || !phone || !email || consent_pd === undefined) {
      logger.warn(`Update profile failed: Missing required fields for telegram_id: ${telegram_id}`);
      return res.status(400).json({ error: 'All fields are required' });
    }

    const { rows } = await pool.query(
      'UPDATE odna_krov.users SET role = $1, full_name = $2, phone = $3, email = $4, consent_pd = $5 WHERE telegram_id = $6 RETURNING telegram_id, role, full_name, phone, email, consent_pd',
      [role, full_name, phone, email, consent_pd, telegram_id]
    );

    if (!rows.length) {
      logger.warn(`User not found for telegram_id: ${telegram_id}`);
      return res.status(404).json({ error: 'User not found' });
    }

    res.json(rows[0]);
  } catch (error) {
    logger.error(`Update profile error for telegram_id ${req.user?.telegram_id || 'unknown'}: ${error.message}`);
    res.status(500).json({ error: 'Internal server error' });
  }
};

const deleteProfile = async (req, res) => {
  try {
    const telegram_id = req.user?.telegram_id;
    if (!telegram_id) {
      logger.error('No telegram_id in request');
      return res.status(401).json({ error: 'Unauthorized: No telegram_id provided' });
    }

    const { rows } = await pool.query(
      'DELETE FROM odna_krov.users WHERE telegram_id = $1 RETURNING telegram_id',
      [telegram_id]
    );

    if (!rows.length) {
      logger.warn(`User not found for telegram_id: ${telegram_id}`);
      return res.status(404).json({ error: 'User not found' });
    }

    res.json({ message: 'User deleted' });
  } catch (error) {
    logger.error(`Delete profile error for telegram_id ${req.user?.telegram_id || 'unknown'}: ${error.message}`);
    res.status(500).json({ error: 'Internal server error' });
  }
};

module.exports = { register, getProfile, updateProfile, deleteProfile };