import React, { memo } from 'react';
import type { PeerFlags as IPeerFlags } from '@/types/model/models.ts';
import styles from './TorrentPeersList.module.css';

interface PeerFlagsProps {
    flags?: IPeerFlags;
}

const getFlagClass = (isActive: boolean | undefined): string => {
    if (isActive) {
        return `${styles.flagBadge} ${styles.flagActive}`;
    }
    return `${styles.flagBadge} ${styles.flagInactive}`;
};

const PeerFlagsComponent: React.FC<PeerFlagsProps> = ({ flags }) => {
    return (
        <div className={styles.flagsContainer}>
            <span className={getFlagClass(flags?.encryption)}>ENC</span>
            <span className={getFlagClass(flags?.utp)}>uTP</span>
            <span className={getFlagClass(flags?.dht)}>DHT</span>
            <span className={getFlagClass(flags?.pex)}>PEX</span>
        </div>
    );
};

export const PeerFlags = memo(PeerFlagsComponent);