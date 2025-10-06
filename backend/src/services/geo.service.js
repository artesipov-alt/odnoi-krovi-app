const axios = require('axios');

async function calculateDistance(lat1, lon1, lat2, lon2) {
  const response = await axios.get(`https://api.routing.yandex.net/v2/distancematrix?apikey=${process.env.YANDEX_API_KEY}&origins=${lat1},${lon1}&destinations=${lat2},${lon2}`);
  return response.data.rows[0].elements[0].distance.value / 1000; // км
}

module.exports = { calculateDistance };
