import React, { memo } from 'react';
import styles from './TorrentPeersList.module.css';

const EmptyPeersListComponent: React.FC = () => {
    return (
        <div className={styles.emptyState}>
            <span className={styles.glitchText}>NO ACTIVE PEERS FOUND</span>
        </div>
    );
};

export const EmptyPeersList = memo(EmptyPeersListComponent);