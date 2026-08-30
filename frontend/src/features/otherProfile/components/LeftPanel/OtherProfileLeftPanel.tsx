import React from 'react';
import { formatDate } from '@/features/profile/utils/profileUtils';
import type { UserFull } from '@/types/model/models.ts';
import styles from './LeftPanel.module.css';

interface InfoBlockProps {
    label: string;
    value: string;
}

const InfoBlock: React.FC<InfoBlockProps> = ({ label, value }) => (
    <div className={styles.infoBlock}>
        <span className={styles.infoLabel}>{label}</span>
        <div className={styles.infoValue}>{value}</div>
    </div>
);

interface OtherProfileLeftPanelProps {
    user: UserFull;
}

export const OtherProfileLeftPanel: React.FC<OtherProfileLeftPanelProps> = ({ user }) => (
    <div className={styles.leftPanel}>
        <div className={styles.avatarRing}>
            <img src="/favicon.svg" alt="Avatar" className={styles.avatarImage} />
        </div>
        <InfoBlock label="Nickname" value={user.nickname} />
        <InfoBlock label="Created At" value={formatDate(user.created_at)} />
        <InfoBlock label="Role" value={user.role ?? 'user'} />
    </div>
);
