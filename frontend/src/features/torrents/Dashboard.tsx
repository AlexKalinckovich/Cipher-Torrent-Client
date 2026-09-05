import React, {useMemo, useState, useCallback, useEffect} from 'react';
import { useNavigate } from 'react-router-dom';
import { message } from 'antd';
import type { ColumnsType } from 'antd/es/table';
import type {TorrentDTO, TorrentIdentity} from '@/types/model/models.ts';
import { useTorrents, useDeleteTorrent } from '@/hooks/useTorrents.ts';
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
import { useTorrentState } from '@/hooks/useTorrentState.ts';
import { torrentService } from '@/api/torrentService.ts';
import { torrentClientService } from '@/api/torrentClientService.ts';
import type { TorrentProgressSnapshot } from '@/api/torrentClientService.ts';
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

    const { mutate: deleteTorrent } = useDeleteTorrent();
    const { startDownload } = useTorrentState();

    const [torrents, setTorrents] = useState<TorrentDTO[]>([]);
    // Client-side progress, keyed by info_hash, overrides the catalog status.
    const [clientProgress, setClientProgress] = useState<Record<string, TorrentProgressSnapshot>>({});

    useEffect(() => {
        if (data) {
            setTorrents(data);
        }
    }, [data]);

    // Track live client-side progress for any active WebTorrent torrents.
    useEffect(() => {
        const cleanupFns: Array<() => void> = [];

        torrents.forEach((t) => {
            const snapshot = torrentClientService.snapshot(t.info_hash);
            if (snapshot) {
                setClientProgress((prev) => ({ ...prev, [t.info_hash]: snapshot }));
            }
            const cleanup = torrentClientService.subscribeProgress(t.info_hash, (s: TorrentProgressSnapshot) => {
                setClientProgress((prev) => ({ ...prev, [s.infoHash]: s }));
            });
            cleanupFns.push(cleanup);
        });

        return () => {
            cleanupFns.forEach((cleanup) => cleanup());
        };
    }, [torrents]);

    const handlePlayAction = useCallback((record: TorrentDTO): void => {
        void torrentClientService.get(record.info_hash)
            .then((t) => {
                if (t) {
                    t.resume?.();
                    void message.success(`Resumed: ${record.name}`);
                } else {
                    void message.info(`Not downloaded in this browser yet: ${record.name}`);
                }
            })
            .catch((err: Error) => {
                void message.error(`Failed to resume ${record.name}: ${err.message}`);
            });
    }, []);

    // 2. Pause Action
    const handlePauseAction = useCallback((record: TorrentDTO): void => {
        void torrentClientService.get(record.info_hash)
            .then((t) => {
                if (t) {
                    t.pause?.();
                    void message.success(`Paused: ${record.name}`);
                } else {
                    void message.info(`Not downloaded in this browser yet: ${record.name}`);
                }
            })
            .catch((err: Error) => {
                void message.error(`Failed to pause ${record.name}: ${err.message}`);
            });
    }, []);


    const handleRemoveAction = useCallback((record: TorrentDTO): void => {
        const identity: TorrentIdentity = {
            info_hash: record.info_hash,
            creator_pub_key: record.creator_public_key,
        };

        deleteTorrent(identity, {
            onSuccess: (): void => {
                void torrentClientService.remove(record.info_hash);
                void message.success(`Removed: ${record.name}`);
            },
            onError: (error: Error): void => {
                void message.error(`Failed to remove ${record.name}: ${error.message}`);
            }
        });
    }, [deleteTorrent]);

    const handleDownloadAction = useCallback((record: TorrentDTO): void => {
        const identity: TorrentIdentity = {
            info_hash: record.info_hash,
            creator_pub_key: record.creator_public_key,
        };

        // Fetch the small .torrent METADATA descriptor from the backend,
        // then download the actual data client-side via WebTorrent.
        torrentService.downloadTorrentFile(identity)
            .then((torrentFile: Blob): Promise<import('webtorrent').Torrent> => startDownload(torrentFile))
            .then((): void => {
                void message.success(`Downloading: ${record.name}`);
            })
            .catch((err: Error): void => {
                void message.error(`Failed to download ${record.name}: ${err.message}`);
            });
    }, [startDownload]);



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
            render: (_: unknown, record: TorrentDTO): React.ReactNode => {
                const live = clientProgress[record.info_hash];
                const status = live && !live.done ? 'downloading' : record.status;
                return <StatusBadge status={status} />;
            }
        },
        {
            title: 'Progress',
            dataIndex: 'progress',
            key: 'progress',
            width: '15%',
            render: (_: unknown, record: TorrentDTO): React.ReactNode => {
                const live = clientProgress[record.info_hash];
                const progress = live ? Math.round(live.progress * 100) : record.progress;
                return <ProgressBar progress={progress} />;
            }
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
                    onDownload={(): void => handleDownloadAction(record)}
                    playClassName={actionStyles.actionIconGreen}
                    pauseClassName={actionStyles.actionIconYellow}
                    removeClassName={actionStyles.actionIconRed}
                    infoClassName={actionStyles.actionIconBlue}
                    downloadClassName={actionStyles.actionIconCyan}
                />
            )
        },
    ], [clientProgress, handleOpenDetails]);

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