import React, { useMemo, useState, useDeferredValue, useCallback } from 'react';
import { Button, Empty, Spin, message, Input } from 'antd';
import { ArrowLeftOutlined, SearchOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useStoreTorrents } from '@/hooks/useStoreTorrents.ts';
import { useAuth } from '@/AuthContext.tsx';
import type { StoreTorrent, TorrentIdentity } from '@/types/model/models.ts';
import { useAddTorrent } from '@/hooks/useTorrents.ts';
import { StoreCard } from './components/StoreCard/StoreCard';
import styles from './Storefront.module.css';

interface PreparedTorrent {
    torrent: StoreTorrent;
    searchText: string; // precomputed lowercase searchable text
    isOwn: boolean;
}

export const Storefront: React.FC = () => {
    const navigate = useNavigate();
    const { user } = useAuth();
    const { data, isLoading, error } = useStoreTorrents();
    const { mutate: addTorrent, isPending: isAdding } = useAddTorrent();
    const [downloadingHash, setDownloadingHash] = React.useState<string | null>(null);

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
                            onDownload={() => handleDownload(item.torrent)}
                        />
                    ))}
                </div>
            )}
        </div>
    );
};