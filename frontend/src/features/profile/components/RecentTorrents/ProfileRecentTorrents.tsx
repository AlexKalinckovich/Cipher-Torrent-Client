import React, { useCallback, useMemo, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { RightOutlined } from '@ant-design/icons';
import { message, Skeleton } from 'antd';
import { useTorrents, useAddTorrent } from '@/hooks/useTorrents.ts';
import type {TorrentDTO, TorrentIdentity} from '@/types/model/models.ts';
import { AddTorrentModal } from '@/features/torrents/components/AddTorrentModal/AddTorrentModal';
import { RecentTorrentsHeader } from './RecentTorrentsHeader';
import { RecentTorrentsList } from './RecentTorrentsList';
import styles from './RecentTorrents.module.css';

const getSortedTorrents = (data: TorrentDTO[] | undefined): TorrentDTO[] => {
    if (!data) {
        return [];
    }
    return [...data].sort((a: TorrentDTO, b: TorrentDTO): number => {
        return new Date(b.added_at).getTime() - new Date(a.added_at).getTime();
    });
};

export const ProfileRecentTorrents: React.FC = () => {
    const navigate = useNavigate();
    const { data, isLoading } = useTorrents();
    const { mutate: addTorrent } = useAddTorrent();
    const [isAddModalOpen, setIsAddModalOpen] = useState<boolean>(false);

    const recentTorrents = useMemo((): TorrentDTO[] => {
        const sorted: TorrentDTO[] = getSortedTorrents(data);
        return sorted.slice(0, 3);
    }, [data]);

    const handleOpenAdd = useCallback((): void => {
        navigate('/store');
    }, [navigate]);

    const handleCloseAdd = useCallback((): void => {
        setIsAddModalOpen(false);
    }, []);

    const handleSubmitAdd = useCallback((request: TorrentIdentity): void => {
        addTorrent(request, {
            onSuccess: (): void => {
                void message.success('NEW TORRENT SUBMITTED TO NETWORK QUEUE');
            },
            onError: (error: Error): void => {
                void message.error(`FAILED TO ADD TORRENT: ${error.message}`);
            }
        });
    }, [addTorrent]);

    const handleRemove = useCallback((e: React.MouseEvent, name: string): void => {
        e.stopPropagation();
        void message.warning(`Removing ${name}...`);
    }, []);

    const handleSelect = useCallback((torrent: TorrentDTO): void => {
        navigate('/inspector', { state: { torrent } });
    }, [navigate]);

    const handleViewAll = useCallback((): void => {
        navigate('/dashboard');
    }, [navigate]);

    if (isLoading) {
        return <Skeleton active paragraph={{ rows: 3 }} />;
    }

    return (
        <div className={styles.recentContainer}>
            <RecentTorrentsHeader onAddClick={handleOpenAdd} />
            <RecentTorrentsList
                torrents={recentTorrents}
                onSelect={handleSelect}
                onRemove={handleRemove}
            />
            <button type="button" className={`${styles.actionBtn} ${styles.viewAllBtn}`} onClick={handleViewAll}>
                View All in Dashboard <RightOutlined />
            </button>
            <AddTorrentModal
                isOpen={isAddModalOpen}
                onClose={handleCloseAdd}
                onSubmit={handleSubmitAdd}
            />
        </div>
    );
};