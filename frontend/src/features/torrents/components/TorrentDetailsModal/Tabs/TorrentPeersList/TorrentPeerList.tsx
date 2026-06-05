import React, { memo } from 'react';
import type { Peer } from '@/types/model/models.ts';
import { PeerRow } from './PeerRow.tsx';
import { EmptyPeersList } from './EmptyPeersList.tsx';
import styles from './TorrentPeersList.module.css';

interface TorrentPeersListProps {
    peers: Peer[];
}

const TorrentPeersListComponent: React.FC<TorrentPeersListProps> = ({ peers }) => {
    if (peers.length === 0) {
        return <EmptyPeersList />;
    }

    return (
        <div className={styles.listContainer}>
            {peers.map((peer: Peer): React.ReactNode => (
                <PeerRow key={peer.peer_id} peer={peer} />
            ))}
        </div>
    );
};

export const TorrentPeersList = memo(TorrentPeersListComponent);