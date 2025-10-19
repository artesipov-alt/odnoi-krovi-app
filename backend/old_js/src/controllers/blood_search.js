const pool = require('../config/db');
const logger = require('winston');
const { sendNotification } = require('../services/notification.service');
const { calculateDistance } = require('../services/geo.service');

async function createSearch(req, res) {
  const { pet_id } = req.body;
  try {
    const { rows: petRows } = await pool.query('SELECT latitude, longitude, blood_group FROM odna_krov.pets WHERE id = $1', [pet_id]);
    if (!petRows.length) return res.status(404).json({ error: 'Pet not found' });
    const petLat = petRows[0].latitude;
    const petLon = petRows[0].longitude;
    const bloodGroup = petRows[0].blood_group;

    const { rows: clinics } = await pool.query('SELECT user_id, latitude, longitude FROM odna_krov.vet_clinics');
    const suitableClinics = [];
    for (const clinic of clinics) {
      const distance = await calculateDistance(petLat, petLon, clinic.latitude, clinic.longitude);
      if (distance < 50) {
        const { rows: stocks } = await pool.query('SELECT id FROM odna_krov.blood_stocks WHERE clinic_id = $1 AND blood_type = $2 AND status = $3', [clinic.user_id, bloodGroup, 'active']);
        if (stocks.length) {
          suitableClinics.push(clinic.user_id);
          sendNotification(clinic.user_id, `Новый поиск крови для питомца ID ${pet_id}, группа ${bloodGroup}. Расстояние: ${distance} км.`);
        }
      }
    }

    await pool.query('INSERT INTO odna_krov.blood_searches (pet_id, status) VALUES ($1, $2)', [pet_id, 'active']);
    logger.info(`Blood search created for pet ${pet_id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action, details) VALUES ((SELECT owner_id FROM odna_krov.pets WHERE id = $1), $2, $3)', [pet_id, 'search_create', { suitableClinics }]);
    sendNotification(req.user.id, `Поиск крови для питомца ID ${pet_id} создан. Найдено подходящих клиник: ${suitableClinics.length}.`);
    res.status(201).json({ success: true, suitableClinics });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

async function getSearches(req, res) {
  const { pet_id } = req.query;
  try {
    const query = pet_id ? 'SELECT * FROM odna_krov.blood_searches WHERE pet_id = $1' : 'SELECT * FROM odna_krov.blood_searches';
    const params = pet_id ? [pet_id] : [];
    const { rows } = await pool.query(query, params);
    logger.info(`Searches retrieved for pet ${pet_id || 'all'}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action) VALUES ($1, $2)', [req.user.id, 'searches_view']);
    res.json(rows);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

async function updateSearch(req, res) {
  const { id } = req.params;
  const { status } = req.body;
  try {
    await pool.query('UPDATE odna_krov.blood_searches SET status = $1 WHERE id = $2', [status, id]);
    logger.info(`Blood search updated with ID ${id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action, details) VALUES ((SELECT owner_id FROM odna_krov.pets WHERE id = (SELECT pet_id FROM odna_krov.blood_searches WHERE id = $1)), $2, $3)', [id, 'search_update', { search_id: id }]);
    res.json({ success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

async function deleteSearch(req, res) {
  const { id } = req.params;
  try {
    await pool.query('DELETE FROM odna_krov.blood_searches WHERE id = $1', [id]);
    logger.info(`Blood search deleted with ID ${id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action, details) VALUES ((SELECT owner_id FROM odna_krov.pets WHERE id = (SELECT pet_id FROM odna_krov.blood_searches WHERE id = $1)), $2, $3)', [id, 'search_delete', { search_id: id }]);
    res.json({ success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

module.exports = { createSearch, getSearches, updateSearch, deleteSearch };