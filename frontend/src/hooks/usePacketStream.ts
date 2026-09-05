import { useState, useEffect, useCallback, useRef } from 'react';
import type { PacketLog } from '@/types/model/models.ts';
import { torrentClientService } from '@/api/torrentClientService.ts';

const MAX_BUFFER_SIZE = 100;

export interface UsePacketStreamResult {
    packets: PacketLog[];
    progress: { download_speed_bps: number; upload_speed_bps: number; progress: number } | null;
    isPaused: boolean;
    togglePause: () => void;
    clearPackets: () => void;
    isConnected: boolean;
}

const enforceBufferLimit = (list: PacketLog[]): PacketLog[] => {
    return list.length > MAX_BUFFER_SIZE ? list.slice(0, MAX_BUFFER_SIZE) : list;
};

/**
 * Streams packets and progress from the client-side WebTorrent torrent,
 * instead of the backend WebSocket.
 *
 * When an infoHash is provided and a matching torrent is active in the
 * client-side WebTorrent service, we attach a packet listener to its peer
 * wires (bittorrent-protocol events) and mirror its progress snapshot.
 */
export const usePacketStream = (infoHash: string | null): UsePacketStreamResult => {
    const [packets, setPackets] = useState<PacketLog[]>([]);
    const [progress, setProgress] = useState<UsePacketStreamResult['progress']>(null);
    const [isPaused, setIsPaused] = useState<boolean>(false);
    const [isConnected, setIsConnected] = useState<boolean>(false);

    const isPausedRef = useRef(isPaused);
    useEffect(() => {
        isPausedRef.current = isPaused;
    }, [isPaused]);

    const handleNewPacket = useCallback((packet: PacketLog): void => {
        if (isPausedRef.current) return;
        setPackets((prev: PacketLog[]) : PacketLog[] => {
            if (prev.some((p: PacketLog) : boolean => p.id === packet.id)) {
                return prev;
            }
            return enforceBufferLimit([packet, ...prev]);
        });
    }, []);

    useEffect(() => {
        if (!infoHash) {
            setIsConnected(false);
            setProgress(null);
            return;
        }

        let cleanupPackets: (() => void) | undefined;

        // Subscribe to live progress from the WebTorrent snapshot.
        const unsubscribeProgress = torrentClientService.subscribeProgress(infoHash, (snapshot) => {
            setProgress({
                download_speed_bps: snapshot.downloadSpeed,
                upload_speed_bps: snapshot.uploadSpeed,
                progress: snapshot.progress,
            });
        });

        // Attach the client-side packet listener.
        torrentClientService.attachPacketListener(infoHash, handleNewPacket).then((cleanup) => {
            cleanupPackets = cleanup;
            setIsConnected(true);
        });

        return () => {
            unsubscribeProgress();
            cleanupPackets?.();
            setIsConnected(false);
            setProgress(null);
        };
    }, [infoHash, handleNewPacket]);

    const togglePause = useCallback(() => setIsPaused((p) => !p), []);
    const clearPackets = useCallback(() => setPackets([]), []);

    return { packets, progress, isPaused, togglePause, clearPackets, isConnected };
};