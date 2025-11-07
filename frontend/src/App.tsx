import 'react-toastify/dist/ReactToastify.css';

import { FC } from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';
import { Slide, ToastContainer } from 'react-toastify';

import Owner from './pages/owner';
import Registration from './pages/registration';
import { useTelegram } from './TelegramProvider';

const App: FC = () => {
    const { isRegistered, user } = useTelegram();

    if (!user) {
        return null;
    }

    return (
        <>
            <Routes>
                <Route path='/owner' element={<Owner user={user} />} />
                <Route path='/registration' element={<Registration user={user} />} />
                <Route path='/' element={isRegistered ? <Navigate to='/owner' /> : <Registration user={user} />} />
            </Routes>
            <ToastContainer
                draggable
                theme='colored'
                hideProgressBar
                autoClose={3000}
                transition={Slide}
                position='top-right'
                closeOnClick={false}
            />
        </>
    );
};

export default App;
