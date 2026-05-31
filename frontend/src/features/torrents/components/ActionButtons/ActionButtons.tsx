import React from 'react';
import { Space } from 'antd';
import { PlayCircleOutlined, PauseCircleOutlined, CloseCircleOutlined } from '@ant-design/icons';
import type { ActionButtonsProps } from '@/features/torrents/types/dashboardTypes';

export const ActionButtons: React.FC<ActionButtonsProps> = ({ onPlay,
                                                                onPause,
                                                                onRemove,
                                                                playClassName,
                                                                pauseClassName,
                                                                removeClassName } : ActionButtonsProps) => {
    return (
        <Space>
            <PlayCircleOutlined className={playClassName} onClick={onPlay} />
            <PauseCircleOutlined className={pauseClassName} onClick={onPause} />
            <CloseCircleOutlined className={removeClassName} onClick={onRemove} />
        </Space>
    );
};
