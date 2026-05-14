import React, { useState } from 'react';
import type { UserFull } from '../../types/model/models.ts';
import { ProfileLeftPanel } from './components/LeftPanel/ProfileLeftPanel';
import { ProfileRightPanel } from './components/RightPanel/ProfileRightPanel';
import { ProfileReputationHeader } from './components/ReputationHeader/ProfileReputationHeader';
import styles from './common.module.css';

const createMockData = (): UserFull => ({
    id: 777,
    email: 'phantom@cipher.net',
    public_key: 'Ed25519_PubKey_Phantom_X9A4B2K8M3L1P0Q5R',
    nickname: 'PhantomRunner',
    created_at: '2026-05-11T23:00:00Z',
    role: 'admin',
    stats: {
        total_uploaded_bytes: 45600000000,
        total_downloaded_bytes: 12400000000,
        reputation_score: 0.94,
        signed_torrents_count: 128,
        active_torrents_count: 14,
        peers_trusted_count: 56
    }
});

const useMockUserProfile = (): UserFull => {
    const [user] = useState<UserFull>(createMockData);
    return user;
};

export const Profile: React.FC = () => {
    const user = useMockUserProfile();
    return (
        <div className={styles.pageContainer}>
            <div className={styles.backgroundBlob1} />
            <div className={styles.backgroundBlob2} />
            <ProfileReputationHeader score={user.stats?.reputation_score} />
            <div className={styles.mainGrid}>
                <ProfileLeftPanel user={user} />
                <ProfileRightPanel stats={user.stats} />
            </div>
        </div>
    );
};