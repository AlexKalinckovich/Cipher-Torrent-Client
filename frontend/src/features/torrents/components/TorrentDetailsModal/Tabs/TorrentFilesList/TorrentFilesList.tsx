import React, { memo } from 'react';
import type { TorrentFile } from '@/types/model/models.ts';
import { FileRow } from './FileRow.tsx';
import { EmptyFilesList } from './EmptyFilesList.tsx';
import styles from './TorrentFilesList.module.css';

interface TorrentFilesListProps {
    files: TorrentFile[];
}

const TorrentFilesListComponent: React.FC<TorrentFilesListProps> = ({ files } : TorrentFilesListProps) => {
    if (files.length === 0) {
        return <EmptyFilesList />;
    }

    return (
        <div className={styles.listContainer}>
            {files.map((file: TorrentFile): React.ReactNode => (
                <FileRow key={file.path} file={file} />
            ))}
        </div>
    );
};

export const TorrentFilesList = memo(TorrentFilesListComponent);