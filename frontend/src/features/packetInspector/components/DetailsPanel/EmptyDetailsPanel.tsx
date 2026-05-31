import React from 'react';
import styles from './DetailsPanel.module.css';

export const EmptyDetailsPanel: React.FC = () => {
    return (
        <div className={styles.detailsPanel}>
            <div className={styles.emptyState}>Select a packet to inspect</div>
        </div>
    );
};