import React, { useMemo } from 'react';
import { useNavigate } from 'react-router-dom';
import { PlusOutlined, DeleteOutlined, RightOutlined } from '@ant-design/icons';
import { message, Skeleton } from 'antd';
import { useTorrents } from '../../hooks/useTorrents';
import { formatBytes } from './profileUtils';
import styles from './Profile.module.css';

export const ProfileRecentTorrents: React.FC = () => {
    const navigate = useNavigate();
    const { data, isLoading } = useTorrents();

    const recentTorrents = useMemo(() => {
        if (!data) return [];
        return [...data]
            .sort((a, b) => new Date(b.added_at).getTime() - new Date(a.added_at).getTime())
            .slice(0, 3);
    }, [data]);

    const handleAdd = (): void => {
        void message.info('Opening file selector...');
    };

    const handleRemove = (e: React.MouseEvent, name: string): void => {
        e.stopPropagation();
        void message.warning(`Removing ${name}...`);
    };

    const handleNavigate = (): void => {
        navigate('/dashboard');
    };

    if (isLoading) return <Skeleton active />;

    return (
        <div className={styles.statsBottomHalf}>
            <div className={styles.sectionHeader}>
                <span className={styles.sectionTitle}>Recent Activity</span>
                <div className={styles.headerActions}>
                    <button type="button" className={styles.actionBtn} onClick={handleAdd}>
                        <PlusOutlined /> Add New
                    </button>
                </div>
            </div>

            <div className={styles.recentList}>
                {recentTorrents.map((torrent) => (
                    <div key={torrent.info_hash} className={styles.recentItem} onClick={handleNavigate}>
                        <div className={styles.itemInfo}>
                            <span className={styles.itemName}>{torrent.name}</span>
                            <div className={styles.itemMeta}>
                                <span>{formatBytes(torrent.size_bytes)}</span>
                                <span>{torrent.status}</span>
                            </div>
                        </div>
                        <DeleteOutlined
                            className={styles.removeIcon}
                            onClick={(e) => handleRemove(e, torrent.name)}
                        />
                    </div>
                ))}
            </div>

            <button type="button" className={`${styles.actionBtn} ${styles.viewAllBtn}`} onClick={handleNavigate}>
                View All in Dashboard <RightOutlined />
            </button>
        </div>
    );
};