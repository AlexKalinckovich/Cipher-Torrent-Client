import React, { memo, useRef, useCallback } from 'react';
import { UploadOutlined } from '@ant-design/icons';
import styles from '../AddTorrentModal.module.css';

interface FileTabProps {
    fileName: string | null;
    onFileSelect: (file: File) => void;
}

const FileTabComponent: React.FC<FileTabProps> = ({ fileName, onFileSelect }) => {
    const inputRef = useRef<HTMLInputElement>(null);

    const handleClick = useCallback((): void => {
        inputRef.current?.click();
    }, []);

    const handleFileChange = useCallback((e: React.ChangeEvent<HTMLInputElement>): void => {
        const files: FileList | null = e.target.files;
        const hasFiles: boolean = files !== null && files.length > 0;

        if (hasFiles) {
            onFileSelect((files as FileList)[0]);
        }
    }, [onFileSelect]);

    return (
        <div className={styles.inputGroup}>
            <span className={styles.label}>TORRENT FILE</span>
            <input
                type="file"
                accept=".torrent"
                className={styles.hiddenInput}
                ref={inputRef}
                onChange={handleFileChange}
            />
            <button type="button" className={styles.fileInputContainer} onClick={handleClick}>
                <span className={styles.fileInputText}>
                    <UploadOutlined style={{ marginRight: 8 }} />
                    {fileName ?? 'CLICK TO SELECT .TORRENT FILE'}
                </span>
            </button>
        </div>
    );
};

export const FileTab = memo(FileTabComponent);