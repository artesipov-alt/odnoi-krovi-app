const pool = require('../config/db');
const logger = require('winston');
const { sendNotification } = require('../services/notification.service');

async function register(req, res) {
  const { telegram_id, role, full_name, phone, email, consent_pd } = req.body;
  console.log('Register called with telegram_id:', telegram_id);
  if (!consent_pd) return res.status(400).json({ error: 'Consent required' });
  try {
    await pool.query('INSERT INTO odna_krov.users (telegram_id, role, full_name, phone, email, consent_pd) VALUES ($1, $2, $3, $4, $5, $6)', [telegram_id, role, full_name, phone, email, consent_pd]);
    logger.info(`User ${telegram_id} registered`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action, details) VALUES ((SELECT id FROM odna_krov.users WHERE telegram_id = $1), $2, $3)', [telegram_id, 'registration', { role }]);
    sendNotification(telegram_id, 'Добро пожаловать в "Одной Крови"!');
    res.status(201).json({ success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

async function getProfile(req, res) {
  const { telegram_id } = req.query;
  try {
    const { rows } = await pool.query('SELECT * FROM odna_krov.users WHERE telegram_id = $1', [telegram_id]);
    if (!rows.length) return res.status(404).json({ error: 'User not found' });
    logger.info(`Profile viewed for ${telegram_id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action) VALUES ($1, $2)', [rows[0].id, 'profile_view']);
    res.json(rows[0]);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

async function updateProfile(req, res) {
  const { telegram_id, full_name, phone, email } = req.body;
  try {
    await pool.query('UPDATE odna_krov.users SET full_name = $1, phone = $2, email = $3 WHERE telegram_id = $4', [full_name, phone, email, telegram_id]);
    logger.info(`Profile updated for ${telegram_id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action) VALUES ((SELECT id FROM odna_krov.users WHERE telegram_id = $1), $2)', [telegram_id, 'profile_update']);
    res.json({ success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

async function deleteProfile(req, res) {
  const { telegram_id } = req.query;
  try {
    await pool.query('DELETE FROM odna_krov.users WHERE telegram_id = $1', [telegram_id]);
    logger.info(`Profile deleted for ${telegram_id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action) VALUES ($1, $2)', [telegram_id, 'profile_delete']);
    res.json({ success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

module.exports = { register, getProfile, updateProfile, deleteProfile };
