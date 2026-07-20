import React from 'react';
import { DownloadOutlined, UploadOutlined, PauseCircleOutlined, ClockCircleOutlined } from '@ant-design/icons';
import type { TorrentDTO } from '@/types/model/models.ts';
import type { BadgeConfig } from '@/features/torrents/types/dashboardTypes';
import styles from './StatusBadge.module.css';

const getBadgeConfig = (status: TorrentDTO['status']): BadgeConfig => {
    const configMap: Record<TorrentDTO['status'], BadgeConfig> = {
        downloading: { icon: <DownloadOutlined />, text: 'Downloading', className: styles.badgeDownloading },
        seeding: { icon: <UploadOutlined />, text: 'Seeding', className: styles.badgeSeeding },
        paused: { icon: <PauseCircleOutlined />, text: 'Paused', className: styles.badgePaused },
        idle: { icon: <ClockCircleOutlined />, text: 'Idle', className: styles.badgeIdle },
    };
    return configMap[status];
};

export const StatusBadge: React.FC<{ status: TorrentDTO['status'] }> = ({ status }) => {
    const config: BadgeConfig = getBadgeConfig(status);
    const className: string = `${styles.statBadge} ${config.className}`;
    return (
        <span className={className}>
            {config.icon} {config.text}
        </span>
    );
};