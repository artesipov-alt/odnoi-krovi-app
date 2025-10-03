const pool = require('../config/db');
const logger = require('winston');
const { S3Client, PutObjectCommand } = require('@aws-sdk/client-s3');

async function createPet(req, res) {
  const { owner_id, name, type, photo, latitude, longitude, blood_group } = req.body;
  try {
    let photo_url = null;
    if (photo) {
      const s3 = new S3Client({
        endpoint: process.env.S3_ENDPOINT,
        credentials: { accessKeyId: process.env.S3_ACCESS_KEY, secretAccessKey: process.env.S3_SECRET_KEY },
        region: 'ru-msk'
      });
      const photoKey = `pet_${Date.now()}.jpg`;
      const params = { Bucket: process.env.S3_BUCKET, Key: photoKey, Body: Buffer.from(photo, 'base64') };
      await s3.send(new PutObjectCommand(params));
      photo_url = `${process.env.S3_ENDPOINT}/${process.env.S3_BUCKET}/${photoKey}`;
    }

    await pool.query('INSERT INTO odna_krov.pets (owner_id, name, type, photo_url, latitude, longitude, blood_group) VALUES ($1, $2, $3, $4, $5, $6, $7)', [owner_id, name, type, photo_url, latitude, longitude, blood_group]);
    logger.info(`Pet created for owner ${owner_id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action) VALUES ($1, $2)', [owner_id, 'pet_create']);
    res.status(201).json({ success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

async function getPets(req, res) {
  const { owner_id } = req.query;
  try {
    const query = owner_id ? 'SELECT * FROM odna_krov.pets WHERE owner_id = $1' : 'SELECT * FROM odna_krov.pets';
    const params = owner_id ? [owner_id] : [];
    const { rows } = await pool.query(query, params);
    logger.info(`Pets retrieved for owner ${owner_id || 'all'}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action) VALUES ($1, $2)', [req.user.id, 'pets_view']);
    res.json(rows);
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

async function updatePet(req, res) {
  const { id } = req.params;
  const { name, type, photo, latitude, longitude, blood_group } = req.body;
  try {
    let photo_url = null;
    if (photo) {
      const s3 = new S3Client({
        endpoint: process.env.S3_ENDPOINT,
        credentials: { accessKeyId: process.env.S3_ACCESS_KEY, secretAccessKey: process.env.S3_SECRET_KEY },
        region: 'ru-msk'
      });
      const photoKey = `pet_${Date.now()}.jpg`;
      const params = { Bucket: process.env.S3_BUCKET, Key: photoKey, Body: Buffer.from(photo, 'base64') };
      await s3.send(new PutObjectCommand(params));
      photo_url = `${process.env.S3_ENDPOINT}/${process.env.S3_BUCKET}/${photoKey}`;
    }
    await pool.query('UPDATE odna_krov.pets SET name = $1, type = $2, photo_url = COALESCE($3, photo_url), latitude = $4, longitude = $5, blood_group = $6 WHERE id = $7', [name, type, photo_url, latitude, longitude, blood_group, id]);
    logger.info(`Pet updated with ID ${id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action, details) VALUES ($1, $2, $3)', [req.user.id, 'pet_update', { pet_id: id }]);
    res.json({ success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

async function deletePet(req, res) {
  const { id } = req.params;
  try {
    await pool.query('DELETE FROM odna_krov.pets WHERE id = $1', [id]);
    logger.info(`Pet deleted with ID ${id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action, details) VALUES ($1, $2, $3)', [req.user.id, 'pet_delete', { pet_id: id }]);
    res.json({ success: true });
  } catch (err) {
    res.status(500).json({ error: err.message });
  }
}

module.exports = { createPet, getPets, updatePet, deletePet };