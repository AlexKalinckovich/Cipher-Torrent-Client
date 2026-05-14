import React, { useState, useCallback } from 'react';
import { EyeOutlined, EyeInvisibleOutlined } from '@ant-design/icons';
import { formatDate } from '@/features/profile/utils/profileUtils';
import type { UserFull } from '@/types/model/models.ts';
import styles from './LeftPanel.module.css';

interface InfoBlockProps {
    label: string;
    value: string;
}

export const InfoBlock: React.FC<InfoBlockProps> = ({ label, value }) => (
    <div className={styles.infoBlock}>
        <span className={styles.infoLabel}>{label}</span>
        <div className={styles.infoValue}>{value}</div>
    </div>
);

interface PublicKeyBlockProps {
    publicKey: string;
}

export const PublicKeyBlock: React.FC<PublicKeyBlockProps> = ({ publicKey }) => {
    const [isVisible, setIsVisible] = useState<boolean>(false);
    const toggleVisibility = useCallback(() => {
        setIsVisible((prev: boolean): boolean => {
            return !prev;
        });
    }, []);
    const displayClass: string = isVisible ? styles.keyTextVisible : styles.keyTextBlurred;
    const icon = isVisible ? <EyeInvisibleOutlined className={styles.keyIcon} /> : <EyeOutlined className={styles.keyIcon} />;
    return (
        <div className={styles.infoBlock}>
            <span className={styles.infoLabel}>Public Key (Peer ID)</span>
            <div className={styles.infoValue}>
                <div className={styles.keyContainer} onClick={toggleVisibility}>
                    <span className={displayClass}>{publicKey}</span>
                    {icon}
                </div>
            </div>
        </div>
    );
};

interface ProfileLeftPanelProps {
    user: UserFull;
}

export const ProfileLeftPanel: React.FC<ProfileLeftPanelProps> = ({ user }) => (
    <div className={styles.leftPanel}>
        <div className={styles.avatarRing}>
            <img src="/favicon.svg" alt="Avatar" className={styles.avatarImage} />
        </div>
        <InfoBlock label="Nickname" value={user.nickname} />
        <PublicKeyBlock publicKey={user.public_key} />
        <InfoBlock label="Created At" value={formatDate(user.created_at)} />
        <InfoBlock label="Role" value={user.role ?? 'user'} />
    </div>
);