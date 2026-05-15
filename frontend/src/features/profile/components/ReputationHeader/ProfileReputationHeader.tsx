import React from 'react';
import styles from './ReputationHeader.module.css';

interface ProfileReputationHeaderProps {
    score?: number;
}

export const ProfileReputationHeader: React.FC<ProfileReputationHeaderProps> = ({ score }) => {
    if (!score) return null;
    return (
        <div className={styles.reputationCard}>
            <span className={styles.repLabel}>Global Reputation</span>
            <span className={styles.repValue}>{score.toFixed(2)}</span>
        </div>
    );
};