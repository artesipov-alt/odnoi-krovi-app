import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useTelegram } from '../context/TelegramContext';
import { toast } from 'react-toastify';
import styles from './Profile.module.css';

interface ProfileData {
  telegram_id: number;
  role: string;
  full_name: string;
  phone: string;
  email: string;
  consent_pd: boolean;
}

const Profile: React.FC = () => {
  const { initData, isRegistered } = useTelegram();
  const [profile, setProfile] = useState<ProfileData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [editing, setEditing] = useState(false);
  const [formData, setFormData] = useState<Partial<ProfileData>>({});

  useEffect(() => {
    const fetchProfile = async () => {
      if (!isRegistered) {
        setLoading(false);
        return;
      }
      try {
        const response = await axios.get('/api/users/profile', {
          headers: { 'X-Telegram-Init-Data': initData },
        });
        setProfile(response.data);
      } catch (err: any) {
        setError('Не удалось загрузить профиль: ' + (err.response?.data?.message || err.message));
      } finally {
        setLoading(false);
      }
    };
    fetchProfile();
  }, [initData, isRegistered]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const response = await axios.put('/api/users/profile', formData, {
        headers: { 'X-Telegram-Init-Data': initData },
      });
      setProfile(response.data);
      setEditing(false);
      toast.success('Профиль обновлён');
    } catch (err: any) {
      toast.error('Ошибка обновления: ' + (err.response?.data?.message || err.message));
    }
  };

  const handleDelete = async () => {
    try {
      await axios.delete('/api/users/profile', {
        headers: { 'X-Telegram-Init-Data': initData },
      });
      toast.success('Профиль удалён');
      window.location.href = '/register';
    } catch (err: any) {
      toast.error('Ошибка удаления: ' + (err.response?.data?.message || err.message));
    }
  };

  if (loading) return <p className={styles.loading}>Загрузка...</p>;
  if (error) return <p className={styles.error}>{error}</p>;
  if (!isRegistered || !profile)
    return <p className={styles.noProfile}>Профиль не найден. Пожалуйста, зарегистрируйтесь.</p>;

  return (
    <div className={styles.container}>
      <h3>Профиль</h3>
      {editing ? (
        <form onSubmit={handleUpdate}>
          <select
            name="role"
            value={formData.role || profile.role}
            onChange={handleInputChange}
            className={styles.input}
          >
            <option value="pet_owner">Владелец питомца</option>
            <option value="donor">Донор</option>
            <option value="clinic_admin">Администратор клиники</option>
          </select>
          <input
            name="full_name"
            value={formData.full_name || profile.full_name}
            onChange={handleInputChange}
            placeholder="Полное имя"
            className={styles.input}
          />
          <input
            name="phone"
            value={formData.phone || profile.phone}
            onChange={handleInputChange}
            placeholder="Телефон"
            className={styles.input}
          />
          <input
            name="email"
            value={formData.email || profile.email}
            onChange={handleInputChange}
            placeholder="Email"
            className={styles.input}
          />
          <button type="submit" className={styles.button}>
            Сохранить
          </button>
          <button
            type="button"
            onClick={() => setEditing(false)}
            className={styles.button}
          >
            Отмена
          </button>
        </form>
      ) : (
        <div>
          <p>Роль: {profile.role === 'pet_owner' ? 'Владелец питомца' : profile.role === 'donor' ? 'Донор' : 'Администратор клиники'}</p>
          <p>Полное имя: {profile.full_name}</p>
          <p>Телефон: {profile.phone}</p>
          <p>Email: {profile.email}</p>
          <p>Согласие на обработку данных: {profile.consent_pd ? 'Да' : 'Нет'}</p>
          <button onClick={() => setEditing(true)} className={styles.button}>
            Редактировать
          </button>
          <button onClick={handleDelete} className={styles.button}>
            Удалить
          </button>
        </div>
      )}
    </div>
  );
};

export default Profile;