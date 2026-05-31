import React, { useMemo, useState, useCallback } from 'react';
import { useNavigate } from 'react-router-dom';
import { message } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import type { Torrent } from '@/types/model/models.ts';
import { useTorrents } from '@/hooks/useTorrents.ts';
import { formatBytes } from '@/utils/DashboardScreenUtils/bytesFormatter.ts';
import { getFilteredTorrents } from '@/utils/DashboardScreenUtils/dashboardDataFilterUtils.ts';
import type { DashboardStats } from './types/dashboardTypes';

import { StatusBadge } from './components/StatusBadge/StatusBadge';
import { ProgressBar } from './components/ProgressBar/ProgressBar';
import { ActionButtons } from './components/ActionButtons/ActionButtons';
import { StatsSection } from './components/StatsSection/StatsSection';
import { FilterBar } from './components/FilterBar/FilterBar';
import { TorrentTable } from './components/TorrentTable/TorrentTable';
import { ErrorState } from './components/ErrorState/ErrorState';
import { ProfileLink } from './components/ProfileLink/ProfileLink';

import filterStyles from './components/FilterBar/FilterBar.module.css';
import actionStyles from './components/ActionButtons/ActionButtons.module.css';
import commonStyles from './common.module.css';

const sumSize = (sum: number, t: Torrent): number => {
    return sum + t.size_bytes;
};

const countStatus = (torrents: Torrent[], status: Torrent['status']): number => {
    return torrents.filter((item: Torrent): boolean => item.status === status).length;
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
    { title: 'Status', dataIndex: 'status', key: 'status', width: '12%', render: (_: unknown, record: Torrent): React.ReactNode => <StatusBadge status={record.status} /> },
    { title: 'Progress', dataIndex: 'progress', key: 'progress', width: '15%', render: (progress: number): React.ReactNode => <ProgressBar progress={progress} /> },
    { title: 'Size', dataIndex: 'size_bytes', key: 'size_bytes', width: '15%', render: formatBytes },
    { title: 'Peers', dataIndex: 'peers_count', key: 'peers_count', width: '10%' },
    { title: 'Actions', key: 'actions', width: '10%', render: (): React.ReactNode => (
            <ActionButtons
                onPlay={handlePlayAction}
                onPause={handlePauseAction}
                onRemove={handleRemoveAction}
                playClassName={actionStyles.actionIconGreen}
                pauseClassName={actionStyles.actionIconYellow}
                removeClassName={actionStyles.actionIconRed}
            />
        ) },
];

export const Dashboard: React.FC = () => {
    const { data, isLoading, error } = useTorrents();
    const navigate = useNavigate();
    const [searchText, setSearchText] = useState<string>('');
    const [statusFilter, setStatusFilter] = useState<string | null>(null);

    const safeData: Torrent[] = useMemo((): Torrent[] => data ?? [], [data]);

    const filteredTorrents: Torrent[] = useMemo((): Torrent[] => {
        return getFilteredTorrents(safeData, statusFilter, searchText);
    }, [safeData, statusFilter, searchText]);

    const stats: DashboardStats = useMemo((): DashboardStats => {
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

    const handleRowClick = useCallback((torrent: Torrent): void => {
        navigate('/inspector', { state: { torrent } });
    }, [navigate]);

    if (error) {
        return <ErrorState error={error} />;
    }

    return (
        <div className={commonStyles.dashboardContainer}>
            <div className={commonStyles.backgroundBlob1} />
            <div className={commonStyles.backgroundBlob2} />
            <div className={commonStyles.backgroundBlob3} />

            <ProfileLink />

            <StatsSection stats={stats} />

            <FilterBar
                currentFilter={statusFilter}
                searchText={searchText}
                onFilterChange={handleFilterChange}
                onSearchChange={handleSearchChange}
                baseClassName={filterStyles.filterBtn}
                activeClassName={`${filterStyles.filterBtn} ${filterStyles.filterBtnActive}`}
            />

            <TorrentTable
                data={filteredTorrents}
                loading={isLoading}
                columns={COLUMNS}
                onRowClick={handleRowClick}
            />
        </div>
    );
};