import React, { useState } from 'react';
import axios from 'axios';
import { useTelegram } from '../context/TelegramContext';
import { toast } from 'react-toastify';
import styles from './PetForm.module.css';

interface Pet {
  id?: number;
  name: string;
  type: string;
  photo_url: string;
  latitude: number;
  longitude: number;
  blood_group: string;
}

interface PetFormProps {
  pet?: Pet | null;
  onSave: () => void;
}

const PetForm: React.FC<PetFormProps> = ({ pet, onSave }) => {
  const { initData } = useTelegram();
  const [formData, setFormData] = useState<Pet>({
    name: pet?.name || '',
    type: pet?.type || '',
    photo_url: pet?.photo_url || '',
    latitude: pet?.latitude || 0,
    longitude: pet?.longitude || 0,
    blood_group: pet?.blood_group || ''
  });
  const [error, setError] = useState<string | null>(null);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: name === 'latitude' || name === 'longitude' ? parseFloat(value) : value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.name || !formData.type || !formData.blood_group) {
      setError('All fields are required');
      return;
    }

    try {
      if (pet?.id) {
        await axios.put(`/api/pets/${pet.id}`, formData, {
          headers: { 'X-Telegram-Init-Data': initData }
        });
        toast.success('Pet updated!');
      } else {
        await axios.post('/api/pets/create', formData, {
          headers: { 'X-Telegram-Init-Data': initData }
        });
        toast.success('Pet created!');
      }
      onSave();
    } catch (err: any) {
      toast.error('Failed to save pet: ' + (err.response?.data?.message || err.message));
    }
  };

  return (
    <div className={styles.container}>
      <h3>{pet?.id ? 'Edit Pet' : 'Add Pet'}</h3>
      {error && <p className={styles.error}>{error}</p>}
      <form onSubmit={handleSubmit}>
        <input
          name="name"
          value={formData.name}
          onChange={handleInputChange}
          placeholder="Pet Name"
          className={styles.input}
        />
        <select
          name="type"
          value={formData.type}
          onChange={handleInputChange}
          className={styles.input}
        >
          <option value="">Select Type</option>
          <option value="cat">Cat</option>
          <option value="dog">Dog</option>
        </select>
        <input
          name="photo_url"
          value={formData.photo_url}
          onChange={handleInputChange}
          placeholder="Photo URL"
          className={styles.input}
        />
        <input
          name="latitude"
          type="number"
          value={formData.latitude}
          onChange={handleInputChange}
          placeholder="Latitude"
          className={styles.input}
        />
        <input
          name="longitude"
          type="number"
          value={formData.longitude}
          onChange={handleInputChange}
          placeholder="Longitude"
          className={styles.input}
        />
        <select
          name="blood_group"
          value={formData.blood_group}
          onChange={handleInputChange}
          className={styles.input}
        >
          <option value="">Select Blood Group</option>
          <option value="A">A</option>
          <option value="B">B</option>
          <option value="AB">AB</option>
          <option value="DEA1.1+">DEA1.1+</option>
          <option value="DEA1.1- ">DEA1.1-</option>
        </select>
        <button type="submit" className={styles.button}>
          {pet?.id ? 'Update' : 'Create'}
        </button>
        <button type="button" onClick={onSave} className={styles.button}>
          Cancel
        </button>
      </form>
    </div>
  );
};

export default PetForm;
