import React, { useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { DashboardOutlined, UserOutlined } from '@ant-design/icons';
import styles from './InspectorNavigation.module.css';

export const InspectorNavigation: React.FC = () => {
    const navigate = useNavigate();

    const handleNavigateDashboard = useCallback((): void => {
        navigate('/dashboard');
    }, [navigate]);

    const handleNavigateProfile = useCallback((): void => {
        navigate('/profile');
    }, [navigate]);

    return (
        <div className={styles.navContainer}>
            <button type="button" className={styles.navBtn} onClick={handleNavigateDashboard}>
                <DashboardOutlined /> Dashboard
            </button>
            <button type="button" className={styles.navBtn} onClick={handleNavigateProfile}>
                <UserOutlined /> Profile
            </button>
        </div>
    );
};