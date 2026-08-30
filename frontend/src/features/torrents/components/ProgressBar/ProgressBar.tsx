import React from 'react';
import type { ProgressBarProps } from '@/features/torrents/types/dashboardTypes';
import styles from './ProgressBar.module.css';

export const ProgressBar: React.FC<ProgressBarProps> = ({ progress }) => {
    const safeProgress: number = progress ?? 0;
    const style: React.CSSProperties = { width: `${safeProgress}%` };
    const formattedProgress = safeProgress.toFixed(2);
    return (
        <div className={styles.progressWrapper}>
            <div className={styles.progressBarBg}>
                <div className={styles.progressBarFill} style={style} />
            </div>
            <span className={styles.progressText}>{formattedProgress}%</span>
        </div>
    );
};
