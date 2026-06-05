import React, { memo } from 'react';
import styles from './DetailsPanel.module.css';

const EmptyDetailsPanelComponent: React.FC = () => {
    return (
        <div className={styles.detailsPanel}>
            <div className={styles.emptyState}>Select a packet to inspect</div>
        </div>
    );
};

export const EmptyDetailsPanel = memo(EmptyDetailsPanelComponent);