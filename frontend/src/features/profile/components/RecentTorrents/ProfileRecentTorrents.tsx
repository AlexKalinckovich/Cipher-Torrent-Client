import React, { useCallback, useMemo } from 'react';
import {type NavigateFunction, useNavigate} from 'react-router-dom';
import { PlusOutlined, DeleteOutlined, RightOutlined } from '@ant-design/icons';
import { message, Skeleton } from 'antd';
import { useTorrents } from '@/hooks/useTorrents.ts';
import { formatBytes } from '@/features/profile/utils/profileUtils';
import type { Torrent } from '@/types/model/Torrent/torrent.ts';
import styles from './RecentTorrents.module.css';

interface RecentTorrentItemProps {
    torrent: Torrent;
    onSelect: (torrent: Torrent) => void;
    onRemove: (e: React.MouseEvent, name: string) => void;
}

const RecentTorrentItem: React.FC<RecentTorrentItemProps> = ({ torrent, onSelect, onRemove }) => {
    const handleItemClick = useCallback((): void => {
        onSelect(torrent);
    }, [torrent, onSelect]);

    const handleRemoveClick = useCallback((e: React.MouseEvent): void => {
        onRemove(e, torrent.name);
    }, [torrent.name, onRemove]);

    return (
        <div className={styles.recentItem} onClick={handleItemClick}>
            <div className={styles.itemInfo}>
                <span className={styles.itemName}>{torrent.name}</span>
                <div className={styles.itemMeta}>
                    <span>{formatBytes(torrent.size_bytes)}</span>
                    <span>{torrent.status}</span>
                </div>
            </div>
            <DeleteOutlined
                className={styles.removeIcon}
                onClick={handleRemoveClick}
            />
        </div>
    );
};

export const ProfileRecentTorrents: React.FC = () => {
    const navigate: NavigateFunction = useNavigate();
    const { data, isLoading } = useTorrents();

    const sortedData: Torrent[] = useMemo((): Torrent[] => {
        if (!data) {
            return [];
        }
        return [...data].sort((a: Torrent, b: Torrent): number => {
            return new Date(b.added_at).getTime() - new Date(a.added_at).getTime();
        });
    }, [data]);

    const recentTorrents: Torrent[] = useMemo((): Torrent[] => {
        return sortedData.slice(0, 3);
    }, [sortedData]);

    const handleAdd = useCallback((): void => {
        void message.info('Opening file selector...');
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
            <div className={styles.sectionHeader}>
                <span className={styles.sectionTitle}>Recent Activity</span>
                <div className={styles.headerActions}>
                    <button type="button" className={styles.actionBtn} onClick={handleAdd}>
                        <PlusOutlined /> Add New
                    </button>
                </div>
            </div>

            <div className={styles.recentList}>
                {recentTorrents.map((torrent: Torrent): React.ReactNode => (
                    <RecentTorrentItem
                        key={torrent.info_hash}
                        torrent={torrent}
                        onSelect={handleSelect}
                        onRemove={handleRemove}
                    />
                ))}
            </div>

            <button type="button" className={`${styles.actionBtn} ${styles.viewAllBtn}`} onClick={handleViewAll}>
                View All in Dashboard <RightOutlined />
            </button>
        </div>
    );
};