import React, { memo } from 'react';
import type { TorrentDTO } from '@/types/model/models.ts';
import { RecentTorrentsItem } from './RecentTorrentsItem.tsx';
import styles from './RecentTorrents.module.css';

interface RecentTorrentsListProps {
    torrents: TorrentDTO[];
    onSelect: (torrent: TorrentDTO) => void;
    onRemove: (e: React.MouseEvent, name: string) => void;
}

const RecentTorrentsListComponent: React.FC<RecentTorrentsListProps> = ({ torrents, onSelect, onRemove }) => {
    return (
        <div className={styles.recentList}>
            {torrents.map((torrent: TorrentDTO): React.ReactNode => (
                <RecentTorrentsItem
                    key={torrent.info_hash}
                    torrent={torrent}
                    onSelect={onSelect}
                    onRemove={onRemove}
                />
            ))}
        </div>
    );
};

export const RecentTorrentsList = memo(RecentTorrentsListComponent);