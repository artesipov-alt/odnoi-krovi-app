import React, { useState, useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import axios from 'axios';
import { useTelegram } from '../context/TelegramContext';
import { YMaps, Map } from '@pbe/react-yandex-maps';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import Profile from '@/components/Profile';
import PetForm from '@/components/PetForm';
import VetClinics from '@/components/VetClinics';
import BloodStocks from '@/components/BloodStocks';
import BloodSearches from '@/components/BloodSearches';
import Donations from '@/components/Donations';
import styles from './OwnerDashboard.module.css';

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
  const { user, initData, isRegistered } = useTelegram();
  const [pets, setPets] = useState<Pet[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [editingPet, setEditingPet] = useState<Pet | null>(null);
  const [showPetForm, setShowPetForm] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    if (!isRegistered) {
      navigate('/register');
      return;
    }

    const fetchPets = async () => {
      try {
        const response = await axios.get('/api/pets', {
          headers: { 'X-Telegram-Init-Data': initData || 'test_init_data' },
        });
        setPets(response.data as Pet[]);
      } catch (err: any) {
        setError('Не удалось загрузить питомцев: ' + (err.response?.data?.error || err.message));
      } finally {
        setLoading(false);
      }
    };

    fetchPets();
  }, [initData, isRegistered, navigate]);

  const handleSavePet = async () => {
    try {
      const response = await axios.get('/api/pets', {
        headers: { 'X-Telegram-Init-Data': initData || 'test_init_data' },
      });
      setPets(response.data as Pet[]);
      setShowPetForm(false);
      setEditingPet(null);
    } catch (err: any) {
      toast.error('Не удалось обновить список питомцев: ' + (err.response?.data?.error || err.message));
    }
  };

  const deletePet = async (id: number) => {
    try {
      await axios.delete(`/api/pets/${id}`, {
        headers: { 'X-Telegram-Init-Data': initData || 'test_init_data' },
      });
      toast.success('Питомец удалён!');
      const response = await axios.get('/api/pets', {
        headers: { 'X-Telegram-Init-Data': initData || 'test_init_data' },
      });
      setPets(response.data as Pet[]);
    } catch (err: any) {
      toast.error('Ошибка удаления: ' + (err.response?.data?.error || err.message));
    }
  };

  if (loading) return <p className={styles.loading}>Загрузка...</p>;
  if (error) return <p className={styles.error}>{error}</p>;

  return (
    <div className={styles.container}>
      <h2>Дашборд владельца</h2>
      <p className={styles.userInfo}>
        Полное имя: {user ? `${user.first_name} ${user.last_name || ''}`.trim() : 'Не авторизован'}
      </p>
      <Link to="/chats">
        <button className={styles.button}>Перейти к чатам</button>
      </Link>
      <Profile />
      <VetClinics />
      <BloodStocks />
      <BloodSearches />
      <Donations />
      <h3>Мои питомцы</h3>
      {showPetForm ? (
        <PetForm pet={editingPet} onSave={handleSavePet} />
      ) : (
        <>
          <button
            onClick={() => setShowPetForm(true)}
            className={styles.button}
          >
            Добавить питомца
          </button>
          {pets.length > 0 ? (
            <div className={styles.petList}>
              {pets.map((pet) => (
                <div key={pet.id} className={styles.petItem}>
                  <p>{pet.name} ({pet.type})</p>
                  <img src={pet.photo_url} alt={pet.name} className={styles.petImage} />
                  <p>Blood Group: {pet.blood_group}</p>
                  <button
                    onClick={() => {
                      setEditingPet(pet);
                      setShowPetForm(true);
                    }}
                    className={styles.button}
                  >
                    Обновить
                  </button>
                  <button
                    onClick={() => deletePet(pet.id)}
                    className={styles.button}
                  >
                    Удалить
                  </button>
                </div>
              ))}
            </div>
          ) : (
            <p className={styles.noPets}>Питомцы не найдены</p>
          )}
        </>
      )}
      <YMaps>
        <Map defaultState={{ center: [55.75, 37.62], zoom: 10 }} width="100%" height="400px" />
      </YMaps>
      <ToastContainer position="top-right" autoClose={3000} />
    </div>
  );
}

export default OwnerDashboard;