import React, { useMemo, useState, useCallback } from 'react';
import { Input, message, Space, Table, Typography } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import {
    CloseCircleOutlined,
    DownloadOutlined,
    PauseCircleOutlined,
    PlayCircleOutlined,
    SearchOutlined,
    UploadOutlined,
} from '@ant-design/icons';
import type { Torrent } from '../../types/model/models.ts';
import { useTorrents } from '../../hooks/useTorrents.ts';
import styles from './Dashboard.module.css';
import { formatBytes } from "../../utils/DashboardScreenUtils/bytesFormatter.ts";
import { getFilteredTorrents } from "../../utils/DashboardScreenUtils/dashboardDataFilterUtils.ts";

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

interface FilterButtonProps {
    label: string;
    isActive: boolean;
    onClick: () => void;
    baseClassName: string;
    activeClassName: string;
}

interface StatsSectionProps {
    stats: DashboardStats;
}

interface ErrorStateProps {
    error: Error;
}

interface FilterBarProps {
    currentFilter: string | null;
    searchText: string;
    onFilterChange: (status: string | null) => void;
    onSearchChange: (value: string) => void;
    baseClassName: string;
    activeClassName: string;
}

interface TorrentTableProps {
    data: Torrent[];
    loading: boolean;
    columns: ColumnsType<Torrent>;
}

interface ProgressBarProps {
    progress: number;
}

interface ActionButtonsProps {
    onPlay: () => void;
    onPause: () => void;
    onRemove: () => void;
    playClassName: string;
    pauseClassName: string;
    removeClassName: string;
}

const sumSize = (sum: number, t: Torrent): number => {
    return sum + t.size_bytes;
};

const countStatus = (torrents: Torrent[], status: Torrent['status']): number => {
    return torrents.filter((item: Torrent) => item.status === status).length;
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
    const configMap: Record<Torrent['status'], BadgeConfig> = {
        downloading: { icon: <DownloadOutlined />, text: 'Downloading', className: styles.badgeDownloading },
        seeding: { icon: <UploadOutlined />, text: 'Seeding', className: styles.badgeSeeding },
        paused: { icon: <PauseCircleOutlined />, text: 'Paused', className: styles.badgePaused },
        error: { icon: <CloseCircleOutlined />, text: 'Error', className: styles.badgeError },
    };
    return configMap[status];
};

const StatusBadge: React.FC<{ status: Torrent['status'] }> = ({ status }) => {
    const config: BadgeConfig = getBadgeConfig(status);
    const className: string = `${styles.statBadge} ${config.className}`;
    return (
        <span className={className}>
      {config.icon} {config.text}
    </span>
    );
};

const ProgressBar: React.FC<ProgressBarProps> = ({ progress }) => {
    const safeProgress: number = progress ?? 0;
    const style: React.CSSProperties = { width: `${safeProgress}%` };
    return (
        <div className={styles.progressWrapper}>
            <div className={styles.progressBarBg}>
                <div className={styles.progressBarFill} style={style} />
            </div>
            <span className={styles.progressText}>{safeProgress}%</span>
        </div>
    );
};

const ActionButtons: React.FC<ActionButtonsProps> = ({ onPlay, onPause, onRemove, playClassName, pauseClassName, removeClassName }) => {
    return (
        <Space>
            <PlayCircleOutlined className={playClassName} onClick={onPlay} />
            <PauseCircleOutlined className={pauseClassName} onClick={onPause} />
            <CloseCircleOutlined className={removeClassName} onClick={onRemove} />
        </Space>
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

const COLUMNS: ColumnsType<Torrent> = [
    { title: 'Name', dataIndex: 'name', key: 'name', ellipsis: true, width: '30%' },
    { title: 'Status', dataIndex: 'status', key: 'status', width: '12%', render: (_: unknown, record: Torrent) => <StatusBadge status={record.status} /> },
    { title: 'Progress', dataIndex: 'progress', key: 'progress', width: '15%', render: (progress: number) => <ProgressBar progress={progress} /> },
    { title: 'Size', dataIndex: 'size_bytes', key: 'size_bytes', width: '15%', render: formatBytes },
    { title: 'Peers', dataIndex: 'peers_count', key: 'peers_count', width: '10%' },
    { title: 'Actions', key: 'actions', width: '10%', render: () => (
            <ActionButtons
                onPlay={handlePlayAction}
                onPause={handlePauseAction}
                onRemove={handleRemoveAction}
                playClassName={styles.actionIconGreen}
                pauseClassName={styles.actionIconYellow}
                removeClassName={styles.actionIconRed}
            />
        ) },
];

const StatCard: React.FC<{ title: string; value: string | number; sub: React.ReactElement }> = ({ title, value, sub }) => {
    return (
        <div className={styles.statCard}>
            <div className={styles.statTitle}>{title}</div>
            <div className={styles.statValue}>{value}</div>
            <div className={styles.statSub}>{sub}</div>
        </div>
    );
};

const StatsSection: React.FC<StatsSectionProps> = ({ stats }) => {
    return (
        <div className={styles.statsGrid}>
            <StatCard title="Total" value={stats.total} sub={<span>Torrents</span>} />
            <StatCard title="Downloading" value={stats.downloading} sub={<><DownloadOutlined /> Active</>} />
            <StatCard title="Seeding" value={stats.seeding} sub={<><UploadOutlined /> Sharing</>} />
            <StatCard title="Total Size" value={formatBytes(stats.totalSize)} sub={<span>Used space</span>} />
        </div>
    );
};

const FilterButton: React.FC<FilterButtonProps> = ({ label, isActive, onClick, baseClassName, activeClassName }) => {
    let className: string = baseClassName;
    if (isActive) {
        className = activeClassName;
    }
    return (
        <button type="button" className={className} onClick={onClick}>
            {label}
        </button>
    );
};

const FilterBar: React.FC<FilterBarProps> = ({ currentFilter, searchText, onFilterChange, onSearchChange, baseClassName, activeClassName }) => {
    const handleFilterAll = useCallback(() => onFilterChange(null), [onFilterChange]);
    const handleFilterDownloading = useCallback(() => onFilterChange('downloading'), [onFilterChange]);
    const handleFilterSeeding = useCallback(() => onFilterChange('seeding'), [onFilterChange]);
    const handleFilterPaused = useCallback(() => onFilterChange('paused'), [onFilterChange]);
    const handleFilterError = useCallback(() => onFilterChange('error'), [onFilterChange]);

    const handleSearchInput = useCallback((e: React.ChangeEvent<HTMLInputElement>) => {
        onSearchChange(e.target.value);
    }, [onSearchChange]);

    return (
        <div className={styles.filtersBar}>
            <div className={styles.filterButtons}>
                <FilterButton label="All" isActive={currentFilter === null} onClick={handleFilterAll} baseClassName={baseClassName} activeClassName={activeClassName} />
                <FilterButton label="Downloading" isActive={currentFilter === 'downloading'} onClick={handleFilterDownloading} baseClassName={baseClassName} activeClassName={activeClassName} />
                <FilterButton label="Seeding" isActive={currentFilter === 'seeding'} onClick={handleFilterSeeding} baseClassName={baseClassName} activeClassName={activeClassName} />
                <FilterButton label="Paused" isActive={currentFilter === 'paused'} onClick={handleFilterPaused} baseClassName={baseClassName} activeClassName={activeClassName} />
                <FilterButton label="Error" isActive={currentFilter === 'error'} onClick={handleFilterError} baseClassName={baseClassName} activeClassName={activeClassName} />
            </div>
            <Input
                placeholder="Search torrents..."
                prefix={<SearchOutlined style={{ color: '#8b5cf6' }} />}
                value={searchText}
                onChange={handleSearchInput}
                className={styles.searchInput}
            />
        </div>
    );
};

const TorrentTable: React.FC<TorrentTableProps> = ({ data, loading, columns }) => {
    const paginationConfig = { pageSize: 10, showTotal: (total: number) => `Total ${total} items` };
    return (
        <div className={styles.tableWrapper}>
            <Table<Torrent>
                dataSource={data}
                columns={columns}
                rowKey="info_hash"
                loading={loading}
                pagination={paginationConfig}
            />
        </div>
    );
};

const ErrorState: React.FC<ErrorStateProps> = ({ error }) => {
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

export const Dashboard: React.FC = () => {
    const { data, isLoading, error } = useTorrents();
    const [searchText, setSearchText] = useState<string>('');
    const [statusFilter, setStatusFilter] = useState<string | null>(null);

    const safeData: Torrent[] = useMemo(() => data ?? [], [data]);

    const filteredTorrents: Torrent[] = useMemo(() => {
        return getFilteredTorrents(safeData, statusFilter, searchText);
    }, [safeData, statusFilter, searchText]);

    const stats: DashboardStats = useMemo(() => {
        return calculateStats(safeData);
    }, [safeData]);

    const handleFilterChange = useCallback((status: string | null): void => {
        setStatusFilter(status);
        if (status) {
            void message.info(`Filtered by: ${status}`, 1);
        } else {
            void message.info('Filter cleared', 1);
        }
    }, []);

    const handleSearchChange = useCallback((value: string): void => {
        setSearchText(value);
    }, []);

    if (error) {
        return <ErrorState error={error} />;
    }

    return (
        <div className={styles.dashboardContainer}>
            <div className={styles.backgroundBlob1} />
            <div className={styles.backgroundBlob2} />
            <div className={styles.backgroundBlob3} />
            <StatsSection stats={stats} />
            <FilterBar
                currentFilter={statusFilter}
                searchText={searchText}
                onFilterChange={handleFilterChange}
                onSearchChange={handleSearchChange}
                baseClassName={styles.filterBtn}
                activeClassName={`${styles.filterBtn} ${styles.filterBtnActive}`}
            />
            <TorrentTable data={filteredTorrents} loading={isLoading} columns={COLUMNS} />
        </div>
    );
};