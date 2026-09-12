import React, { useCallback, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button, Empty, Spin, message, Upload } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import type { UploadProps } from 'antd';
import { useStoreTorrents } from '@/hooks/useStoreTorrents.ts';
import { useAddTorrent } from '@/hooks/useTorrents.ts';
import { useTorrentState } from '@/hooks/useTorrentState.ts';
import { useAuth } from '@/AuthContext.tsx';
import type { StoreTorrent, TorrentIdentity } from '@/types/model/models.ts';
import { torrentService } from '@/api/torrentService.ts';
import { StoreCard } from './components/StoreCard/StoreCard';
import styles from './Storefront.module.css';

export const Storefront: React.FC = () => {
    const navigate = useNavigate();
    const { user } = useAuth();
    const { data, isLoading, error } = useStoreTorrents();
    const { mutate: addTorrent, isPending: isAdding } = useAddTorrent();
    const { startDownload } = useTorrentState();
    const [uploading, setUploading] = useState(false);
    const [selectedHash, setSelectedHash] = useState<string | null>(null);
    const [downloadingHash, setDownloadingHash] = useState<string | null>(null);

    const handleUpload: UploadProps['customRequest'] = useCallback(async (options) => {
        const file = options.file as File;
        setUploading(true);
        try {
            // Currently the "Upload your torrent" flow is a placeholder that
            // navigates to the dashboard for now, or here we could POST the
            // file via torrentService.createTorrent(file).
            void message.info(`Received file: ${file.name} — upload flow to be wired up`);
        } finally {
            setUploading(false);
        }
    }, []);

    const handleAddOwnTorrent = useCallback((): void => {
        // Placeholder: "+" button that will prompt to upload the user's own torrent file.
        void message.info('Choose a .torrent file to publish it to the store');
    }, []);

    const handleSelect = useCallback((torrent: StoreTorrent): void => {
        const identity: TorrentIdentity = {
            info_hash: torrent.info_hash,
            creator_pub_key: torrent.creator_public_key
        };

        setSelectedHash(torrent.info_hash);
        addTorrent(identity, {
            onSuccess: (): void => {
                void message.success(`Added "${torrent.name}" to your torrents`);
                navigate('/dashboard');
            },
            onError: (error: Error): void => {
                void message.error(`Failed to add "${torrent.name}": ${error.message}`);
                setSelectedHash(null);
            }
        });
    }, [addTorrent, navigate]);

    const handleDownload = useCallback((torrent: StoreTorrent): void => {
        const identity: TorrentIdentity = {
            info_hash: torrent.info_hash,
            creator_pub_key: torrent.creator_public_key
        };

        setDownloadingHash(torrent.info_hash);
        torrentService.downloadTorrentFile(identity)
            .then((torrentFile: Blob): Promise<import('webtorrent').Torrent> => startDownload(torrentFile))
            .then((): void => {
                void message.success(`Downloading "${torrent.name}"...`);
                navigate('/dashboard');
            })
            .catch((err: Error): void => {
                void message.error(`Failed to download "${torrent.name}": ${err.message}`);
                setDownloadingHash(null);
            });
    }, [startDownload, navigate]);

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
                <h1 className={styles.title}>Torrent Storefront</h1>
                <Upload
                    accept=".torrent"
                    showUploadList={false}
                    customRequest={handleUpload}
                    disabled={uploading}
                >
                    <Button
                        type="primary"
                        icon={<PlusOutlined />}
                        loading={uploading}
                        onClick={handleAddOwnTorrent}
                    >
                        Upload your torrent
                    </Button>
                </Upload>
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
                                (isAdding && selectedHash === torrent.info_hash) ||
                                downloadingHash === torrent.info_hash
                            }
                            onSelect={() => handleSelect(torrent)}
                            onDownload={() => handleDownload(torrent)}
                        />
                    ))}
                </div>
            )}
        </div>
    );
};