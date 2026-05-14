import React, { useMemo } from 'react';
import type { UserStats } from '../../../../types/model/models.ts';
import type { StatItem } from '../../types/profileTypes';
import { buildTrafficData, buildActivityData } from '../../utils/profileUtils';
import { ProfileStatsChart } from '../StatsChart/ProfileStatsChart';
import { ProfileLegend } from '../StatsChart/ProfileLegend';
import { ProfileRecentTorrents } from '../RecentTorrents/ProfileRecentTorrents';
import styles from './RightPanel.module.css';

interface ProfileRightPanelProps {
    stats?: UserStats;
}

export const ProfileRightPanel: React.FC<ProfileRightPanelProps> = ({ stats }) => {
    const trafficData: StatItem[] = useMemo(() => buildTrafficData(stats), [stats]);
    const activityData: StatItem[] = useMemo(() => buildActivityData(stats), [stats]);

    return (
        <div className={styles.rightPanel}>
            <div className={styles.statsTopHalf}>
                <ProfileLegend items={trafficData} isByteFormat={true} />
                <ProfileStatsChart data={trafficData} />
                <div className={styles.activityGrid}>
                    {activityData.map(renderActivityCard)}
                </div>
            </div>

            <ProfileRecentTorrents />
        </div>
    );
};

const renderActivityCard = (item: StatItem): React.ReactElement => (
    <div key={item.name} className={styles.activityCard}>
        <div
            className={styles.activityValue}
            style={{ color: item.color, textShadow: `0 0 10px ${item.color}55` }}
        >
            {item.value}
        </div>
        <div className={styles.activityLabel}>{item.name}</div>
    </div>
);