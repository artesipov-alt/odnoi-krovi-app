const pool = require('../config/db');
const logger = require('winston');
const { S3Client, PutObjectCommand } = require('@aws-sdk/client-s3');
const TelegramBot = require('node-telegram-bot-api');

const bot = new TelegramBot(process.env.BOT_TOKEN);

async function createPet(req, res) {
  const { owner_id, name, type, photo, latitude, longitude, blood_group } = req.body;
  try {
    if (!owner_id || !name || !type) {
      logger.warn('Missing required fields: owner_id, name, type');
      return res.status(400).json({ error: 'Missing required fields: owner_id, name, type' });
    }
    if (parseInt(owner_id) !== req.user.id) {
      logger.warn(`Forbidden: User ${req.user.id} tried to create pet for owner ${owner_id}`);
      return res.status(403).json({ error: 'Forbidden: You can only create pets for yourself' });
    }

    let photo_url = null;
    if (photo) {
      const s3 = new S3Client({
        endpoint: process.env.S3_ENDPOINT,
        credentials: { accessKeyId: process.env.S3_ACCESS_KEY, secretAccessKey: process.env.S3_SECRET_KEY },
        region: 'ru-msk'
      });
      const photoKey = `pet_${Date.now()}.jpg`;
      const params = { Bucket: process.env.S3_BUCKET, Key: photoKey, Body: Buffer.from(photo, 'base64'), ACL: 'public-read' };
      await s3.send(new PutObjectCommand(params));
      photo_url = `${process.env.S3_ENDPOINT}/${process.env.S3_BUCKET}/${photoKey}`;
    }

    const { rows } = await pool.query(
      'INSERT INTO odna_krov.pets (owner_id, name, type, photo_url, latitude, longitude, blood_group) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id',
      [owner_id, name, type, photo_url, latitude, longitude, blood_group]
    );
    const pet_id = rows[0].id;

    logger.info(`Pet created for owner ${owner_id}, pet_id: ${pet_id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action, details) VALUES ($1, $2, $3)', 
      [owner_id, 'pet_create', JSON.stringify({ pet_id })]);
    await bot.sendMessage(req.user.id, `Питомец ${name} добавлен!`);
    res.status(201).json({ success: true, pet_id });
  } catch (err) {
    logger.error(`Error creating pet: ${err.message}`);
    res.status(500).json({ error: err.message });
  }
}

async function getPets(req, res) {
  const { owner_id } = req.query;
  try {
    if (owner_id && parseInt(owner_id) !== req.user.id) {
      return res.status(403).json({ error: 'Forbidden: You can only view your own pets' });
    }

    const query = owner_id ? 'SELECT * FROM odna_krov.pets WHERE owner_id = $1' : 'SELECT * FROM odna_krov.pets';
    const params = owner_id ? [owner_id] : [];
    const { rows } = await pool.query(query, params);

    logger.info(`Pets retrieved for owner ${owner_id || 'all'}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action, details) VALUES ($1, $2, $3)', 
      [req.user.id, 'pets_view', JSON.stringify({ owner_id: owner_id || 'all' })]);
    res.json(rows);
  } catch (err) {
    logger.error(`Error retrieving pets: ${err.message}`);
    res.status(500).json({ error: err.message });
  }
}

async function updatePet(req, res) {
  const { id } = req.params;
  const { name, type, photo, latitude, longitude, blood_group } = req.body;
  try {
    if (!name || !type) {
      return res.status(400).json({ error: 'Missing required fields: name, type' });
    }
    const { rows } = await pool.query('SELECT owner_id FROM odna_krov.pets WHERE id = $1', [id]);
    if (rows.length === 0) {
      return res.status(404).json({ error: 'Pet not found' });
    }
    if (rows[0].owner_id !== req.user.id) {
      return res.status(403).json({ error: 'Forbidden: You can only update your own pets' });
    }

    let photo_url = null;
    if (photo) {
      const s3 = new S3Client({
        endpoint: process.env.S3_ENDPOINT,
        credentials: { accessKeyId: process.env.S3_ACCESS_KEY, secretAccessKey: process.env.S3_SECRET_KEY },
        region: 'ru-msk'
      });
      const photoKey = `pet_${Date.now()}.jpg`;
      const params = { Bucket: process.env.S3_BUCKET, Key: photoKey, Body: Buffer.from(photo, 'base64'), ACL: 'public-read' };
      await s3.send(new PutObjectCommand(params));
      photo_url = `${process.env.S3_ENDPOINT}/${process.env.S3_BUCKET}/${photoKey}`;
    }

    await pool.query(
      'UPDATE odna_krov.pets SET name = $1, type = $2, photo_url = COALESCE($3, photo_url), latitude = $4, longitude = $5, blood_group = $6 WHERE id = $7',
      [name, type, photo_url, latitude, longitude, blood_group, id]
    );

    logger.info(`Pet updated with ID ${id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action, details) VALUES ($1, $2, $3)', 
      [req.user.id, 'pet_update', JSON.stringify({ pet_id: id })]);
    await bot.sendMessage(req.user.id, `Питомец с ID ${id} обновлён!`);
    res.json({ success: true, pet_id: id });
  } catch (err) {
    logger.error(`Error updating pet: ${err.message}`);
    res.status(500).json({ error: err.message });
  }
}

async function deletePet(req, res) {
  const { id } = req.params;
  try {
    const { rows } = await pool.query('SELECT owner_id FROM odna_krov.pets WHERE id = $1', [id]);
    if (rows.length === 0) {
      return res.status(404).json({ error: 'Pet not found' });
    }
    if (rows[0].owner_id !== req.user.id) {
      return res.status(403).json({ error: 'Forbidden: You can only delete your own pets' });
    }

    await pool.query('DELETE FROM odna_krov.pets WHERE id = $1', [id]);
    logger.info(`Pet deleted with ID ${id}`);
    await pool.query('INSERT INTO odna_krov.logs (user_id, action, details) VALUES ($1, $2, $3)', 
      [req.user.id, 'pet_delete', JSON.stringify({ pet_id: id })]);
    await bot.sendMessage(req.user.id, `Питомец с ID ${id} удалён!`);
    res.json({ success: true, pet_id: id });
  } catch (err) {
    logger.error(`Error deleting pet: ${err.message}`);
    res.status(500).json({ error: err.message });
  }
}

module.exports = { createPet, getPets, updatePet, deletePet };