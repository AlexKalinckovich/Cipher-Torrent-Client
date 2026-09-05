import React from 'react';
import { PlayCircleOutlined, PauseCircleOutlined, DeleteOutlined, InfoCircleOutlined, DownloadOutlined } from '@ant-design/icons';
import styles from './ActionButtons.module.css';

interface ActionButtonsProps {
    onPlay: () => void;
    onPause: () => void;
    onRemove: () => void;
    onInfo: () => void;
    onDownload: () => void;
    playClassName: string;
    pauseClassName: string;
    removeClassName: string;
    infoClassName: string;
    downloadClassName: string;
}

export const ActionButtons: React.FC<ActionButtonsProps> = ({
                                                                onPlay, onPause, onRemove, onInfo, onDownload,
                                                                playClassName, pauseClassName, removeClassName, infoClassName, downloadClassName
                                                            }) => {
    return (
        <div className={styles.actionContainer}>
            <button type="button" className={playClassName} onClick={onPlay}>
                <PlayCircleOutlined />
            </button>
            <button type="button" className={pauseClassName} onClick={onPause}>
                <PauseCircleOutlined />
            </button>
            <button type="button" className={infoClassName} onClick={onInfo}>
                <InfoCircleOutlined />
            </button>
            <button type="button" className={removeClassName} onClick={onRemove}>
                <DeleteOutlined />
            </button>
            <button type="button" className={downloadClassName} onClick={onDownload}>
                <DownloadOutlined />
            </button>
        </div>
    );
};
