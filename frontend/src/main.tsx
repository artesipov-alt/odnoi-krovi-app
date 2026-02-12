import './styles.less';

import { QueryClientProvider } from '@tanstack/react-query';
import ReactDOM from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';

import { queryClient } from 'api/queryClient';
import ErrorBoundary from 'components/ErrorBoundary/ErrorBoundary';

import App from './App';
import { TelegramProvider } from './TelegramProvider';

const rootElement = document.getElementById('root');

const root = ReactDOM.createRoot(rootElement!);
root.render(
    <TelegramProvider>
        <BrowserRouter>
            <ErrorBoundary>
                <QueryClientProvider client={queryClient}>
                    <App />
                </QueryClientProvider>
            </ErrorBoundary>
        </BrowserRouter>
    </TelegramProvider>,
);
