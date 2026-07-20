import React from 'react';
import { ProfileLeftPanel } from './components/LeftPanel/ProfileLeftPanel';
import { ProfileRightPanel } from './components/RightPanel/ProfileRightPanel';
import styles from './common.module.css';
import { useAuth } from '@/AuthContext.tsx';
import {Spin} from "antd";

export const Profile: React.FC = () => {
    const { user, isLoading } = useAuth();

    if (isLoading) {
        return (
            <div className={styles.pageContainer} style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                <Spin size="large" />
            </div>
        );
    }

    if (!user) {
        return <div className={styles.pageContainer}>Please log in to view your profile.</div>;
    }

    return (
        <div className={styles.pageContainer}>
            <div className={styles.backgroundBlob1} />
            <div className={styles.backgroundBlob2} />
            <div className={styles.mainGrid}>
                <ProfileLeftPanel user={user} />
                <ProfileRightPanel stats={user.stats} />
            </div>
        </div>
    );
};