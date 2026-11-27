import React from 'react';
import type { TorrentDTO, Peer } from '@/types/model/models.ts';
import { TorrentPeersList } from '@/features/torrents/components/TorrentDetailsModal/Tabs/TorrentPeersList/TorrentPeerList';
import { randomIp } from '@/utils/random.ts';
import styles from './Tabs.module.css';

interface TabProps {
    torrent: TorrentDTO;
}

const MOCK_PEERS: Peer[] = [
    {
        peer_id: 'peer_1',
        ip: randomIp(),
        port: 51413,
        client_name: 'qBittorrent/4.5.2',
        download_speed_bps: 150000,
        upload_speed_bps: 45000,
    },
    {
        peer_id: 'peer_2',
        ip: randomIp(),
        port: 6881,
        client_name: 'Transmission/3.0',
        download_speed_bps: 320000,
        upload_speed_bps: 12000,
    }
];

export const PeersTab: React.FC<TabProps> = ({ torrent }) => {
    console.log(torrent.storage_path)
    return (
        <div className={styles.tabContainer}>
            <TorrentPeersList peers={MOCK_PEERS} />
        </div>
    );
};