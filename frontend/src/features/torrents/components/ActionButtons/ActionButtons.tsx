import React from 'react';
import { PlayCircleOutlined, PauseCircleOutlined, DeleteOutlined, InfoCircleOutlined } from '@ant-design/icons';
import styles from './ActionButtons.module.css';

interface ActionButtonsProps {
    onPlay: () => void;
    onPause: () => void;
    onRemove: () => void;
    onInfo: () => void;
    playClassName: string;
    pauseClassName: string;
    removeClassName: string;
    infoClassName: string;
}

export const ActionButtons: React.FC<ActionButtonsProps> = ({
                                                                onPlay, onPause, onRemove, onInfo,
                                                                playClassName, pauseClassName, removeClassName, infoClassName
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
        </div>
    );
};