import { useEffect, useCallback } from 'react';
import { useQueryClient } from '@tanstack/react-query';
import { wsClient } from '../api/webSocketClient.ts';
import type { GlobalEvent } from '../api/webSocketClient.ts';
import type { Torrent, TorrentProgressEvent } from '../types/model/models.ts';

const mapTorrentProgress = (t: Torrent, event: TorrentProgressEvent): Torrent => {
    if (t.info_hash === event.info_hash) {
        return {
            ...t,
            progress: event.progress,
            download_speed_bps: event.download_speed,
            upload_speed_bps: event.upload_speed
        };
    }
    return t;
};

const updateCacheArray = (old: Torrent[] | undefined, event: TorrentProgressEvent): Torrent[] => {
    if (!old) {
        return [];
    }
    return old.map((t: Torrent): Torrent => mapTorrentProgress(t, event));
};

export const useTorrentRealtime = (): void => {
    const queryClient = useQueryClient();

    const handleEvent = useCallback((event: GlobalEvent): void => {
        if (event.event_type === 'torrent:progress') {
            queryClient.setQueryData<Torrent[]>(['torrents'], (old: Torrent[] | undefined): Torrent[] => {
                return updateCacheArray(old, event);
            });
        }
    }, [queryClient]);

    useEffect((): (() => void) => {
        wsClient.subscribe(handleEvent);
        return (): void => {
            wsClient.unsubscribe(handleEvent);
        };
    }, [handleEvent]);
};