import { useCallback, useEffect, useRef, useState } from 'react';
import { torrentClientService, type TorrentProgressSnapshot } from '../api/torrentClientService';

export interface ClientTorrentState {
    snapshot: TorrentProgressSnapshot | null;
    isDownloading: boolean;
}

/**
 * React state hook around the singleton WebTorrent client.
 * Holds per-infoHash progress/speed/peer snapshots in memory so the
 * dashboard and other components can render live client-side state.
 */
export const useTorrentState = (infoHash?: string | null) => {
    const [snapshot, setSnapshot] = useState<TorrentProgressSnapshot | null>(null);
    const [isDownloading, setIsDownloading] = useState<boolean>(false);
    const keyRef = useRef<string | null>(infoHash ?? null);

    useEffect(() => {
        keyRef.current = infoHash ?? null;
    }, [infoHash]);

    useEffect(() => {
        const key = infoHash ?? null;
        if (!key) {
            setSnapshot(null);
            return;
        }

        // Hydrate immediately from any in-memory torrent.
        const existing = torrentClientService.snapshot(key);
        if (existing) {
            setSnapshot(existing);
            setIsDownloading(!existing.done);
        }

        torrentClientService.subscribe(key, {
            onProgress: (s: TorrentProgressSnapshot) => {
                setSnapshot(s);
                setIsDownloading(!s.done);
            },
            onDone: (s: TorrentProgressSnapshot) => {
                setSnapshot(s);
                setIsDownloading(false);
            },
            onError: () => {
                setIsDownloading(false);
            },
        });

        return () => {
            torrentClientService.unsubscribe(key);
        };
    }, [infoHash]);

    const startDownload = useCallback(async (torrentId: string | Uint8Array | File | Blob) => {
        setIsDownloading(true);
        const torrent = await torrentClientService.add(torrentId, {
            onProgress: (s: TorrentProgressSnapshot) => {
                setSnapshot(s);
                setIsDownloading(!s.done);
            },
            onDone: (s: TorrentProgressSnapshot) => {
                setSnapshot(s);
                setIsDownloading(false);
            },
        });
        keyRef.current = torrent.infoHash;
        return torrent;
    }, []);

    const pause = useCallback(async () => {
        if (keyRef.current) {
            const t = await torrentClientService.get(keyRef.current);
            t?.pause?.();
            setIsDownloading(false);
        }
    }, []);

    const resume = useCallback(async () => {
        if (keyRef.current) {
            const t = await torrentClientService.get(keyRef.current);
            t?.resume?.();
            setIsDownloading(true);
        }
    }, []);

    return { snapshot, isDownloading, startDownload, pause, resume };
};