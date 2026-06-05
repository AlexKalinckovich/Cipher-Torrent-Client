import React from 'react';
import type { Torrent } from '@/types/model/models.ts';
import { useTorrentFilesMock } from '@/hooks/useTorrentFilesMock';
import { TorrentFilesList } from '@/features/torrents/components/TorrentDetailsModal/Tabs/TorrentFilesList/TorrentFilesList';
import styles from './Tabs.module.css';

interface TabProps {
    torrent: Torrent;
}

export const FilesTab: React.FC<TabProps> = ({ torrent }) => {
    const files = useTorrentFilesMock(torrent);

    return (
        <div className={styles.tabContainer}>
            <TorrentFilesList files={files} />
        </div>
    );
};