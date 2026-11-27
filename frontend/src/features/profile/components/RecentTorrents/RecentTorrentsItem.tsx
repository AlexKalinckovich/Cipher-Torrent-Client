import React, { useCallback, memo } from 'react';
import { DeleteOutlined } from '@ant-design/icons';
import { formatBytes } from '@/features/profile/utils/profileUtils.ts';
import type { TorrentDTO } from '@/types/model/models.ts';
import styles from './RecentTorrents.module.css';

interface RecentTorrentItemProps {
    torrent: TorrentDTO;
    onSelect: (torrent: TorrentDTO) => void;
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
        <div className={styles.recentItem}>
            <button type="button" className={styles.itemInfo} onClick={handleItemClick}>
                <span className={styles.itemName}>{torrent.name}</span>
                <div className={styles.itemMeta}>
                    <span>{formatBytes(torrent.size_bytes)}</span>
                    <span>{torrent.status}</span>
                </div>
            </button>
            <button
                type="button"
                className={styles.removeIcon}
                onClick={handleRemoveClick}
                aria-label={`Remove ${torrent.name}`}
            >
                <DeleteOutlined />
            </button>
        </div>
    );
};

const areItemsEqual = (prevProps: RecentTorrentItemProps, nextProps: RecentTorrentItemProps): boolean => {
    return prevProps.torrent.info_hash === nextProps.torrent.info_hash;
};

export const RecentTorrentsItem = memo(RecentTorrentItemComponent, areItemsEqual);