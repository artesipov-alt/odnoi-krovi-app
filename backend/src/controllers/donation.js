const pool = require('../config/db');
const logger = require('winston');
const { sendNotification } = require('../services/notification.service');

async function planDonation(req, res) {
  const { donor_pet_id, clinic_id, date } = req.body;
  try {
    const { rows } = await pool.query('INSERT INTO odna_krov.donations (donor_pet_id, clinic_id, date, status) VALUES ($1, $2, $3, $4) RETURNING id', [donor_pet_id, clinic_id, date, 'planned']);
    logger.info(`Donation planned with ID ${rows[0].id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action, details) VALUES ((SELECT owner_id FROM odna_krov.pets WHERE id = $1), $2, $3)', [donor_pet_id, 'donation_plan', { donation_id: rows[0].id }]);
    sendNotification(clinic_id, `Новая запланированная донация от питомца ID ${donor_pet_id} на ${date}.`);
    const { rows: owner } = await pool.query('SELECT owner_id FROM odna_krov.pets WHERE id = $1', [donor_pet_id]);
    sendNotification(owner[0].owner_id, `Ваша донация для питомца ID ${donor_pet_id} запланирована на ${date}.`);
    res.status(201).json({ id: rows[0].id, success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

async function getDonations(req, res) {
  const { donor_pet_id, clinic_id } = req.query;
  let query = 'SELECT * FROM odna_krov.donations WHERE 1=1';
  const params = [];
  if (donor_pet_id) { query += ' AND donor_pet_id = $1'; params.push(donor_pet_id); }
  if (clinic_id) { query += ` AND clinic_id = $${params.length + 1}`; params.push(clinic_id); }
  try {
    const { rows } = await pool.query(query, params);
    logger.info(`Donations retrieved for pet ${donor_pet_id || 'all'} and clinic ${clinic_id || 'all'}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action) VALUES ($1, $2)', [req.user.id, 'donations_view']);
    res.json(rows);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

async function updateDonationStatus(req, res) {
  const { id } = req.params;
  const { status } = req.body;
  try {
    await pool.query('UPDATE odna_krov.donations SET status = $1 WHERE id = $2', [status, id]);
    logger.info(`Donation status updated to ${status} for ID ${id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action, details) VALUES ((SELECT clinic_id FROM odna_krov.donations WHERE id = $1), $2, $3)', [id, 'donation_update', { status }]);
    const { rows } = await pool.query('SELECT donor_pet_id, clinic_id FROM odna_krov.donations WHERE id = $1', [id]);
    sendNotification(rows[0].clinic_id, `Статус донации ID ${id} изменен на ${status}.`);
    const { rows: owner } = await pool.query('SELECT owner_id FROM odna_krov.pets WHERE id = $1', [rows[0].donor_pet_id]);
    sendNotification(owner[0].owner_id, `Статус вашей донации ID ${id} изменен на ${status}.`);
    res.json({ success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

async function deleteDonation(req, res) {
  const { id } = req.params;
  try {
    await pool.query('DELETE FROM odna_krov.donations WHERE id = $1', [id]);
    logger.info(`Donation deleted with ID ${id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action, details) VALUES ($1, $2, $3)', [req.user.id, 'donation_delete', { donation_id: id }]);
    res.json({ success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

module.exports = { planDonation, getDonations, updateDonationStatus, deleteDonation };