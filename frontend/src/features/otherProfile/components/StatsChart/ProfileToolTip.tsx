import React from 'react';
import type { StatItem } from '@/features/profile/types/profileTypes';
import styles from './StatsChart.module.css';

interface ProfileToolTipProps {
    active?: boolean;
    payload?: Array<{
        payload: {
            name: string;
            value: number;
        };
    }>;
}

export const renderCustomTooltip = ({ active, payload }: ProfileToolTipProps): React.ReactElement | null => {
    if (active && payload && payload.length > 0) {
        const data = payload[0].payload;
        return (
            <div className={styles.chartTooltip}>
                <div className={styles.tooltipName}>{data.name}</div>
                <div className={styles.tooltipValue}>{data.value.toLocaleString()}</div>
            </div>
        );
    }
    return null;
};
