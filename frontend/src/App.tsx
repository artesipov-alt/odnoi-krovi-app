import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import OwnerDashboard from '@/pages/OwnerDashboard';
import ProfileForm from '@/components/ProfileForm'; // Исправлено: корректный импорт
import ChatMessages from '@/components/ChatMessages';
import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { useTelegramAuth } from '@/services/telegram.service';

const App: React.FC = () => {
  const { isRegistered } = useTelegramAuth();

  return (
    <Router>
      <Routes>
        <Route path="/register" element={<ProfileForm />} />
        <Route path="/owner" element={<OwnerDashboard />} />
        <Route path="/chat/:chatId" element={<ChatMessages />} />
        <Route
          path="/"
          element={isRegistered ? <Navigate to="/owner" /> : <ProfileForm />}
        />
      </Routes>
      <ToastContainer position="top-right" autoClose={3000} />
    </Router>
  );
};

export default App;