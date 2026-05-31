import React, { useCallback } from 'react';
import { Input } from 'antd';
import { PlayCircleOutlined, PauseCircleOutlined, DeleteOutlined, SearchOutlined } from '@ant-design/icons';
import type { InspectorToolbarProps } from '@/features/packetInspector/types/packetInspectorTypes';
import styles from './InspectorToolbar.module.css';

const getStatusClass = (isPaused: boolean): string => {
    if (isPaused) {
        return `${styles.statusDot} ${styles.statusPaused}`;
    }
    return `${styles.statusDot} ${styles.statusActive}`;
};

export const InspectorToolbar: React.FC<InspectorToolbarProps> = ({
                                                                      isPaused,
                                                                      onTogglePause,
                                                                      onClear,
                                                                      searchText,
                                                                      onSearchChange
                                                                  }) => {
    const handleSearch: (e: React.ChangeEvent<HTMLInputElement>) => void = useCallback((e: React.ChangeEvent<HTMLInputElement>): void => {
        onSearchChange(e.target.value);
    }, [onSearchChange]);

    return (
        <div className={styles.header}>
            <div className={styles.statusWrapper}>
                <div className={getStatusClass(isPaused)} />
                <h1 className={styles.title}>Packet Inspector</h1>
            </div>
            <div className={styles.controls}>
                <Input
                    placeholder="Filter IP or Type..."
                    prefix={<SearchOutlined />}
                    value={searchText}
                    onChange={handleSearch}
                    className={styles.searchInput}
                />
                <button type="button" className={styles.btn} onClick={onTogglePause}>
                    {isPaused ? <><PlayCircleOutlined /> Resume</> : <><PauseCircleOutlined /> Pause</>}
                </button>
                <button type="button" className={styles.btn} onClick={onClear}>
                    <DeleteOutlined /> Clear
                </button>
            </div>
        </div>
    );
};