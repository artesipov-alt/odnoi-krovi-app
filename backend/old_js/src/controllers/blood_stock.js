const pool = require('../config/db');
const logger = require('winston');
const { sendNotification } = require('../services/notification.service');

async function createStock(req, res) {
  const { clinic_id, pet_type, blood_type, volume_ml, price_rub, expiration_date, status } = req.body;
  try {
    const { rows } = await pool.query('INSERT INTO odna_krov.blood_stocks (clinic_id, pet_type, blood_type, volume_ml, price_rub, expiration_date, status) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id', [clinic_id, pet_type, blood_type, volume_ml, price_rub, expiration_date, status || 'active']);
    logger.info(`Blood stock created with ID ${rows[0].id} for clinic ${clinic_id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action, details) VALUES ($1, $2, $3)', [clinic_id, 'stock_create', { stock_id: rows[0].id }]);
    sendNotification(clinic_id, `Новый запас крови создан: ${blood_type}, объем ${volume_ml} мл.`);
    res.status(201).json({ id: rows[0].id, success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

async function getStocks(req, res) {
  const { clinic_id, pet_type } = req.query;
  let query = 'SELECT * FROM odna_krov.blood_stocks WHERE 1=1';
  const params = [];
  if (clinic_id) { query += ' AND clinic_id = $1'; params.push(clinic_id); }
  if (pet_type) { query += ` AND pet_type = $${params.length + 1}`; params.push(pet_type); }
  try {
    const { rows } = await pool.query(query, params);
    logger.info(`Stocks retrieved for clinic ${clinic_id || 'all'}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action) VALUES ($1, $2)', [req.user.id, 'stocks_view']);
    res.json(rows);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

async function updateStock(req, res) {
  const { id } = req.params;
  const { volume_ml, price_rub, expiration_date, status } = req.body;
  try {
    await pool.query('UPDATE odna_krov.blood_stocks SET volume_ml = $1, price_rub = $2, expiration_date = $3, status = $4 WHERE id = $5', [volume_ml, price_rub, expiration_date, status, id]);
    logger.info(`Blood stock updated with ID ${id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action, details) VALUES ((SELECT clinic_id FROM odna_krov.blood_stocks WHERE id = $1), $2, $3)', [id, 'stock_update', { stock_id: id }]);
    res.json({ success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

async function deleteStock(req, res) {
  const { id } = req.params;
  try {
    await pool.query('DELETE FROM odna_krov.blood_stocks WHERE id = $1', [id]);
    logger.info(`Blood stock deleted with ID ${id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action, details) VALUES ((SELECT clinic_id FROM odna_krov.blood_stocks WHERE id = $1), $2, $3)', [id, 'stock_delete', { stock_id: id }]);
    res.json({ success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

async function bookStock(req, res) {
  const { id, booker_id } = req.body;
  try {
    await pool.query('UPDATE odna_krov.blood_stocks SET status = $1 WHERE id = $2', ['booked', id]);
    logger.info(`Blood stock booked with ID ${id} by ${booker_id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action, details) VALUES ((SELECT clinic_id FROM odna_krov.blood_stocks WHERE id = $1), $2, $3)', [id, 'stock_book', { booker_id }]);
    const { rows } = await pool.query('SELECT clinic_id FROM odna_krov.blood_stocks WHERE id = $1', [id]);
    sendNotification(rows[0].clinic_id, `Запас крови ID ${id} забронирован пользователем ${booker_id}.`);
    sendNotification(booker_id, `Вы забронировали запас крови ID ${id}.`);
    res.json({ success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

module.exports = { createStock, getStocks, updateStock, deleteStock, bookStock };