import React from 'react';
import type { Torrent } from '@/types/model/models.ts';
import { useTorrentPeersMock } from '@/hooks/useTorrentPeersMock';
import { TorrentPeersList } from '@/features/torrents/components/TorrentDetailsModal/Tabs/TorrentPeersList/TorrentPeerList';
import styles from './Tabs.module.css';

interface TabProps {
    torrent: Torrent;
}

export const PeersTab: React.FC<TabProps> = ({ torrent }) => {
    const peers = useTorrentPeersMock(torrent.info_hash);

    return (
        <div className={styles.tabContainer}>
            <TorrentPeersList peers={peers} />
        </div>
    );
};