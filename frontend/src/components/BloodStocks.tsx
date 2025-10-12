import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { useTelegramAuth } from '@/services/telegram.service';
import { toast } from 'react-toastify';
import { Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, Button } from '@mui/material';
import styles from './BloodStocks.module.css';

interface BloodStock {
  id: number;
  clinic_id: number;
  blood_group: string;
  quantity: number;
  expiration_date: string;
}

interface UserProfile {
  role: 'pet_owner' | 'donor' | 'clinic_admin';
}

const BloodStocks: React.FC = () => {
  const { initData } = useTelegramAuth();
  const [stocks, setStocks] = useState<BloodStock[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [formData, setFormData] = useState<Partial<BloodStock>>({});
  const [userRole, setUserRole] = useState<string | null>(null);
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false);
  const [deleteId, setDeleteId] = useState<number | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const profileResponse = await axios.get('/api/users/profile', {
          headers: { 'X-Telegram-Init-Data': initData },
        });
        setUserRole(profileResponse.data.role);

        const stocksResponse = await axios.get('/api/blood_stocks', {
          headers: { 'X-Telegram-Init-Data': initData },
        });
        setStocks(stocksResponse.data);
      } catch (err: any) {
        setError('Failed to load data: ' + (err.response?.data?.message || err.message));
      } finally {
        setLoading(false);
      }
    };
    if (initData) {
      fetchData();
    }
  }, [initData]);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: name === 'clinic_id' || name === 'quantity' ? parseInt(value) || '' : value });
  };

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.clinic_id || !formData.blood_group || !formData.quantity || !formData.expiration_date) {
      setError('All fields are required');
      return;
    }
    try {
      const response = await axios.post('/api/blood_stocks/create', formData, {
        headers: { 'X-Telegram-Init-Data': initData },
      });
      setStocks([...stocks, response.data]);
      setFormData({});
      toast.success('Blood stock created');
    } catch (err: any) {
      toast.error('Create failed: ' + (err.response?.data?.message || err.message));
    }
  };

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingId || !formData.clinic_id || !formData.blood_group || !formData.quantity || !formData.expiration_date) {
      setError('All fields are required');
      return;
    }
    try {
      const response = await axios.put(`/api/blood_stocks/${editingId}`, formData, {
        headers: { 'X-Telegram-Init-Data': initData },
      });
      setStocks(stocks.map((stock) => (stock.id === editingId ? response.data : stock)));
      setEditingId(null);
      setFormData({});
      toast.success('Blood stock updated');
    } catch (err: any) {
      toast.error('Update failed: ' + (err.response?.data?.message || err.message));
    }
  };

  const handleDelete = async () => {
    if (!deleteId) return;
    try {
      await axios.delete(`/api/blood_stocks/${deleteId}`, {
        headers: { 'X-Telegram-Init-Data': initData },
      });
      setStocks(stocks.filter((stock) => stock.id !== deleteId));
      toast.success('Blood stock deleted');
    } catch (err: any) {
      toast.error('Delete failed: ' + (err.response?.data?.message || err.message));
    } finally {
      setDeleteDialogOpen(false);
      setDeleteId(null);
    }
  };

  const handleBook = async (id: number) => {
    try {
      await axios.post('/api/blood_stocks/book', { id }, {
        headers: { 'X-Telegram-Init-Data': initData },
      });
      toast.success('Blood stock booked');
      const response = await axios.get('/api/blood_stocks', {
        headers: { 'X-Telegram-Init-Data': initData },
      });
      setStocks(response.data);
    } catch (err: any) {
      toast.error('Book failed: ' + (err.response?.data?.message || err.message));
    }
  };

  const openDeleteDialog = (id: number) => {
    setDeleteId(id);
    setDeleteDialogOpen(true);
  };

  const closeDeleteDialog = () => {
    setDeleteDialogOpen(false);
    setDeleteId(null);
  };

  if (loading) return <p className={styles.loading}>Loading...</p>;
  if (error) return <p className={styles.error}>{error}</p>;

  return (
    <div className={styles.container}>
      <h3>Blood Stocks</h3>
      {userRole === 'clinic_admin' ? (
        <form onSubmit={editingId ? handleUpdate : handleCreate}>
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
            name="expiration_date"
            type="date"
            value={formData.expiration_date || ''}
            onChange={handleInputChange}
            placeholder="Expiration Date"
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
      ) : (
        <p className={styles.noAccess}>Only clinic admins can manage blood stocks.</p>
      )}
      <div className={styles.stockList}>
        {stocks.length > 0 ? (
          stocks.map((stock) => (
            <div key={stock.id} className={styles.stockItem}>
              <p>
                Clinic ID: {stock.clinic_id}, Blood Group: {stock.blood_group}, Quantity: {stock.quantity} ml,
                Expires: {new Date(stock.expiration_date).toLocaleDateString()}
              </p>
              {userRole === 'clinic_admin' && (
                <>
                  <button
                    onClick={() => {
                      setEditingId(stock.id);
                      setFormData(stock);
                    }}
                    className={styles.button}
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => openDeleteDialog(stock.id)}
                    className={styles.button}
                  >
                    Delete
                  </button>
                </>
              )}
              <button onClick={() => handleBook(stock.id)} className={styles.button}>
                Book
              </button>
            </div>
          ))
        ) : (
          <p className={styles.noStocks}>No blood stocks found</p>
        )}
      </div>
      <Dialog
        open={deleteDialogOpen}
        onClose={closeDeleteDialog}
        aria-labelledby="delete-dialog-title"
        aria-describedby="delete-dialog-description"
      >
        <DialogTitle id="delete-dialog-title">Confirm Delete</DialogTitle>
        <DialogContent>
          <DialogContentText id="delete-dialog-description">
            Are you sure you want to delete this blood stock?
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={closeDeleteDialog} color="primary">
            Cancel
          </Button>
          <Button onClick={handleDelete} color="error" autoFocus>
            Delete
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};

export default BloodStocks;