import React, { memo } from 'react';
import type { Peer } from '@/types/model/models.ts';
import { TorrentSpeed } from '@/features/torrents/components/TorrentSpeed/TorrentSpeed';
import styles from './TorrentPeersList.module.css';

interface PeerRowProps {
    peer: Peer;
}

const PeerRowComponent: React.FC<PeerRowProps> = ({ peer }) => {
    return (
        <div className={styles.row}>
            <div className={styles.endpoint}>
                <span>{peer.ip}</span>
                <span className={styles.port}>:{peer.port}</span>
            </div>
            <span className={styles.clientName}>{peer.client_name}</span>
            <TorrentSpeed rx={peer.download_speed_bps} tx={peer.upload_speed_bps} />
        </div>
    );
};

const arePeerRowPropsEqual = (prevProps: PeerRowProps, nextProps: PeerRowProps): boolean => {
    return prevProps.peer.peer_id === nextProps.peer.peer_id;
};

export const PeerRow = memo(PeerRowComponent, arePeerRowPropsEqual);