import React from 'react';
import type {TorrentDTO, TorrentFile} from '@/types/model/models.ts';
import { TorrentFilesList } from '@/features/torrents/components/TorrentDetailsModal/Tabs/TorrentFilesList/TorrentFilesList';
import styles from './Tabs.module.css';

interface TabProps {
    torrent: TorrentDTO;
}

export const FilesTab: React.FC<TabProps> = ({ torrent }) => {
    const files: TorrentFile[] = torrent.files ?? [];

    return (
        <div className={styles.tabContainer}>
            <TorrentFilesList files={files} />
        </div>
    );
};