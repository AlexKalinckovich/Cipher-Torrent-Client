import React from 'react';
import { formatStatValue } from '@/features/profile/utils/profileUtils';
import type { StatItem } from '@/features/profile/types/profileTypes';
import styles from './StatsChart.module.css';

interface ProfileTooltipPayload {
    payload: StatItem;
}

export interface ProfileTooltipProps {
    active?: boolean;
    payload?: ProfileTooltipPayload[];
}

const extractTooltipData = (payload: ProfileTooltipPayload[]): StatItem | undefined => {
    return payload[0]?.payload;
};

const isValidTooltip = (active?: boolean, payload?: ProfileTooltipPayload[]): boolean => {
    return Boolean(active && payload && payload.length > 0);
};

const buildTooltipContent = (data: StatItem): React.ReactElement => {
    return (
        <div className={styles.chartTooltip}>
            <div className={styles.tooltipName}>{data.name}</div>
            <div className={styles.tooltipValue}>{formatStatValue(data.name, data.value)}</div>
        </div>
    );
};

const handleValidTooltip = (data?: StatItem): React.ReactElement | null => {
    if (!data) {
        return null;
    }
    return buildTooltipContent(data);
};

export const renderCustomTooltip = (props: unknown): React.ReactElement | null => {
    const tooltipProps = props as ProfileTooltipProps;
    const isActive: boolean = isValidTooltip(tooltipProps.active, tooltipProps.payload);

    if (!isActive) {
        return null;
    }

    return handleValidTooltip(extractTooltipData(tooltipProps.payload as ProfileTooltipPayload[]));
};