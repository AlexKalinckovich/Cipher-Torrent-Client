import React from 'react';
import { Typography } from 'antd';
import { CloseCircleOutlined } from '@ant-design/icons';
import type { ErrorStateProps } from '@/features/torrents/types/dashboardTypes';
import styles from './ErrorState.module.css';

export const ErrorState: React.FC<ErrorStateProps> = ({ error }) => {
    return (
        <div style={{ display: 'flex', justifyContent: 'center', height: '100vh', alignItems: 'center' }}>
            <div className={styles.errorState}>
                <CloseCircleOutlined style={{ fontSize: 64, color: '#f87171', marginBottom: 24 }} />
                <Typography.Title level={3} style={{ color: '#fff' }}>Failed to load torrents</Typography.Title>
                <Typography.Text style={{ color: '#c4b5fd' }}>{error.message}</Typography.Text>
            </div>
        </div>
    );
};
