const { Pool } = require('pg');

const pool = new Pool({
  host: process.env.DB_HOST,
  port: process.env.DB_PORT,
  user: process.env.DB_USER,
  password: process.env.DB_PASS,
  database: process.env.DB_NAME,
  ssl: { rejectUnauthorized: false } // Для IP-хоста без SSL
});

module.exports = pool;

pool.query('SELECT 1', (err, res) => {
  if (err) console.error('DB connection error:', err);
  else console.log('DB connected successfully');
});