import React, { useCallback, useMemo, useState } from 'react';
import {type NavigateFunction, useNavigate} from 'react-router-dom';
import { RightOutlined } from '@ant-design/icons';
import { message, Skeleton } from 'antd';
import { useTorrents } from '@/hooks/useTorrents.ts';
import type { Torrent, TorrentAddRequest } from '@/types/model/models.ts';
import { AddTorrentModal } from '@/features/torrents/components/AddTorrentModal/AddTorrentModal';
import { RecentTorrentsHeader } from './RecentTorrentsHeader';
import { RecentTorrentsList } from './RecentTorrentsList';
import styles from './RecentTorrents.module.css';

const getSortedTorrents = (data: Torrent[] | undefined): Torrent[] => {
    if (!data) {
        return [];
    }
    return [...data].sort((a: Torrent, b: Torrent): number => {
        return new Date(b.added_at).getTime() - new Date(a.added_at).getTime();
    });
};

export const ProfileRecentTorrents: React.FC = () => {
    const navigate: NavigateFunction = useNavigate();
    const { data, isLoading } = useTorrents();
    const [isAddModalOpen, setIsAddModalOpen] = useState<boolean>(false);

    const recentTorrents: Torrent[] = useMemo((): Torrent[] => {
        const sorted: Torrent[] = getSortedTorrents(data);
        return sorted.slice(0, 3);
    }, [data]);

    const handleOpenAdd = useCallback((): void => {
        setIsAddModalOpen(true);
    }, []);

    const handleCloseAdd = useCallback((): void => {
        setIsAddModalOpen(false);
    }, []);

    const handleSubmitAdd = useCallback((request: TorrentAddRequest): void => {
        void message.success('NEW TORRENT SUBMITTED TO NETWORK QUEUE');
        console.log('Payload:', request);
    }, []);

    const handleRemove = useCallback((e: React.MouseEvent, name: string): void => {
        e.stopPropagation();
        void message.warning(`Removing ${name}...`);
    }, []);

    const handleSelect = useCallback((torrent: Torrent): void => {
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