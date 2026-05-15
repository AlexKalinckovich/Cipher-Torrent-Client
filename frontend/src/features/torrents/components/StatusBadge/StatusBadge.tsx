import React from 'react';
import { DownloadOutlined, UploadOutlined, PauseCircleOutlined, CloseCircleOutlined } from '@ant-design/icons';
import type { Torrent } from '@/types/model/models.ts';
import type { BadgeConfig } from '@/features/torrents/types/dashboardTypes';
import styles from './StatusBadge.module.css';

const getBadgeConfig = (status: Torrent['status']): BadgeConfig => {
    const configMap: Record<Torrent['status'], BadgeConfig> = {
        downloading: { icon: <DownloadOutlined />, text: 'Downloading', className: styles.badgeDownloading },
        seeding: { icon: <UploadOutlined />, text: 'Seeding', className: styles.badgeSeeding },
        paused: { icon: <PauseCircleOutlined />, text: 'Paused', className: styles.badgePaused },
        error: { icon: <CloseCircleOutlined />, text: 'Error', className: styles.badgeError },
    };
    return configMap[status];
};

export const StatusBadge: React.FC<{ status: Torrent['status'] }> = ({ status }) => {
    const config: BadgeConfig = getBadgeConfig(status);
    const className: string = `${styles.statBadge} ${config.className}`;
    return (
        <span className={className}>
            {config.icon} {config.text}
        </span>
    );
};
