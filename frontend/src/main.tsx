import './index.css';
import { App } from './App.tsx';
import React from 'react';
import ReactDOM from 'react-dom/client';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ConfigProvider, theme } from 'antd';
import { api } from './api/axiosClient';
import { MockServer } from './mocks/mock-servert.ts';

const queryClient = new QueryClient();


const darkTheme = {
    algorithm: theme.darkAlgorithm,
    token: {
        colorPrimary: '#b37feb',     
        colorBgBase: '#0f0f14',       
        colorBgContainer: '#1a1a24',  
        colorBorderSecondary: '#2a2a35',
        colorTextBase: '#e5e5eb',
        colorTextSecondary: '#a0a0b0',
        borderRadius: 8,
    },
    components: {
        Table: {
            headerBg: '#1f1f2a',
            rowHoverBg: '#2a2a35',
            borderColor: '#2a2a35',
        },
        Card: {
            actionsBg: '#1a1a24',
        },
    },
};


const initializeEnvironment = (): void => {
    new MockServer(api);
    renderApp();
};

const renderApp = (): void => {
    ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
        <React.StrictMode>
            <QueryClientProvider client={queryClient}>
                <ConfigProvider theme={darkTheme}>
                    <App />
                </ConfigProvider>
            </QueryClientProvider>
        </React.StrictMode>
    );
};

initializeEnvironment();