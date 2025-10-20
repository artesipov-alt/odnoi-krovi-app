const crypto = require('crypto');

function validateInitData(initData, botToken) {
  const data = new URLSearchParams(initData);
  const hash = data.get('hash');
  data.delete('hash');
  const keys = Array.from(data.keys()).sort();
  const dataCheckString = keys.map(k => `${k}=${data.get(k)}`).join('\n');
  const secret = crypto.createHmac('sha256', 'WebAppData').update(botToken).digest();
  const calculatedHash = crypto.createHmac('sha256', secret).update(dataCheckString).digest('hex');
  return calculatedHash === hash;
}

module.exports = { validateInitData };