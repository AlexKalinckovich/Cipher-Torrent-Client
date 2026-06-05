import React, { useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { UserOutlined } from '@ant-design/icons';
import styles from './ProfileLink.module.css';

export const ProfileLink: React.FC = () => {
    const navigate = useNavigate();

    const handleNavigate = useCallback((): void => {
        navigate('/profile');
    }, [navigate]);

    return (
        <div className={styles.container}>
            <button type="button" className={styles.btn} onClick={handleNavigate}>
                <UserOutlined /> Profile
            </button>
        </div>
    );
};