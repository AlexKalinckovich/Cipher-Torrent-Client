import React, { useMemo, useState, useDeferredValue, useCallback } from 'react';
import { Button, Empty, Spin, message, Input, Upload } from 'antd';
import { ArrowLeftOutlined, SearchOutlined, UploadOutlined } from '@ant-design/icons';
import { useQueryClient } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { useStoreTorrents } from '@/hooks/useStoreTorrents.ts';
import { useAuth } from '@/AuthContext.tsx';
import type { StoreTorrent, TorrentDTO, TorrentIdentity } from '@/types/model/models.ts';
import { useAddTorrent, useCreateTorrent, useSignTorrent, useDeletePublishedTorrent } from '@/hooks/useTorrents.ts';
import { StoreCard } from './components/StoreCard/StoreCard';
import styles from './Storefront.module.css';

interface PreparedTorrent {
    torrent: StoreTorrent;
    searchText: string; // precomputed lowercase searchable text
    isOwn: boolean;
}

export const Storefront: React.FC = () => {
    const navigate = useNavigate();
    const queryClient = useQueryClient();
    const { user } = useAuth();
    const { data, isLoading, error } = useStoreTorrents();
    const { mutate: addTorrent, isPending: isAdding } = useAddTorrent();
    const { mutateAsync: createTorrent, isPending: isCreating } = useCreateTorrent();
    const { mutateAsync: signTorrent, isPending: isSigning } = useSignTorrent();
    const { mutate: deletePublished } = useDeletePublishedTorrent();
    const [downloadingHash, setDownloadingHash] = React.useState<string | null>(null);
    const [deletingHash, setDeletingHash] = React.useState<string | null>(null);

    // Search state: input is urgent, list renders from the deferred value.
    const [query, setQuery] = useState('');
    const deferredQuery = useDeferredValue(query);
    const isStale = query !== deferredQuery;

    // Precompute searchable text ONCE when data loads (cheap per render afterward).
    const prepared = useMemo<PreparedTorrent[]>(() => {
        const items = data ?? [];
        return items.map((torrent) => ({
            torrent,
            searchText: `${torrent.name} ${torrent.creator_public_key}`.toLowerCase(),
            isOwn: user?.public_key === torrent.creator_public_key,
        }));
    }, [data, user?.public_key]);

    // Filtering runs against the DEFERRED query -> interruptible, doesn't block typing.
    const results = useMemo<PreparedTorrent[]>(() => {
        const q = deferredQuery.trim().toLowerCase();
        if (!q) {
            return prepared;
        }
        return prepared.filter((it) => it.searchText.includes(q));
    }, [prepared, deferredQuery]);

    const handleDownload = useCallback((torrent: StoreTorrent): void => {
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
    }, [addTorrent, navigate]);

    // Delete a torrent the current user published: removes the MinIO object and
    // the torrents row (cascades to all users' libraries). Only shown for own torrents.
    const handleDelete = useCallback((torrent: StoreTorrent): void => {
        const identity: TorrentIdentity = {
            info_hash: torrent.info_hash,
            creator_pub_key: torrent.creator_public_key,
        };

        setDeletingHash(torrent.info_hash);
        deletePublished(identity, {
            onSuccess: (): void => {
                void message.success(`Deleted "${torrent.name}" from the store`);
                setDeletingHash(null);
            },
            onError: (error: Error): void => {
                void message.error(`Failed to delete "${torrent.name}": ${error.message}`);
                setDeletingHash(null);
            },
        });
    }, [deletePublished]);

    // Upload the user's own torrent: POST /create (sets the authenticated user as
    // creator), then POST /sign to automatically sign it with the uploader's keys.
    const handleUpload = useCallback(async (file: File): Promise<void> => {
        try {
            const created: TorrentDTO = await createTorrent(file);
            const identity: TorrentIdentity = {
                info_hash: created.info_hash,
                creator_pub_key: created.creator_public_key,
            };
            await signTorrent(identity);
            void message.success(`Uploaded and signed "${created.name}"`);
            void queryClient.invalidateQueries({ queryKey: ['store-torrents'] });
        } catch (err) {
            const errorMessage = err instanceof Error ? err.message : 'Unknown error';
            void message.error(`Failed to upload torrent: ${errorMessage}`);
        }
    }, [createTorrent, signTorrent, queryClient]);

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
                <div className={styles.headerActions}>
                    <Upload
                        accept=".torrent"
                        showUploadList={false}
                        customRequest={async ({ file }) => {
                            await handleUpload(file as File);
                        }}
                    >
                        <Button
                            type="primary"
                            icon={<UploadOutlined />}
                            loading={isCreating || isSigning}
                        >
                            Upload my torrent
                        </Button>
                    </Upload>
                </div>
            </div>

            <div className={styles.searchBar}>
                <Input
                    allowClear
                    size="large"
                    prefix={<SearchOutlined />}
                    placeholder="Search by torrent name or creator public key…"
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                />
                {isStale && <span className={styles.searchHint}>Filtering…</span>}
            </div>

            {prepared.length === 0 ? (
                <div className={styles.emptyState}>
                    <Empty description="No torrents published yet" />
                </div>
            ) : results.length === 0 ? (
                <div className={styles.emptyState}>
                    <Empty description={`No torrents match "${query}"`} />
                </div>
            ) : (
                <div
                    className={styles.grid}
                    style={{ opacity: isStale ? 0.6 : 1, transition: 'opacity 150ms' }}
                >
                    {results.map((item) => (
                        <StoreCard
                            key={`${item.torrent.info_hash}-${item.torrent.creator_public_key}`}
                            torrent={item.torrent}
                            isOwn={item.isOwn}
                            isLoading={
                                (isAdding && downloadingHash === item.torrent.info_hash) ||
                                downloadingHash === item.torrent.info_hash
                            }
                            isDeleting={deletingHash === item.torrent.info_hash}
                            onDownload={() => handleDownload(item.torrent)}
                            onDelete={() => handleDelete(item.torrent)}
                        />
                    ))}
                </div>
            )}
        </div>
    );
};