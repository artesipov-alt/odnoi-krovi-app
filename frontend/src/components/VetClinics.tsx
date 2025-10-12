import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useTelegramAuth } from '@/services/telegram.service';
import { toast } from 'react-toastify';
import styles from './VetClinics.module.css';

interface VetClinic {
  id: number;
  name: string;
  address: string;
  phone: string;
}

const VetClinics: React.FC = () => {
  const { initData } = useTelegramAuth();
  const [clinics, setClinics] = useState<VetClinic[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [formData, setFormData] = useState<Partial<VetClinic>>({});

  useEffect(() => {
    const fetchClinics = async () => {
      try {
        const response = await axios.get('/api/vet_clinics', {
          headers: { 'X-Telegram-Init-Data': initData }
        });
        setClinics(response.data);
      } catch (err: any) {
        setError('Failed to load vet clinics: ' + (err.response?.data?.message || err.message));
      } finally {
        setLoading(false);
      }
    };
    if (initData) {
      fetchClinics();
    }
  }, [initData]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.name || !formData.address || !formData.phone) {
      setError('All fields are required');
      return;
    }
    try {
      const response = await axios.post('/api/vet_clinics/create', formData, {
        headers: { 'X-Telegram-Init-Data': initData }
      });
      setClinics([...clinics, response.data]);
      setFormData({});
      toast.success('Vet clinic created');
    } catch (err: any) {
      toast.error('Create failed: ' + (err.response?.data?.message || err.message));
    }
  };

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingId || !formData.name || !formData.address || !formData.phone) {
      setError('All fields are required');
      return;
    }
    try {
      const response = await axios.put(`/api/vet_clinics/${editingId}`, formData, {
        headers: { 'X-Telegram-Init-Data': initData }
      });
      setClinics(clinics.map(clinic => (clinic.id === editingId ? response.data : clinic)));
      setEditingId(null);
      setFormData({});
      toast.success('Vet clinic updated');
    } catch (err: any) {
      toast.error('Update failed: ' + (err.response?.data?.message || err.message));
    }
  };

  const handleDelete = async (id: number) => {
    if (window.confirm('Are you sure you want to delete this vet clinic?')) {
      try {
        await axios.delete(`/api/vet_clinics/${id}`, {
          headers: { 'X-Telegram-Init-Data': initData }
        });
        setClinics(clinics.filter(clinic => clinic.id !== id));
        toast.success('Vet clinic deleted');
      } catch (err: any) {
        toast.error('Delete failed: ' + (err.response?.data?.message || err.message));
      }
    }
  };

  if (loading) return <p className={styles.loading}>Loading...</p>;
  if (error) return <p className={styles.error}>{error}</p>;

  return (
    <div className={styles.container}>
      <h3>Vet Clinics</h3>
      <form onSubmit={editingId ? handleUpdate : handleCreate}>
        <input
          name="name"
          value={formData.name || ''}
          onChange={handleInputChange}
          placeholder="Clinic Name"
          className={styles.input}
        />
        <input
          name="address"
          value={formData.address || ''}
          onChange={handleInputChange}
          placeholder="Address"
          className={styles.input}
        />
        <input
          name="phone"
          value={formData.phone || ''}
          onChange={handleInputChange}
          placeholder="Phone"
          className={styles.input}
        />
        <button type="submit" className={styles.button}>
          {editingId ? 'Update' : 'Create'}
        </button>
        {editingId && (
          <button
            type="button"
            onClick={() => {
              setEditingId(null);
              setFormData({});
            }}
            className={styles.button}
          >
            Cancel
          </button>
        )}
      </form>
      <div className={styles.clinicList}>
        {clinics.length > 0 ? (
          clinics.map(clinic => (
            <div key={clinic.id} className={styles.clinicItem}>
              <p>{clinic.name} - {clinic.address} - {clinic.phone}</p>
              <button
                onClick={() => {
                  setEditingId(clinic.id);
                  setFormData(clinic);
                }}
                className={styles.button}
              >
                Edit
              </button>
              <button onClick={() => handleDelete(clinic.id)} className={styles.button}>
                Delete
              </button>
            </div>
          ))
        ) : (
          <p className={styles.noClinics}>No vet clinics found</p>
        )}
      </div>
    </div>
  );
};

export default VetClinics;