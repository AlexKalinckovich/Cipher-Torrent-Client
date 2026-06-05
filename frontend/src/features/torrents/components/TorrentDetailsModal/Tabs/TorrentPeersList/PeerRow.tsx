import React, { memo } from 'react';
import type { Peer } from '@/types/model/models.ts';
import { TorrentSpeed } from '@/features/torrents/components/TorrentSpeed/TorrentSpeed';
import { PeerFlags } from './PeerFlags';
import styles from './TorrentPeersList.module.css';

interface PeerRowProps {
    peer: Peer;
}

const getReputationClass = (score: number | undefined): string => {
    const safeScore: number = score ?? 0;
    if (safeScore >= 0.5) {
        return `${styles.repScore} ${styles.repGood}`;
    }
    return `${styles.repScore} ${styles.repBad}`;
};

const formatTime = (isoString: string | undefined): string => {
    if (!isoString) {
        return 'UNKNOWN';
    }
    return new Date(isoString).toLocaleTimeString();
};

const PeerRowComponent: React.FC<PeerRowProps> = ({ peer }) => {
    return (
        <div className={styles.row}>
            <div className={styles.endpoint}>
                <span>{peer.ip}</span>
                <span className={styles.port}>:{peer.port}</span>
            </div>
            <span className={styles.clientName}>{peer.client_name}</span>
            <TorrentSpeed rx={peer.download_speed_bps} tx={peer.upload_speed_bps} />
            <div className={styles.datesCol}>
                <div className={styles.dateRow}>
                    <span className={styles.dateLabel}>CON</span>
                    <span className={styles.dateValue}>{formatTime(peer.connected_since)}</span>
                </div>
                <div className={styles.dateRow}>
                    <span className={styles.dateLabel}>LST</span>
                    <span className={styles.dateValue}>{formatTime(peer.last_seen)}</span>
                </div>
            </div>
            <PeerFlags flags={peer.flags} />
            <span className={getReputationClass(peer.reputation_score)}>
                {peer.reputation_score?.toFixed(2) ?? '0.00'}
            </span>
        </div>
    );
};

const arePeerRowPropsEqual = (prevProps: PeerRowProps, nextProps: PeerRowProps): boolean => {
    return prevProps.peer.peer_id === nextProps.peer.peer_id;
};

export const PeerRow = memo(PeerRowComponent, arePeerRowPropsEqual);