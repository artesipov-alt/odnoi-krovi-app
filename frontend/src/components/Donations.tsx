import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useTelegram } from '../context/TelegramContext';
import { toast } from 'react-toastify';
import styles from './Donations.module.css';

interface Donation {
  id: number;
  pet_id: number;
  clinic_id: number;
  blood_group: string;
  quantity: number;
  donation_date: string;
}

const Donations: React.FC = () => {
  const { initData } = useTelegram();
  const [donations, setDonations] = useState<Donation[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [formData, setFormData] = useState<Partial<Donation>>({});

  useEffect(() => {
    const fetchDonations = async () => {
      try {
        const response = await axios.get('/api/donations', {
          headers: { 'X-Telegram-Init-Data': initData },
        });
        setDonations(response.data);
      } catch (err: any) {
        setError('Failed to load donations: ' + (err.response?.data?.message || err.message));
      } finally {
        setLoading(false);
      }
    };
    if (initData) {
      fetchDonations();
    }
  }, [initData]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: name === 'pet_id' || name === 'clinic_id' || name === 'quantity' ? parseInt(value) : value });
  };

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.pet_id || !formData.clinic_id || !formData.blood_group || !formData.quantity || !formData.donation_date) {
      setError('All fields are required');
      return;
    }
    try {
      const response = await axios.post('/api/donations/plan', formData, {
        headers: { 'X-Telegram-Init-Data': initData },
      });
      setDonations([...donations, response.data]);
      setFormData({});
      toast.success('Donation created');
    } catch (err: any) {
      toast.error('Create failed: ' + (err.response?.data?.message || err.message));
    }
  };

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingId || !formData.pet_id || !formData.clinic_id || !formData.blood_group || !formData.quantity || !formData.donation_date) {
      setError('All fields are required');
      return;
    }
    try {
      const response = await axios.put(`/api/donations/${editingId}/status`, formData, {
        headers: { 'X-Telegram-Init-Data': initData },
      });
      setDonations(donations.map(donation => (donation.id === editingId ? response.data : donation)));
      setEditingId(null);
      setFormData({});
      toast.success('Donation updated');
    } catch (err: any) {
      toast.error('Update failed: ' + (err.response?.data?.message || err.message));
    }
  };

  const handleDelete = async (id: number) => {
    try {
      await axios.delete(`/api/donations/${id}`, {
        headers: { 'X-Telegram-Init-Data': initData },
      });
      setDonations(donations.filter(donation => donation.id !== id));
      toast.success('Donation deleted');
    } catch (err: any) {
      toast.error('Delete failed: ' + (err.response?.data?.message || err.message));
    }
  };

  if (loading) return <p className={styles.loading}>Loading...</p>;
  if (error) return <p className={styles.error}>{error}</p>;

  return (
    <div className={styles.container}>
      <h3>Donations</h3>
      <form onSubmit={editingId ? handleUpdate : handleCreate}>
        <input
          name="pet_id"
          type="number"
          value={formData.pet_id || ''}
          onChange={handleInputChange}
          placeholder="Pet ID"
          className={styles.input}
        />
        <input
          name="clinic_id"
          type="number"
          value={formData.clinic_id || ''}
          onChange={handleInputChange}
          placeholder="Clinic ID"
          className={styles.input}
        />
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
          name="donation_date"
          type="date"
          value={formData.donation_date || ''}
          onChange={handleInputChange}
          placeholder="Donation Date"
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
      <div className={styles.donationList}>
        {donations.length > 0 ? (
          donations.map(donation => (
            <div key={donation.id} className={styles.donationItem}>
              <p>
                Pet ID: {donation.pet_id}, Clinic ID: {donation.clinic_id}, Blood Group: {donation.blood_group}, Quantity: {donation.quantity} ml, Date: {donation.donation_date}
              </p>
              <button
                onClick={() => {
                  setEditingId(donation.id);
                  setFormData(donation);
                }}
                className={styles.button}
              >
                Edit
              </button>
              <button onClick={() => handleDelete(donation.id)} className={styles.button}>
                Delete
              </button>
            </div>
          ))
        ) : (
          <p className={styles.noDonations}>No donations found</p>
        )}
      </div>
    </div>
  );
};

export default Donations;