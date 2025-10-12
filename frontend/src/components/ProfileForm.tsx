import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { useTelegramAuth, registerUser } from '@/services/telegram.service';
import { toast } from 'react-toastify';
import styles from './ProfileForm.module.css'; // Исправлено: импорт из ProfileForm.module.css

interface RegisterData {
  telegram_id: number;
  role: 'pet_owner' | 'donor';
  full_name: string;
  phone: string;
  email: string;
  consent_pd: boolean;
}

const ProfileForm: React.FC = () => {
  const { user, initData, isRegistered } = useTelegramAuth();
  const navigate = useNavigate();
  const [formData, setFormData] = useState<RegisterData>({
    telegram_id: user?.id || 0,
    role: 'pet_owner',
    full_name: user ? `${user.first_name} ${user.last_name || ''}`.trim() : '',
    phone: '',
    email: '',
    consent_pd: false,
  });
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false); // Добавлено: состояние загрузки

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => { // Исправлено: упрощена типизация
    const { name, value, type } = e.target;
    setFormData({
      ...formData,
      [name]: type === 'checkbox' ? (e.target as HTMLInputElement).checked : value,
    });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.role || !formData.full_name || !formData.phone || !formData.email || !formData.consent_pd) {
      setError('Все поля обязательны, включая согласие на обработку данных');
      return;
    }
    setLoading(true);
    try {
      await registerUser(formData);
      const profileResponse = await axios.get('/api/users/profile', {
        headers: { 'X-Telegram-Init-Data': initData },
      });
      if (profileResponse.data) {
        toast.success('Регистрация прошла успешно');
        navigate('/owner');
      }
    } catch (err: any) {
      toast.error('Ошибка регистрации: ' + (err.response?.data?.message || err.message));
    } finally {
      setLoading(false);
    }
  };

  if (!user || !initData) {
    return <p>Загрузка данных Telegram...</p>;
  }

  if (isRegistered) {
    navigate('/owner');
    return null;
  }

  return (
    <div className={styles.container}>
      <h2>Регистрация</h2>
      {error && <p className={styles.error}>{error}</p>}
      {loading && <p>Загрузка...</p>}
      <form onSubmit={handleSubmit}>
        <select
          name="role"
          value={formData.role}
          onChange={handleInputChange}
          className={styles.input}
        >
          <option value="pet_owner">Владелец питомца</option>
          <option value="donor">Донор</option>
        </select>
        <input
          name="full_name"
          value={formData.full_name}
          onChange={handleInputChange}
          placeholder="Полное имя"
          className={styles.input}
        />
        <input
          name="phone"
          value={formData.phone}
          onChange={handleInputChange}
          placeholder="Телефон"
          className={styles.input}
        />
        <input
          name="email"
          value={formData.email}
          onChange={handleInputChange}
          placeholder="Email"
          className={styles.input}
        />
        <label>
          <input
            name="consent_pd"
            type="checkbox"
            checked={formData.consent_pd}
            onChange={handleInputChange}
          />
          Согласие на обработку персональных данных
        </label>
        <button type="submit" className={styles.button} disabled={loading}>
          Зарегистрироваться
        </button>
      </form>
    </div>
  );
};

export default ProfileForm;