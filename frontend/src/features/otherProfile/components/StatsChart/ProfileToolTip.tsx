import React from 'react';
import type { TooltipContentProps } from 'recharts';
import styles from './StatsChart.module.css';

interface TooltipData {
    name: string;
    value: number;
}

export const renderCustomTooltip = ({
    active,
    payload,
}: TooltipContentProps): React.ReactElement | null => {
    if (active && payload && payload.length > 0) {
        const data = payload[0].payload as TooltipData;
        return (
            <div className={styles.chartTooltip}>
                <div className={styles.tooltipName}>{data.name}</div>
                <div className={styles.tooltipValue}>{data.value.toLocaleString()}</div>
            </div>
        );
    }
    return null;
};