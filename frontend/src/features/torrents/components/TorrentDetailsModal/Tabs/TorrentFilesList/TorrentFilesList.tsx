import React, { memo } from 'react';
import type { TorrentFile } from '@/types/model/models.ts';
import { FileRow } from './FileRow';
import { EmptyFilesList } from './EmptyFilesList';
import styles from './TorrentFilesList.module.css';

interface TorrentFilesListProps {
    files: TorrentFile[];
}

const TorrentFilesListComponent: React.FC<TorrentFilesListProps> = ({ files }) => {
    if (files.length === 0) {
        return <EmptyFilesList />;
    }

    return (
        <div className={styles.listContainer}>
            {files.map((file) => (
                <FileRow key={file.path} file={file} />
            ))}
        </div>
    );
};

export const TorrentFilesList = memo(TorrentFilesListComponent);