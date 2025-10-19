const pool = require('../config/db');
const logger = require('winston');
const { sendNotification } = require('../services/notification.service');

async function createChat(req, res) {
  const { recipient_user_id, donor_user_id } = req.body;
  try {
    const { rows } = await pool.query('INSERT INTO odna_krov.chats (recipient_user_id, donor_user_id) VALUES ($1, $2) RETURNING id', [recipient_user_id, donor_user_id]);
    logger.info(`Chat created with ID ${rows[0].id} between ${recipient_user_id} and ${donor_user_id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action, details) VALUES ($1, $2, $3)', [donor_user_id, 'chat_create', { chat_id: rows[0].id }]);
    sendNotification(recipient_user_id, `Новый чат от донора ${donor_user_id}.`);
    res.status(201).json({ id: rows[0].id, success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

async function getChats(req, res) {
  const { user_id } = req.query;
  try {
    const { rows } = await pool.query('SELECT * FROM odna_krov.chats WHERE recipient_user_id = $1 OR donor_user_id = $1', [user_id]);
    logger.info(`Chats retrieved for user ${user_id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action) VALUES ($1, $2)', [user_id, 'chats_view']);
    res.json(rows);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

async function deleteChat(req, res) {
  const { id } = req.params;
  try {
    await pool.query('DELETE FROM odna_krov.chats WHERE id = $1', [id]);
    logger.info(`Chat deleted with ID ${id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action, details) VALUES ($1, $2, $3)', [req.user.id, 'chat_delete', { chat_id: id }]);
    res.json({ success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

module.exports = { createChat, getChats, deleteChat };