import React, {useMemo, useState} from 'react';
import {Input, message, Space, Table, Typography} from 'antd';
import type {ColumnsType} from 'antd/es/table';
import {
    CloseCircleOutlined,
    DownloadOutlined,
    PauseCircleOutlined,
    PlayCircleOutlined,
    SearchOutlined,
    UploadOutlined
} from '@ant-design/icons';
import type {Torrent} from '../../types/model/models.ts';
import {useTorrents} from '../../hooks/useTorrents.ts';
import styles from './Dashboard.module.css';
import {formatBytes} from "../../utils/DashboardScreenUtils/bytesFormatter.ts";
import {getFilteredTorrents} from "../../utils/DashboardScreenUtils/dashboardDataFilterUtils.ts";

interface DashboardStats {
    total: number;
    downloading: number;
    seeding: number;
    paused: number;
    error: number;
    totalSize: number;
}

interface BadgeConfig {
    icon: React.ReactElement;
    text: string;
    className: string;
}

const sumSize = (sum: number, t: Torrent): number => {
    return sum + t.size_bytes;
};

const countStatus = (torrents: Torrent[], status: Torrent['status']): number => {
    return torrents.filter((t: Torrent) => t.status === status).length;
};

const calculateStats = (torrents: Torrent[]): DashboardStats => {
    return {
        total: torrents.length,
        downloading: countStatus(torrents, 'downloading'),
        seeding: countStatus(torrents, 'seeding'),
        paused: countStatus(torrents, 'paused'),
        error: countStatus(torrents, 'error'),
        totalSize: torrents.reduce(sumSize, 0),
    };
};

const getBadgeConfig = (status: Torrent['status']): BadgeConfig => {
    switch (status) {
        case 'downloading':
            return { icon: <DownloadOutlined />, text: 'Downloading', className: styles.badgeDownloading };
        case 'seeding':
            return { icon: <UploadOutlined />, text: 'Seeding', className: styles.badgeSeeding };
        case 'paused':
            return { icon: <PauseCircleOutlined />, text: 'Paused', className: styles.badgePaused };
        case 'error':
            return { icon: <CloseCircleOutlined />, text: 'Error', className: styles.badgeError };
    }
};

const renderStatusBadge = (status: Torrent['status']): React.ReactElement => {
    const config: BadgeConfig = getBadgeConfig(status);
    return (
        <span className={`${styles.statBadge} ${config.className}`}>
            {config.icon} {config.text}
        </span>
    );
};

const renderProgress = (progress: number): React.ReactElement => {
    const safeProgress = progress ?? 0;
    return (
        <div className={styles.progressWrapper}>
            <div className={styles.progressBarBg}>
                <div className={styles.progressBarFill} style={{ width: `${safeProgress}%` }} />
            </div>
            <span className={styles.progressText}>{safeProgress}%</span>
        </div>
    );
};

const handlePlayAction = (): void => {
    void message.info('Play action');
};

const handlePauseAction = (): void => {
    void message.info('Pause action');
};

const handleRemoveAction = (): void => {
    void message.info('Remove action');
};

const renderActions = (): React.ReactElement => {
    return (
        <Space>
            <PlayCircleOutlined style={{ cursor: 'pointer', color: '#34d399', fontSize: 18 }} onClick={handlePlayAction} />
            <PauseCircleOutlined style={{ cursor: 'pointer', color: '#fbbf24', fontSize: 18 }} onClick={handlePauseAction} />
            <CloseCircleOutlined style={{ cursor: 'pointer', color: '#f87171', fontSize: 18 }} onClick={handleRemoveAction} />
        </Space>
    );
};

const getColumns = (): ColumnsType<Torrent> => {
    return [
        { title: 'Name', dataIndex: 'name', key: 'name', ellipsis: true, width: '30%' },
        { title: 'Status', dataIndex: 'status', key: 'status', width: '12%', render: renderStatusBadge },
        { title: 'Progress', dataIndex: 'progress', key: 'progress', width: '15%', render: renderProgress },
        { title: 'Size', dataIndex: 'size_bytes', key: 'size_bytes', width: '15%', render: formatBytes },
        { title: 'Peers', dataIndex: 'peers_count', key: 'peers_count', width: '10%' },
        { title: 'Actions', key: 'actions', width: '10%', render: renderActions }
    ];
};

const StatsSection: React.FC<{ stats: DashboardStats }> = ({ stats }) => {
    return (
        <div className={styles.statsGrid}>
            <div className={styles.statCard}><div className={styles.statTitle}>Total</div><div className={styles.statValue}>{stats.total}</div><div className={styles.statSub}>Torrents</div></div>
            <div className={styles.statCard}><div className={styles.statTitle}>Downloading</div><div className={styles.statValue}>{stats.downloading}</div><div className={styles.statSub}><DownloadOutlined /> Active</div></div>
            <div className={styles.statCard}><div className={styles.statTitle}>Seeding</div><div className={styles.statValue}>{stats.seeding}</div><div className={styles.statSub}><UploadOutlined /> Sharing</div></div>
            <div className={styles.statCard}><div className={styles.statTitle}>Total Size</div><div className={styles.statValue}>{formatBytes(stats.totalSize)}</div><div className={styles.statSub}>Used space</div></div>
        </div>
    );
};

const renderErrorState = (error: Error): React.ReactElement => {
    return (
        <div className={styles.dashboardContainer}>
            <div className={styles.errorState}>
                <CloseCircleOutlined style={{ fontSize: 64, color: '#f87171', marginBottom: 24 }} />
                <Typography.Title level={3} style={{ color: '#fff' }}>Failed to load torrents</Typography.Title>
                <Typography.Text style={{ color: '#c4b5fd' }}>{error.message}</Typography.Text>
            </div>
        </div>
    );
};

const getFilterButtonClass = (currentFilter: string | null, targetFilter: string | null): string => {
    if (currentFilter === targetFilter) return `${styles.filterBtn} ${styles.filterBtnActive}`;
    return styles.filterBtn;
};

export const Dashboard: React.FC = () => {
    const { data, isLoading, error } = useTorrents();
    const [searchText, setSearchText] = useState<string>('');
    const [statusFilter, setStatusFilter] = useState<string | null>(null);

    const safeData: Torrent[] = useMemo(() => data ?? [], [data]);

    const filteredTorrents = useMemo(() => {
        return getFilteredTorrents(safeData, statusFilter, searchText);
    }, [safeData, searchText, statusFilter]);

    const stats = useMemo(() => {
        return calculateStats(safeData);
    }, [safeData]);

    const handleFilterChange = (status: string | null): void => {
        setStatusFilter(status);
        void message.info(`Filtered by: ${status || 'All'}`, 1);
    };

    const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>): void => {
        setSearchText(e.target.value);
    };

    if (error) {
        return renderErrorState(error);
    }

    return (
        <div className={styles.dashboardContainer}>
            <div className={styles.backgroundBlob1} />
            <div className={styles.backgroundBlob2} />
            <div className={styles.backgroundBlob3} />

            <StatsSection stats={stats} />

            <div className={styles.filtersBar}>
                <div className={styles.filterButtons}>
                    <button type="button" className={getFilterButtonClass(statusFilter, null)} onClick={() => handleFilterChange(null)}>All</button>
                    <button type="button" className={getFilterButtonClass(statusFilter, 'downloading')} onClick={() => handleFilterChange('downloading')}>Downloading</button>
                    <button type="button" className={getFilterButtonClass(statusFilter, 'seeding')} onClick={() => handleFilterChange('seeding')}>Seeding</button>
                    <button type="button" className={getFilterButtonClass(statusFilter, 'paused')} onClick={() => handleFilterChange('paused')}>Paused</button>
                    <button type="button" className={getFilterButtonClass(statusFilter, 'error')} onClick={() => handleFilterChange('error')}>Error</button>
                </div>
                <Input
                    placeholder="Search torrents..."
                    prefix={<SearchOutlined style={{ color: '#8b5cf6' }} />}
                    value={searchText}
                    onChange={handleSearchChange}
                    className={styles.searchInput}
                />
            </div>

            <div className={styles.tableWrapper}>
                <Table<Torrent>
                    dataSource={filteredTorrents}
                    columns={getColumns()}
                    rowKey="info_hash"
                    loading={isLoading}
                    pagination={{ pageSize: 10, showTotal: (total) => `Total ${total} items` }}
                />
            </div>
        </div>
    );
};