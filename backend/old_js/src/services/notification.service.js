const TelegramBot = require('node-telegram-bot-api');
const bot = new TelegramBot(process.env.BOT_TOKEN, { polling: false });

async function sendNotification(telegramId, message) {
  try {
    await bot.sendMessage(telegramId, message, { parse_mode: 'HTML' });
  } catch (err) {
    console.error('Notification error:', err);
  }
}

module.exports = { sendNotification };