import React, { memo } from 'react';
import { Button, Tooltip } from 'antd';
import { DeleteOutlined, DownloadOutlined, UserOutlined } from '@ant-design/icons';
import type { StoreTorrent } from '@/types/model/models.ts';
import styles from './StoreCard.module.css';

interface StoreCardProps {
    torrent: StoreTorrent;
    isOwn: boolean;
    isLoading: boolean;
    isDeleting: boolean;
    onDownload: () => void;
    onDelete: () => void;
}

const truncateKey = (key: string, len = 16): string => {
    if (key.length <= len) {
        return key;
    }
    return `${key.slice(0, len / 2)}…${key.slice(-len / 2)}`;
};

export const StoreCard: React.FC<StoreCardProps> = memo(function StoreCard({
    torrent,
    isOwn,
    isLoading,
    isDeleting,
    onDownload,
    onDelete,
}) {
    return (
        <div className={styles.card}>
            <div className={styles.cardHeader}>
                <span className={styles.cardName} title={torrent.name}>
                    {torrent.name}
                </span>
                {isOwn && <span className={styles.ownBadge}>Yours</span>}
            </div>

            <Tooltip title={torrent.creator_public_key}>
                <div className={styles.creator}>
                    <UserOutlined />
                    <span>{truncateKey(torrent.creator_public_key)}</span>
                </div>
            </Tooltip>

            <div className={styles.meta}>
                <span>{formatBytes(torrent.size_bytes)}</span>
                <span>{new Date(torrent.added_at).toLocaleDateString()}</span>
            </div>

            {isOwn ? (
                <div className={styles.buttonRow}>
                    <Button
                        type="primary"
                        icon={<DownloadOutlined />}
                        loading={isLoading}
                        block
                        onClick={onDownload}
                    >
                        Download
                    </Button>
                    <Button
                        danger
                        icon={<DeleteOutlined />}
                        loading={isDeleting}
                        block
                        onClick={onDelete}
                    >
                        Delete
                    </Button>
                </div>
            ) : (
                <Button
                    type="primary"
                    icon={<DownloadOutlined />}
                    loading={isLoading}
                    block
                    onClick={onDownload}
                >
                    Download
                </Button>
            )}
        </div>
    );
});

const formatBytes = (bytes: number): string => {
    if (bytes === 0) {
        return '0 B';
    }
    const units = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = Math.floor(Math.log(bytes) / Math.log(1024));
    return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${units[i]}`;
};
