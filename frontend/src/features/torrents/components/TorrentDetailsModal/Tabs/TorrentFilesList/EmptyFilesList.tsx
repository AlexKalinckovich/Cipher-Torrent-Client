import React, { memo } from 'react';
import styles from './TorrentFilesList.module.css';

const EmptyFilesListComponent: React.FC = () => {
    return (
        <div className={styles.emptyState}>
            NO FILES METADATA AVAILABLE
        </div>
    );
};

export const EmptyFilesList = memo(EmptyFilesListComponent);