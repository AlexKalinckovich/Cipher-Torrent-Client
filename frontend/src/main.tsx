import './index.css'
import {App} from './App.tsx'
import React from 'react';
import ReactDOM from 'react-dom/client';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { api } from './api/axiosClient';
import { MockServer } from './mocks/mock-servert.ts';
const queryClient = new QueryClient();

const initializeEnvironment = (): void => {
    new MockServer(api);
    renderApp();
};

const renderApp = (): void => {
    ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
        <React.StrictMode>
            <QueryClientProvider client={queryClient}>
                <App />
            </QueryClientProvider>
        </React.StrictMode>
    );
};

initializeEnvironment();