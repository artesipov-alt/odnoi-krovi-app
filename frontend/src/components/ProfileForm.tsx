import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import { useTelegram } from '../context/TelegramContext';
import { registerUser, checkProfile } from 'services/telegram.service';
import { toast } from 'react-toastify';
import styles from './ProfileForm.module.css';

interface RegisterData {
  telegram_id: number;
  role: 'pet_owner' | 'donor' | 'clinic_admin';
  full_name: string;
  phone: string;
  email: string;
  consent_pd: boolean;
}

const ProfileForm: React.FC = () => {
  const { user, initData, isRegistered, recheckRegistered } = useTelegram();
  const navigate = useNavigate();
  const [formData, setFormData] = useState<RegisterData>({
    telegram_id: user?.id || 314638947,
    role: 'pet_owner',
    full_name: user ? `${user.first_name} ${user.last_name || ''}`.trim() : 'Тест Пользователь',
    phone: '',
    email: '',
    consent_pd: false,
  });
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (isRegistered) {
      navigate('/owner');
    }
  }, [isRegistered, navigate]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
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
    setError(null);
    try {
      await registerUser(formData, initData);
      await recheckRegistered();
      const isProfileCreated = await checkProfile(initData);
      if (isProfileCreated) {
        toast.success('Регистрация прошла успешно');
        navigate('/owner');
      } else {
        throw new Error('Профиль не найден после регистрации');
      }
    } catch (err: any) {
      const errorMessage = err.response?.data?.error || err.message || 'Неизвестная ошибка';
      toast.error('Ошибка регистрации: ' + errorMessage);
      setError(errorMessage);
    } finally {
      setLoading(false);
    }
  };

  if (!user || !initData) {
    return <p>Загрузка данных Telegram...</p>;
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
          <option value="clinic_admin">Администратор клиники</option>
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
