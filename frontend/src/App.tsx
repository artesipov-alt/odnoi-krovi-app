import 'react-toastify/dist/ReactToastify.css';

import { FC } from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';
import { Slide, ToastContainer } from 'react-toastify';

import { useAuth } from './hooks/useAuth';
// import { useTelegram } from './TelegramProvider';
import { useGetUserById } from './hooks/useGetUserById';
import Adding from './pages/adding';
import Owner from './pages/owner';
import Registration from './pages/registration';
import Search from './pages/search';

const App: FC = () => {
    // const { isRegistered, user } = useAuth();
    const { userId, initialize } = useAuth();

    const { data: user, isLoading } = useGetUserById(userId);

    if (!user || isLoading) {
        return null;
    }

    return (
        <>
            <Routes>
                <Route path='/owner' element={<Owner userId={user.id} />} />
                <Route path='/adding' element={<Adding userId={user.id} />} />
                <Route path='/search/:id' element={<Search userId={user.id} />} />
                <Route
                    path='/registration'
                    element={<Registration initialize={initialize} userId={user.id} fullName={user.fullName} />}
                />
                <Route
                    path='/'
                    element={
                        user.phone ? (
                            <Navigate to='/owner' />
                        ) : (
                            <Registration initialize={initialize} userId={user.id} fullName={user.fullName} />
                        )
                    }
                />
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
