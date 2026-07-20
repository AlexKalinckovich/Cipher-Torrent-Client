import React, {useMemo, useState, useCallback, useEffect} from 'react';
import { useNavigate } from 'react-router-dom';
import { message } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import type {ProgressLog, TorrentDTO, TorrentIdentity} from '@/types/model/models.ts';
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
import { TorrentDetailsModal } from './components/TorrentDetailsModal/TorrentDetailsModal';
import filterStyles from './components/FilterBar/FilterBar.module.css';
import actionStyles from './components/ActionButtons/ActionButtons.module.css';
import commonStyles from './common.module.css';
import {
    usePauseTorrent,
    useResumeTorrent,
    useDeleteTorrent
} from '@/hooks/useTorrents.ts';
import {dashboardProgressService} from "@/api/dashboardProgressService.ts";
const sumSize = (sum: number, t: TorrentDTO): number => {
    return sum + t.size_bytes;
};

const countStatus = (torrents: TorrentDTO[], status: TorrentDTO['status']): number => {
    return torrents.filter((item: TorrentDTO): boolean => item.status === status).length;
};

const calculateStats = (torrents: TorrentDTO[]): DashboardStats => {
    return {
        total: torrents.length,
        downloading: countStatus(torrents, 'downloading'),
        seeding: countStatus(torrents, 'seeding'),
        paused: countStatus(torrents, 'paused'),
        totalSize: torrents.reduce(sumSize, 0),
    };
};



const notifyFilter = (status: string | null): void => {
    if (status) {
        void message.info(`Filtered by: ${status}`, 1);
        return;
    }
    void message.info('Filter cleared', 1);
};

export const Dashboard: React.FC = () => {
    const { data, isLoading, error } = useTorrents();
    const navigate = useNavigate();

    const { mutate: pauseTorrent } = usePauseTorrent();
    const { mutate: resumeTorrent } = useResumeTorrent();
    const { mutate: deleteTorrent } = useDeleteTorrent();

    const [torrents, setTorrents] = useState<TorrentDTO[]>([]);

    useEffect(() => {
        if (data) {
            setTorrents(data);
        }
    }, [data]);

    useEffect(() => {
        const unsubscribe = dashboardProgressService.subscribe((progress: ProgressLog) => {
            setTorrents((prev) =>
                prev.map((t) => {
                    if (
                        t.info_hash === progress.info_hash &&
                        t.creator_public_key === progress.creator_public_key
                    ) {
                        return { ...t, progress: progress.progress };
                    }
                    return t;
                })
            );
        });

        return () => {
            unsubscribe();
        };
    }, []);

    const handlePlayAction = useCallback((record: TorrentDTO): void => {
        const identity: TorrentIdentity = {
            info_hash: record.info_hash,
            creator_pub_key: record.creator_public_key,
        };

        resumeTorrent(identity, {
            onSuccess: (): void => {
                void message.success(`Resumed: ${record.name}`);
            },
            onError: (error: Error): void => {
                void message.error(`Failed to resume ${record.name}: ${error.message}`);
            }
        });
    }, [resumeTorrent]);

    // 2. Pause Action
    const handlePauseAction = useCallback((record: TorrentDTO): void => {
        const identity: TorrentIdentity = {
            info_hash: record.info_hash,
            creator_pub_key: record.creator_public_key,
        };

        pauseTorrent(identity, {
            onSuccess: (): void => {
                void message.success(`Paused: ${record.name}`);
            },
            onError: (error: Error): void => {
                void message.error(`Failed to pause ${record.name}: ${error.message}`);
            }
        });
    }, [pauseTorrent]);


    const handleRemoveAction = useCallback((record: TorrentDTO): void => {
        const identity: TorrentIdentity = {
            info_hash: record.info_hash,
            creator_pub_key: record.creator_public_key,
        };

        deleteTorrent(identity, {
            onSuccess: (): void => {
                void message.success(`Removed: ${record.name}`);
            },
            onError: (error: Error): void => {
                void message.error(`Failed to remove ${record.name}: ${error.message}`);
            }
        });
    }, [deleteTorrent]);



    const [searchText, setSearchText] = useState<string>('');
    const [statusFilter, setStatusFilter] = useState<string | null>(null);
    const [detailsTorrent, setDetailsTorrent] = useState<TorrentDTO | null>(null);

    const filteredTorrents: TorrentDTO[] = useMemo((): TorrentDTO[] => {
        return getFilteredTorrents(torrents, statusFilter, searchText);
    }, [torrents, statusFilter, searchText]);

    const stats: DashboardStats = useMemo((): DashboardStats => {
        return calculateStats(torrents);
    }, [torrents]);

    const handleFilterChange = useCallback((status: string | null): void => {
        setStatusFilter(status);
        notifyFilter(status);
    }, []);

    const handleSearchChange = useCallback((value: string): void => {
        setSearchText(value);
    }, []);

    const handleRowClick = useCallback((torrent: TorrentDTO): void => {
        navigate(`/inspector?infoHash=${torrent.info_hash}`, { state: { torrent } });
    }, [navigate]);

    const handleOpenDetails = useCallback((torrent: TorrentDTO): void => {
        setDetailsTorrent(torrent);
    }, []);

    const handleCloseDetails = useCallback((): void => {
        setDetailsTorrent(null);
    }, []);



    const COLUMNS: ColumnsType<TorrentDTO> = useMemo((): ColumnsType<TorrentDTO> => [
        { title: 'Name', dataIndex: 'name', key: 'name', ellipsis: true, width: '30%' },
        {
            title: 'Status',
            dataIndex: 'status',
            key: 'status',
            width: '10%',
            render: (_: unknown, record: TorrentDTO): React.ReactNode => <StatusBadge status={record.status} />
        },
        {
            title: 'Progress',
            dataIndex: 'progress',
            key: 'progress',
            width: '15%',
            render: (progress: number): React.ReactNode => <ProgressBar progress={progress} />
        },
        {
            title: 'Size',
            dataIndex: 'size_bytes',
            key: 'size_bytes',
            width: '10%',
            render: formatBytes
        },
        {
            title: 'Actions',
            key: 'actions',
            width: '15%',
            render: (_: unknown, record: TorrentDTO): React.ReactNode => (
                <ActionButtons
                    onPlay={(): void => handlePlayAction(record)}
                    onPause={(): void => handlePauseAction(record)}
                    onRemove={(): void => handleRemoveAction(record)}
                    onInfo={(): void => handleOpenDetails(record)}
                    playClassName={actionStyles.actionIconGreen}
                    pauseClassName={actionStyles.actionIconYellow}
                    removeClassName={actionStyles.actionIconRed}
                    infoClassName={actionStyles.actionIconBlue}
                />
            )
        },
    ], [handleOpenDetails]);

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
            <TorrentDetailsModal
                torrent={detailsTorrent}
                onClose={handleCloseDetails}
            />
        </div>
    );
};