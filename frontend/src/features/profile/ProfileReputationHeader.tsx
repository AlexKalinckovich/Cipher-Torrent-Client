import React from 'react';
import styles from './Profile.module.css';

interface ProfileReputationHeaderProps {
  score?: number;
}

export const ProfileReputationHeader: React.FC<ProfileReputationHeaderProps> = ({ score }) => {
  if (!score) return null;
  return (
    <div className={styles.reputationFloat}>
      <span className={styles.repLabel}>Global Reputation</span>
      <span className={styles.repValue}>{score.toFixed(2)}</span>
    </div>
  );
};