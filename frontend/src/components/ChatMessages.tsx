import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import axios from 'axios';
import { useTelegram } from '../context/TelegramContext';
import { toast } from 'react-toastify';
import styles from './ChatMessages.module.css';

interface Message {
  id: number;
  sender_id: number;
  content: string;
  created_at: string;
}

const ChatMessages: React.FC = () => {
  const { chatId } = useParams<{ chatId: string }>();
  const { initData, user } = useTelegram();
  const [messages, setMessages] = useState<Message[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [newMessage, setNewMessage] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    const fetchMessages = async () => {
      try {
        const response = await axios.get('/api/chat_messages', {
          headers: { 'X-Telegram-Init-Data': initData || 'test_init_data' },
          params: { chat_id: chatId },
        });
        setMessages(response.data);
      } catch (err: any) {
        setError('Не удалось загрузить сообщения: ' + (err.response?.data?.error || err.message));
      } finally {
        setLoading(false);
      }
    };
    if (initData && chatId) {
      fetchMessages();
    }
  }, [initData, chatId]);

  const handleSendMessage = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newMessage.trim()) {
      setError('Сообщение не может быть пустым');
      return;
    }
    try {
      const response = await axios.post(
        '/api/chat_messages/send',
        { content: newMessage, chat_id: chatId },
        { headers: { 'X-Telegram-Init-Data': initData || 'test_init_data' } }
      );
      setMessages([...messages, response.data]);
      setNewMessage('');
      toast.success('Сообщение отправлено');
    } catch (err: any) {
      toast.error('Ошибка отправки: ' + (err.response?.data?.error || err.message));
    }
  };

  if (loading) return <p className={styles.loading}>Загрузка...</p>;
  if (error) return <p className={styles.error}>{error}</p>;

  return (
    <div className={styles.container}>
      <h3>Чат</h3>
      <button onClick={() => navigate('/owner')} className={styles.backButton}>
        Назад в дашборд
      </button>
      <div className={styles.messageList}>
        {messages.length > 0 ? (
          messages.map((message) => (
            <div
              key={message.id}
              className={`${styles.messageItem} ${
                message.sender_id === user?.id ? styles.sent : styles.received
              }`}
            >
              <p>{message.content}</p>
              <p className={styles.timestamp}>
                {new Date(message.created_at).toLocaleString()}
              </p>
            </div>
          ))
        ) : (
          <p className={styles.noMessages}>Сообщений пока нет</p>
        )}
      </div>
      <form onSubmit={handleSendMessage}>
        <input
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
          placeholder="Введите сообщение"
          className={styles.input}
        />
        <button type="submit" className={styles.button}>
          Отправить
        </button>
      </form>
    </div>
  );
};

export default ChatMessages;