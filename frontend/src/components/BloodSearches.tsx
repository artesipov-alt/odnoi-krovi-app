import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useTelegram } from '../context/TelegramContext';
import { toast } from 'react-toastify';
import styles from './BloodSearches.module.css';

interface BloodSearch {
  id: number;
  blood_group: string;
  quantity: number;
  location: string;
  status: string;
}

const BloodSearches: React.FC = () => {
  const { initData } = useTelegram();
  const [searches, setSearches] = useState<BloodSearch[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [formData, setFormData] = useState<Partial<BloodSearch>>({});

  useEffect(() => {
    const fetchSearches = async () => {
      try {
        const response = await axios.get('/api/blood_searches', {
          headers: { 'X-Telegram-Init-Data': initData }
        });
        setSearches(response.data);
      } catch (err: any) {
        setError('Failed to load blood searches: ' + (err.response?.data?.message || err.message));
      } finally {
        setLoading(false);
      }
    };
    if (initData) {
      fetchSearches();
    }
  }, [initData]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: name === 'quantity' ? parseInt(value) : value });
  };

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.blood_group || !formData.quantity || !formData.location || !formData.status) {
      setError('All fields are required');
      return;
    }
    try {
      const response = await axios.post('/api/blood_searches/create', formData, {
        headers: { 'X-Telegram-Init-Data': initData }
      });
      setSearches([...searches, response.data]);
      setFormData({});
      toast.success('Blood search created');
    } catch (err: any) {
      toast.error('Create failed: ' + (err.response?.data?.message || err.message));
    }
  };

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingId || !formData.blood_group || !formData.quantity || !formData.location || !formData.status) {
      setError('All fields are required');
      return;
    }
    try {
      const response = await axios.put(`/api/blood_searches/${editingId}`, formData, {
        headers: { 'X-Telegram-Init-Data': initData }
      });
      setSearches(searches.map(search => (search.id === editingId ? response.data : search)));
      setEditingId(null);
      setFormData({});
      toast.success('Blood search updated');
    } catch (err: any) {
      toast.error('Update failed: ' + (err.response?.data?.message || err.message));
    }
  };

  const handleDelete = async (id: number) => {
    if (window.confirm('Are you sure you want to delete this blood search?')) {
      try {
        await axios.delete(`/api/blood_searches/${id}`, {
          headers: { 'X-Telegram-Init-Data': initData }
        });
        setSearches(searches.filter(search => search.id !== id));
        toast.success('Blood search deleted');
      } catch (err: any) {
        toast.error('Delete failed: ' + (err.response?.data?.message || err.message));
      }
    }
  };

  if (loading) return <p className={styles.loading}>Loading...</p>;
  if (error) return <p className={styles.error}>{error}</p>;

  return (
    <div className={styles.container}>
      <h3>Blood Searches</h3>
      <form onSubmit={editingId ? handleUpdate : handleCreate}>
        <select
          name="blood_group"
          value={formData.blood_group || ''}
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
        <input
          name="quantity"
          type="number"
          value={formData.quantity || ''}
          onChange={handleInputChange}
          placeholder="Quantity (ml)"
          className={styles.input}
        />
        <input
          name="location"
          value={formData.location || ''}
          onChange={handleInputChange}
          placeholder="Location"
          className={styles.input}
        />
        <select
          name="status"
          value={formData.status || ''}
          onChange={handleInputChange}
          className={styles.input}
        >
          <option value="">Select Status</option>
          <option value="open">Open</option>
          <option value="closed">Closed</option>
        </select>
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
      <div className={styles.searchList}>
        {searches.length > 0 ? (
          searches.map(search => (
            <div key={search.id} className={styles.searchItem}>
              <p>
                Blood Group: {search.blood_group}, Quantity: {search.quantity} ml, Location: {search.location}, Status: {search.status}
              </p>
              <button
                onClick={() => {
                  setEditingId(search.id);
                  setFormData(search);
                }}
                className={styles.button}
              >
                Edit
              </button>
              <button onClick={() => handleDelete(search.id)} className={styles.button}>
                Delete
              </button>
            </div>
          ))
        ) : (
          <p className={styles.noSearches}>No blood searches found</p>
        )}
      </div>
    </div>
  );
};

export default BloodSearches;