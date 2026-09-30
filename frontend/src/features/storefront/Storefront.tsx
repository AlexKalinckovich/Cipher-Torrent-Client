import React from 'react';
import { Button, Empty, Spin, message } from 'antd';
import { ArrowLeftOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useStoreTorrents } from '@/hooks/useStoreTorrents.ts';
import { useAuth } from '@/AuthContext.tsx';
import type { StoreTorrent, TorrentIdentity } from '@/types/model/models.ts';
import { useAddTorrent } from '@/hooks/useTorrents.ts';
import { StoreCard } from './components/StoreCard/StoreCard';
import styles from './Storefront.module.css';

export const Storefront: React.FC = () => {
    const navigate = useNavigate();
    const { user } = useAuth();
    const { data, isLoading, error } = useStoreTorrents();
    const { mutate: addTorrent, isPending: isAdding } = useAddTorrent();
    const [downloadingHash, setDownloadingHash] = React.useState<string | null>(null);

    // Download = add the torrent to the current user's library.
    // The backend (Anacrolix engine) downloads the actual content server-side.
    const handleDownload = (torrent: StoreTorrent): void => {
        const identity: TorrentIdentity = {
            info_hash: torrent.info_hash,
            creator_pub_key: torrent.creator_public_key,
        };

        setDownloadingHash(torrent.info_hash);
        addTorrent(identity, {
            onSuccess: (): void => {
                void message.success(`Added "${torrent.name}" to your downloads`);
                navigate('/dashboard');
            },
            onError: (error: Error): void => {
                void message.error(`Failed to download "${torrent.name}": ${error.message}`);
                setDownloadingHash(null);
            },
        });
    };

    if (isLoading) {
        return (
            <div className={styles.pageContainer}>
                <Spin size="large" />
            </div>
        );
    }

    if (error) {
        return (
            <div className={styles.pageContainer}>
                <Empty description={`Failed to load store: ${error.message}`} />
            </div>
        );
    }

    const torrents: StoreTorrent[] = data ?? [];

    return (
        <div className={styles.pageContainer}>
            <div className={styles.backgroundBlob1} />
            <div className={styles.backgroundBlob2} />

            <div className={styles.header}>
                <Button
                    type="text"
                    icon={<ArrowLeftOutlined />}
                    onClick={() => navigate(-1)}
                >
                    Back
                </Button>
                <h1 className={styles.title}>Torrent Storefront</h1>
            </div>

            {torrents.length === 0 ? (
                <div className={styles.emptyState}>
                    <Empty description="No torrents published yet" />
                </div>
            ) : (
                <div className={styles.grid}>
                    {torrents.map((torrent) => (
                        <StoreCard
                            key={`${torrent.info_hash}-${torrent.creator_public_key}`}
                            torrent={torrent}
                            isOwn={user?.public_key === torrent.creator_public_key}
                            isLoading={
                                (isAdding && downloadingHash === torrent.info_hash) ||
                                downloadingHash === torrent.info_hash
                            }
                            onDownload={() => handleDownload(torrent)}
                        />
                    ))}
                </div>
            )}
        </div>
    );
};