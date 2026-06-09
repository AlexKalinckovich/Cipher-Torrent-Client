import React, { memo } from 'react';
import { DownloadOutlined, UploadOutlined } from '@ant-design/icons';
import { formatBytes } from '@/utils/DashboardScreenUtils/bytesFormatter.ts';
import type { StatsSectionProps } from '@/features/torrents/types/dashboardTypes';
import styles from './StatsSection.module.css';

const StatCard: React.FC<{ title: string; value: string | number; sub: React.ReactElement }> = ({ title, value, sub }) => {
    return (
        <div className={styles.statCard}>
            <div className={styles.statTitle}>{title}</div>
            <div className={styles.statValue}>{value}</div>
            <div className={styles.statSub}>{sub}</div>
        </div>
    );
};

const StatsSectionComponent: React.FC<StatsSectionProps> = ({ stats }) => {
    return (
        <div className={styles.statsGrid}>
            <StatCard title="Total" value={stats.total} sub={<span>Torrents</span>} />
            <StatCard title="Downloading" value={stats.downloading} sub={<><DownloadOutlined /> Active</>} />
            <StatCard title="Seeding" value={stats.seeding} sub={<><UploadOutlined /> Sharing</>} />
            <StatCard title="Total Size" value={formatBytes(stats.totalSize)} sub={<span>Used space</span>} />
        </div>
    );
};

const areStatsEqual = (prevProps: StatsSectionProps, nextProps: StatsSectionProps): boolean => {
    return prevProps.stats.total === nextProps.stats.total &&
        prevProps.stats.downloading === nextProps.stats.downloading &&
        prevProps.stats.seeding === nextProps.stats.seeding &&
        prevProps.stats.paused === nextProps.stats.paused &&
        prevProps.stats.error === nextProps.stats.error &&
        prevProps.stats.totalSize === nextProps.stats.totalSize;
};

export const StatsSection = memo(StatsSectionComponent, areStatsEqual);