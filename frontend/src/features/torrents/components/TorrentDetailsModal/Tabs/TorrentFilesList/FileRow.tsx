import React, { memo } from 'react';
import { FileOutlined } from '@ant-design/icons';
import type { TorrentFile } from '@/types/model/models.ts';
import { formatBytes } from '@/utils/DashboardScreenUtils/bytesFormatter.ts';
import styles from './TorrentFilesList.module.css';

interface FileRowProps {
    file: TorrentFile;
}

const FileRowComponent: React.FC<FileRowProps> = ({ file }) => {
    return (
        <div className={styles.row}>
            <div className={styles.filePath}>
                <FileOutlined className={styles.fileIcon} />
                {file.path}
            </div>
            <span className={styles.fileSize}>{formatBytes(file.size_bytes)}</span>
        </div>
    );
};

const areFileRowEqual = (prevProps: FileRowProps, nextProps: FileRowProps): boolean => {
    return prevProps.file.path === nextProps.file.path;
};

export const FileRow = memo(FileRowComponent, areFileRowEqual);