import React, { useState } from 'react';
import axios, { AxiosError } from 'axios';
import { useTelegramAuth } from '../services/telegram.service';

interface WebAppUser {
  id: number;
  first_name: string;
  last_name?: string;
}

function ProfileForm() {
  const { user, initData } = useTelegramAuth();
  const [role, setRole] = useState('');
  const [full_name, setFullName] = useState(user ? `${user.first_name} ${user.last_name || ''}`.trim() : '');
  const [phone, setPhone] = useState('');
  const [email, setEmail] = useState('');
  const [consent_pd, setConsentPd] = useState(false);

  const handleSubmit = async () => {
    if (!consent_pd) {
      alert('Consent required');
      return;
    }
    if (!user) {
      alert('User not authenticated');
      return;
    }
    try {
      await axios.post('/api/users/register', {
        telegram_id: user.id,
        role,
        full_name,
        phone,
        email,
        consent_pd
      }, {
        headers: { 'X-Telegram-Init-Data': initData }
      });
      alert('Registered!');
    } catch (err) {
      const error = err as AxiosError;
      alert(error.message || 'Registration failed');
    }
  };

  return (
    <div className="container">
      <h2>Регистрация</h2>
      <select value={role} onChange={(e: React.ChangeEvent<HTMLSelectElement>) => setRole(e.target.value)}>
        <option value="">Выберите роль</option>
        <option value="pet_owner">Владелец питомца</option>
        <option value="vet_clinic">Ветклиника</option>
        <option value="animal_center">Центр содержания</option>
      </select>
      <input
        type="text"
        placeholder="Полное имя"
        value={full_name}
        onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFullName(e.target.value)}
      />
      <input
        type="text"
        placeholder="Телефон"
        value={phone}
        onChange={(e: React.ChangeEvent<HTMLInputElement>) => setPhone(e.target.value)}
      />
      <input
        type="email"
        placeholder="Email"
        value={email}
        onChange={(e: React.ChangeEvent<HTMLInputElement>) => setEmail(e.target.value)}
      />
      <input
        type="checkbox"
        checked={consent_pd}
        onChange={(e: React.ChangeEvent<HTMLInputElement>) => setConsentPd(e.target.checked)}
      /> Согласие на обработку ПДн
      <button onClick={handleSubmit}>Зарегистрироваться</button>
    </div>
  );
}

export default ProfileForm;