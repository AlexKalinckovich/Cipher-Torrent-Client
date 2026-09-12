import React, { memo, useCallback } from 'react';
import { DownloadOutlined, UserOutlined } from '@ant-design/icons';
import type { StoreTorrent } from '@/types/model/models.ts';
import { formatBytes } from '@/utils/DashboardScreenUtils/bytesFormatter.ts';
import styles from './StoreCard.module.css';

interface StoreCardProps {
    torrent: StoreTorrent;
    isOwn: boolean;
    onSelect: () => void;
}

const truncateKey = (key: string, chars: number = 16): string => {
    if (key.length <= chars) {
        return key;
    }
    return `${key.slice(0, chars)}...${key.slice(-6)}`;
};

const StoreCardComponent: React.FC<StoreCardProps> = ({ torrent, isOwn, onSelect }) => {
    const handleClick = useCallback((): void => {
        onSelect();
    }, [onSelect]);

    return (
        <div className={styles.card} onClick={handleClick} role="button" tabIndex={0}>
            <div className={styles.iconContainer}>
                <DownloadOutlined className={styles.icon} />
            </div>
            <div className={styles.body}>
                <h3 className={styles.name} title={torrent.name}>
                    {torrent.name}
                </h3>
                <div className={styles.meta}>
                    <span className={styles.size}>{formatBytes(torrent.size_bytes)}</span>
                    <span className={styles.added}>
                        {new Date(torrent.added_at).toLocaleDateString()}
                    </span>
                </div>
                <div className={styles.creator} title={torrent.creator_public_key}>
                    <UserOutlined />
                    <span className={styles.creatorKey}>
                        {isOwn ? 'You' : truncateKey(torrent.creator_public_key)}
                    </span>
                    {isOwn && <span className={styles.ownBadge}>yours</span>}
                </div>
            </div>
        </div>
    );
};

export const StoreCard = memo(StoreCardComponent);