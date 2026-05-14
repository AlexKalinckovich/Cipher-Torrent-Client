import React from 'react';
import { formatBytes } from '@/features/profile/utils/profileUtils';
import type { StatItem } from '@/features/profile/types/profileTypes';
import styles from './StatsChart.module.css';

interface ProfileLegendProps {
    items: StatItem[];
    isByteFormat: boolean;
}

const renderValue = (item: StatItem, isByte: boolean): string => {
    if (isByte) return formatBytes(item.value);
    return item.value.toString();
};

export const ProfileLegend: React.FC<ProfileLegendProps> = ({ items, isByteFormat }) => {
    return (
        <div className={styles.legendSection}>
            {items.map((item) => (
                <div key={item.name} className={styles.legendItem}>
                    <div className={styles.colorDot} style={{ backgroundColor: item.color }} />
                    <span className={styles.legendName}>{item.name}</span>
                    <span className={styles.legendVal}>{renderValue(item, isByteFormat)}</span>
                </div>
            ))}
        </div>
    );
};