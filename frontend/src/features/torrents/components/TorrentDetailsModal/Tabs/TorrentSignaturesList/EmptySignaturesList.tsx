import React, { memo } from 'react';
import styles from './TorrentSignaturesList.module.css';

const EmptySignaturesListComponent: React.FC = () => {
    return (
        <div className={styles.emptyState}>
            NO CRYPTOGRAPHIC SIGNATURES DETECTED
        </div>
    );
};

export const EmptySignaturesList = memo(EmptySignaturesListComponent);