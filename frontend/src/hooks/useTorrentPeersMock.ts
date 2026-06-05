import { useMemo } from 'react';
import type { Peer } from '@/types/model/models.ts';

const generateMockPeers = (infoHash: string): Peer[] => {
    const now: number = Date.now();
    return [
        {
            peer_id: `${infoHash.substring(0, 5)}_Client_A`,
            ip: '192.168.1.45',
            port: 45901,
            client_name: 'qBittorrent/4.3.9',
            supports_ut_reputation: true,
            reputation_score: 0.92,
            download_speed_bps: 1250000,
            upload_speed_bps: 45000,
            flags: { encryption: true, utp: true, dht: true, pex: true },
            connected_since: new Date(now - 3600000).toISOString(),
            last_seen: new Date(now - 2000).toISOString()
        },
        {
            peer_id: `${infoHash.substring(0, 5)}_Client_B`,
            ip: '88.213.45.12',
            port: 51413,
            client_name: 'Transmission/3.00',
            supports_ut_reputation: false,
            reputation_score: 0.45,
            download_speed_bps: 0,
            upload_speed_bps: 890000,
            flags: { encryption: true, utp: false, dht: true, pex: false },
            connected_since: new Date(now - 7200000).toISOString(),
            last_seen: new Date(now - 15000).toISOString()
        },
        {
            peer_id: `${infoHash.substring(0, 5)}_Client_C`,
            ip: '104.22.14.99',
            port: 6881,
            client_name: 'Deluge/2.0.3',
            supports_ut_reputation: true,
            reputation_score: 0.78,
            download_speed_bps: 5400000,
            upload_speed_bps: 1200000,
            flags: { encryption: false, utp: true, dht: false, pex: true },
            connected_since: new Date(now - 1800000).toISOString(),
            last_seen: new Date(now - 500).toISOString()
        }
    ];
};

export const useTorrentPeersMock = (infoHash: string): Peer[] => {
    return useMemo((): Peer[] => {
        if (!infoHash) {
            return [];
        }
        return generateMockPeers(infoHash);
    }, [infoHash]);
};