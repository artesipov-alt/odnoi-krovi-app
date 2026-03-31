import 'react-toastify/dist/ReactToastify.css';

import { FC } from 'react';
import { Navigate, Route, Routes } from 'react-router-dom';
import { Slide, ToastContainer } from 'react-toastify';

import { useAuth } from './hooks/useAuth';
import { useGetUserById } from './hooks/useGetUserById';
import About from './pages/about';
import AboutHistory from './pages/about/History';
import AboutTech from './pages/about/Tech';
import AboutThanks from './pages/about/Thanks';
import Adding from './pages/adding';
import Bonuses from './pages/bonuses';
import Owner from './pages/owner';
import Profile from './pages/profile';
import RecipientsList from './pages/recipientsList';
import Registration from './pages/registration';
import Search from './pages/search';

const App: FC = () => {
    const { userId, initialize } = useAuth();

    const { data: user, isLoading } = useGetUserById(userId);

    if (!user || isLoading) {
        return null;
    }

    return (
        <>
            <Routes>
                <Route path='/owner' element={<Owner userId={user.id} />} />
                <Route path='/profile' element={<Profile userId={user.id} />} />
                <Route path='/adding' element={<Adding userId={user.id} />} />
                <Route path='/search/:id' element={<Search userId={user.id} />} />
                <Route path='/recipientsList' element={<RecipientsList userId={user.id} />} />
                <Route path='/about' element={<About />} />
                <Route path='/about/history' element={<AboutHistory />} />
                <Route path='/about/thanks' element={<AboutThanks />} />
                <Route path='/about/tech' element={<AboutTech />} />
                <Route path='/bonuses' element={<Bonuses />} />
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
