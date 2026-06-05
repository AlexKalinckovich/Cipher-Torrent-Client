import React, { useCallback, memo } from 'react';
import { DeleteOutlined } from '@ant-design/icons';
import { formatBytes } from '@/features/profile/utils/profileUtils.ts';
import type { Torrent } from '@/types/model/models.ts';
import styles from './RecentTorrents.module.css';

interface RecentTorrentItemProps {
    torrent: Torrent;
    onSelect: (torrent: Torrent) => void;
    onRemove: (e: React.MouseEvent, name: string) => void;
}

const RecentTorrentItemComponent: React.FC<RecentTorrentItemProps> = ({ torrent, onSelect, onRemove }) => {
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

const areItemsEqual = (prevProps: RecentTorrentItemProps, nextProps: RecentTorrentItemProps): boolean => {
    return prevProps.torrent.info_hash === nextProps.torrent.info_hash;
};

export const RecentTorrentsItem = memo(RecentTorrentItemComponent, areItemsEqual);