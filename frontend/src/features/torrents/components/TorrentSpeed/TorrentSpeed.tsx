import React from 'react';
import { formatSpeed } from '@/utils/DashboardScreenUtils/speedFormatter';
import styles from './TorrentSpeed.module.css';

interface TorrentSpeedProps {
    rx?: number;
    tx?: number;
}

export const TorrentSpeed: React.FC<TorrentSpeedProps> = ({ rx, tx } : TorrentSpeedProps) => {
    return (
        <div className={styles.speedContainer}>
            <div className={styles.speedRow}>
                <span className={styles.rxIcon}>▼</span>
                <span className={styles.speedText}>{formatSpeed(rx)}</span>
            </div>
            <div className={styles.speedRow}>
                <span className={styles.txIcon}>▲</span>
                <span className={styles.speedText}>{formatSpeed(tx)}</span>
            </div>
        </div>
    );
};