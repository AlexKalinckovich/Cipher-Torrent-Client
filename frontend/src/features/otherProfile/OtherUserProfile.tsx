import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Spin, Alert } from 'antd';
import { OtherProfileLeftPanel } from './components/LeftPanel/OtherProfileLeftPanel';
import { OtherProfileRightPanel } from './components/RightPanel/OtherProfileRightPanel';
import styles from './common.module.css';
import { userService } from '@/api/userService.ts';
import type { UserFull } from '@/types/model/models.ts';

export const OtherUserProfile: React.FC = () => {
    const { publicKey } = useParams<{ publicKey: string }>();
    const navigate = useNavigate();
    const [user, setUser] = useState<UserFull | null>(null);
    const [isLoading, setIsLoading] = useState<boolean>(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchUserProfile = async () => {
            if (!publicKey) {
                setError('Public key is required');
                setIsLoading(false);
                return;
            }

            try {
                const fetchedUser = await userService.getUserByPublicKey(publicKey);
                setUser(fetchedUser);
                setError(null);
            } catch (err) {
                setError('Failed to load user profile');
                console.error('Error fetching user profile:', err);
            } finally {
                setIsLoading(false);
            }
        };

        fetchUserProfile();
    }, [publicKey]);

    if (isLoading) {
        return (
            <div className={styles.pageContainer} style={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
                <Spin size="large" />
            </div>
        );
    }

    if (error || !user) {
        return (
            <div className={styles.pageContainer}>
                <Alert
                    message="Error"
                    description={error || 'User not found'}
                    type="error"
                    showIcon
                    action={
                        <button
                            onClick={() => navigate('/dashboard')}
                            style={{
                                background: 'none',
                                border: 'none',
                                color: '#1890ff',
                                cursor: 'pointer'
                            }}
                        >
                            Back to Dashboard
                        </button>
                    }
                />
            </div>
        );
    }

    return (
        <div className={styles.pageContainer}>
            <div className={styles.backgroundBlob1} />
            <div className={styles.backgroundBlob2} />
            <div className={styles.mainGrid}>
                <OtherProfileLeftPanel user={user} />
                <OtherProfileRightPanel stats={user.stats} />
            </div>
        </div>
    );
};
