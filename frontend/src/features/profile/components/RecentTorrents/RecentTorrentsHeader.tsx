import React, { memo } from 'react';
import { PlusOutlined } from '@ant-design/icons';
import styles from './RecentTorrents.module.css';

interface RecentTorrentsHeaderProps {
    onAddClick: () => void;
}

const RecentTorrentsHeaderComponent: React.FC<RecentTorrentsHeaderProps> = ({ onAddClick }) => {
    return (
        <div className={styles.sectionHeader}>
            <span className={styles.sectionTitle}>Recent Activity</span>
            <div className={styles.headerActions}>
                <button type="button" className={styles.actionBtn} onClick={onAddClick}>
                    <PlusOutlined /> Add New
                </button>
            </div>
        </div>
    );
};

export const RecentTorrentsHeader = memo(RecentTorrentsHeaderComponent);