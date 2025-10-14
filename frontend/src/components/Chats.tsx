import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { useTelegram } from '../context/TelegramContext';
import { toast } from 'react-toastify';
import styles from './Chats.module.css';

interface Chat {
  id: number;
  participant_id: number;
  participant_name: string;
  last_message: string;
  last_message_time: string;
}

const Chats: React.FC = () => {
  const { initData } = useTelegram();
  const [chats, setChats] = useState<Chat[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [newChatParticipantId, setNewChatParticipantId] = useState<number | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchChats = async () => {
      try {
        const response = await axios.get('/api/chats', {
          headers: { 'X-Telegram-Init-Data': initData },
        });
        setChats(response.data);
      } catch (err: any) {
        setError('Failed to load chats: ' + (err.response?.data?.message || err.message));
      } finally {
        setLoading(false);
      }
    };
    if (initData) {
      fetchChats();
    }
  }, [initData]);

  const handleCreateChat = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newChatParticipantId) {
      setError('Participant ID is required');
      return;
    }
    try {
      const response = await axios.post(
        '/api/chats/create',
        { participant_id: newChatParticipantId },
        { headers: { 'X-Telegram-Init-Data': initData } }
      );
      setChats([...chats, response.data]);
      setNewChatParticipantId(null);
      toast.success('Chat created');
    } catch (err: any) {
      toast.error('Create failed: ' + (err.response?.data?.message || err.message));
    }
  };

  const handleOpenChat = (chatId: number) => {
    navigate(`/chat/${chatId}`);
  };

  if (loading) return <p className={styles.loading}>Loading...</p>;
  if (error) return <p className={styles.error}>{error}</p>;

  return (
    <div className={styles.container}>
      <h3>Chats</h3>
      <form onSubmit={handleCreateChat}>
        <input
          type="number"
          value={newChatParticipantId || ''}
          onChange={(e) => setNewChatParticipantId(parseInt(e.target.value) || null)}
          placeholder="Participant Telegram ID"
          className={styles.input}
        />
        <button type="submit" className={styles.button}>
          Create Chat
        </button>
      </form>
      <div className={styles.chatList}>
        {chats.length > 0 ? (
          chats.map((chat) => (
            <div
              key={chat.id}
              className={styles.chatItem}
              onClick={() => handleOpenChat(chat.id)}
            >
              <p>
                {chat.participant_name} (ID: {chat.participant_id})
              </p>
              <p>Last Message: {chat.last_message || 'No messages yet'}</p>
              <p>Time: {new Date(chat.last_message_time).toLocaleString()}</p>
            </div>
          ))
        ) : (
          <p className={styles.noChats}>No chats found</p>
        )}
      </div>
    </div>
  );
};

export default Chats;