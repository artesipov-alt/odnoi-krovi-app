const pool = require('../config/db');
const logger = require('winston');

async function createVetClinic(req, res) {
  const { user_id, name, phone, website, location, latitude, longitude } = req.body;
  try {
    await pool.query('INSERT INTO odna_krov.vet_clinics (user_id, name, phone, website, location, latitude, longitude) VALUES ($1, $2, $3, $4, $5, $6, $7)', [user_id, name, phone, website, location, latitude, longitude]);
    logger.info(`Vet clinic created for user ${user_id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action) VALUES ($1, $2)', [user_id, 'vet_clinic_create']);
    res.status(201).json({ success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

async function getVetClinics(req, res) {
  const { user_id } = req.query;
  try {
    const query = user_id ? 'SELECT * FROM odna_krov.vet_clinics WHERE user_id = $1' : 'SELECT * FROM odna_krov.vet_clinics';
    const params = user_id ? [user_id] : [];
    const { rows } = await pool.query(query, params);
    logger.info(`Vet clinics retrieved for user ${user_id || 'all'}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action) VALUES ($1, $2)', [req.user.id, 'vet_clinics_view']);
    res.json(rows);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

async function updateVetClinic(req, res) {
  const { id } = req.params;
  const { name, phone, website, location, latitude, longitude } = req.body;
  try {
    await pool.query('UPDATE odna_krov.vet_clinics SET name = $1, phone = $2, website = $3, location = $4, latitude = $5, longitude = $6 WHERE id = $7', [name, phone, website, location, latitude, longitude, id]);
    logger.info(`Vet clinic updated with ID ${id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action, details) VALUES ($1, $2, $3)', [req.user.id, 'vet_clinic_update', { clinic_id: id }]);
    res.json({ success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

async function deleteVetClinic(req, res) {
  const { id } = req.params;
  try {
    await pool.query('DELETE FROM odna_krov.vet_clinics WHERE id = $1', [id]);
    logger.info(`Vet clinic deleted with ID ${id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action, details) VALUES ($1, $2, $3)', [req.user.id, 'vet_clinic_delete', { clinic_id: id }]);
    res.json({ success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

module.exports = { createVetClinic, getVetClinics, updateVetClinic, deleteVetClinic };