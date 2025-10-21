import 'react-toastify/dist/ReactToastify.css';

import OwnerDashboard from 'pages/OwnerDashboard';
import React from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';
import { ToastContainer } from 'react-toastify';

import ChatMessages from 'components/ChatMessages';
import Chats from 'components/Chats';
import ProfileForm from 'components/ProfileForm';

import { useTelegram } from './context/TelegramContext';

const App: React.FC = () => {
    const { isRegistered } = useTelegram();

    return (
        <>
            <Routes>
                <Route path='/register' element={<ProfileForm />} />
                <Route path='/owner' element={<OwnerDashboard />} />
                <Route path='/chats' element={<Chats />} />
                <Route path='/chat/:chatId' element={<ChatMessages />} />
                <Route path='/' element={isRegistered ? <Navigate to='/owner' /> : <ProfileForm />} />
            </Routes>
            <ToastContainer position='top-right' autoClose={3000} />
        </>
    );
};

export default App;
