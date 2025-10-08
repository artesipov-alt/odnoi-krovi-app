import React, { useState, useEffect } from 'react';
import axios, { AxiosError } from 'axios';
import { useTelegramAuth } from '../services/telegram.service';
import { YMaps, Map } from '@pbe/react-yandex-maps';

interface WebAppUser {
  id: number;
  first_name: string;
  last_name?: string;
}

interface Pet {
  id: number;
  name: string;
  type: string;
  photo_url: string;
  latitude: number;
  longitude: number;
  blood_group: string;
}

function OwnerDashboard() {
  const { user, initData } = useTelegramAuth();
  const [pets, setPets] = useState<Pet[]>([]);

  useEffect(() => {
    if (!initData) return;
    axios.get('/api/pets', { headers: { 'X-Telegram-Init-Data': initData } })
      .then(res => setPets(res.data as Pet[]))
      .catch((err: AxiosError) => console.error(err.message));
  }, [initData]);

  const updatePet = async (id: number) => {
    const name = prompt('New pet name:', '');
    const type = prompt('Pet type (cat/dog):', '');
    const latitude = parseFloat(prompt('Latitude:', '') || '0');
    const longitude = parseFloat(prompt('Longitude:', '') || '0');
    const blood_group = prompt('Blood group (A/B/AB/DEA1.1+/DEA1.1-):', '');

    if (!name || !type || !blood_group) {
      alert('All fields are required');
      return;
    }

    try {
      await axios.put(`/api/pets/${id}`, { name, type, latitude, longitude, blood_group }, {
        headers: { 'X-Telegram-Init-Data': initData }
      });
      alert('Pet updated!');
      axios.get('/api/pets', { headers: { 'X-Telegram-Init-Data': initData } })
        .then(res => setPets(res.data as Pet[]));
    } catch (err) {
      const error = err as AxiosError;
      alert(error.message || 'Update failed');
    }
  };

  const deletePet = async (id: number) => {
    try {
      await axios.delete(`/api/pets/${id}`, { headers: { 'X-Telegram-Init-Data': initData } });
      alert('Pet deleted!');
      axios.get('/api/pets', { headers: { 'X-Telegram-Init-Data': initData } })
        .then(res => setPets(res.data as Pet[]));
    } catch (err) {
      const error = err as AxiosError;
      alert(error.message || 'Deletion failed');
    }
  };

  return (
    <div className="container">
      <h2>Дашборд владельца</h2>
      <p>Полное имя: {user ? `${user.first_name} ${user.last_name || ''}`.trim() : 'Не авторизован'}</p>
      <h3>Мои питомцы</h3>
      {pets && pets.length > 0 ? (
        pets.map((pet: Pet) => (
          <div key={pet.id}>
            <p>{pet.name} ({pet.type})</p>
            <img src={pet.photo_url} alt={pet.name} style={{ maxWidth: '100px' }} />
            <button onClick={() => updatePet(pet.id)}>Обновить</button>
            <button onClick={() => deletePet(pet.id)}>Удалить</button>
          </div>
        ))
      ) : (
        <p>Питомцы не найдены</p>
      )}
      <YMaps>
        <Map defaultState={{ center: [55.75, 37.62], zoom: 10 }} width="100%" height="400px" />
      </YMaps>
    </div>
  );
}

export default OwnerDashboard;