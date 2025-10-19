const pool = require('../config/db');
const logger = require('winston');
const { sendNotification } = require('../services/notification.service');

async function sendMessage(req, res) {
  const { chat_id, sender_id, message } = req.body;
  try {
    await pool.query('INSERT INTO odna_krov.chat_messages (chat_id, sender_id, message) VALUES ($1, $2, $3)', [chat_id, sender_id, message]);
    logger.info(`Message sent in chat ${chat_id} by ${sender_id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action, details) VALUES ($1, $2, $3)', [sender_id, 'message_send', { chat_id }]);
    const { rows } = await pool.query('SELECT recipient_user_id, donor_user_id FROM odna_krov.chats WHERE id = $1', [chat_id]);
    const recipient = rows[0].recipient_user_id === sender_id ? rows[0].donor_user_id : rows[0].recipient_user_id;
    sendNotification(recipient, `Новое сообщение в чате ${chat_id}: "${message.substring(0, 50)}..."`);
    res.status(201).json({ success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

async function getMessages(req, res) {
  const { chat_id } = req.query;
  try {
    const { rows } = await pool.query('SELECT * FROM odna_krov.chat_messages WHERE chat_id = $1 ORDER BY timestamp ASC', [chat_id]);
    logger.info(`Messages retrieved for chat ${chat_id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action, details) VALUES ($1, $2, $3)', [req.user.id, 'messages_view', { chat_id }]);
    res.json(rows);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

async function deleteMessage(req, res) {
  const { id } = req.params;
  try {
    await pool.query('DELETE FROM odna_krov.chat_messages WHERE id = $1', [id]);
    logger.info(`Message deleted with ID ${id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action, details) VALUES ($1, $2, $3)', [req.user.id, 'message_delete', { message_id: id }]);
    res.json({ success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

module.exports = { sendMessage, getMessages, deleteMessage };