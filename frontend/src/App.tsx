import React from 'react';
import { Routes, Route, Navigate } from 'react-router-dom'; // Удалён BrowserRouter, так как он уже в index.tsx
import OwnerDashboard from '@/pages/OwnerDashboard';
import ProfileForm from '@/components/ProfileForm';
import ChatMessages from '@/components/ChatMessages';
import { ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import { useTelegramAuth } from '@/services/telegram.service';

const App: React.FC = () => {
  const { isRegistered } = useTelegramAuth();

  return (
    <>
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
    </>
  );
};

export default App;