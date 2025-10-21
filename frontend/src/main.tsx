import './styles.css';

import ReactDOM from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';

import ErrorBoundary from 'components/ErrorBoundary';

import App from './App';
import { TelegramProvider } from './context/TelegramContext';

const rootElement = document.getElementById('root');

const root = ReactDOM.createRoot(rootElement!);
root.render(
    <TelegramProvider>
        <BrowserRouter>
            <ErrorBoundary>
                <App />
            </ErrorBoundary>
        </BrowserRouter>
    </TelegramProvider>,
);
