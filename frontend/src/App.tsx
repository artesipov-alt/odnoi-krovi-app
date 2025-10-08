import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import ProfileForm from './components/ProfileForm.tsx';
import OwnerDashboard from './pages/OwnerDashboard.tsx';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<ProfileForm />} />
        <Route path="/owner" element={<OwnerDashboard />} />
      </Routes>
    </Router>
  );
}

export default App;