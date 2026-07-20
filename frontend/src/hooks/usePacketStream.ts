import { useState, useEffect, useCallback, useRef } from 'react';
import type { PacketLog } from '@/types/model/models.ts';
import { packetStreamService, type StreamEvent } from '@/api/packetStreamService.ts';

const MAX_BUFFER_SIZE = 100;

interface UsePacketStreamResult {
    packets: PacketLog[];
    progress: { download_speed_bps: number; upload_speed_bps: number; progress: number } | null;
    isPaused: boolean;
    togglePause: () => void;
    clearPackets: () => void;
}

const enforceBufferLimit = (list: PacketLog[]): PacketLog[] => {
    return list.length > MAX_BUFFER_SIZE ? list.slice(0, MAX_BUFFER_SIZE) : list;
};

export const usePacketStream = (infoHash: string | null): UsePacketStreamResult => {
    const [packets, setPackets] = useState<PacketLog[]>([]);
    const [progress, setProgress] = useState<UsePacketStreamResult['progress']>(null);
    const [isPaused, setIsPaused] = useState<boolean>(false);

    const isPausedRef = useRef(isPaused);
    useEffect(() => {
        isPausedRef.current = isPaused;
    }, [isPaused]);

    const handleNewEvent = useCallback((event: StreamEvent): void => {
        if (isPausedRef.current) return;

        if (event.event_type === 'packet') {
            setPackets((prev: PacketLog[]) : PacketLog[] => {
                if (prev.some((p: PacketLog) : boolean => p.id === event.packet.id)) {
                    return prev;
                }
                const newList: PacketLog[] = [event.packet, ...prev];
                return enforceBufferLimit(newList);
            });
        } else if (event.event_type === 'torrent:progress') {
            console.log("progress_event")
            setProgress({
                download_speed_bps: event.progress.download_speed_bps,
                upload_speed_bps: event.progress.upload_speed_bps,
                progress: event.progress.progress,
            });
        }
    }, []);

    useEffect(() => {
        if (!infoHash) {
            packetStreamService.disconnect();
            return;
        }

        packetStreamService.connect(infoHash, handleNewEvent);

        return () => {
            packetStreamService.disconnect();
        };
    }, [infoHash, handleNewEvent]);

    const togglePause = useCallback(() => setIsPaused((p) => !p), []);
    const clearPackets = useCallback(() => setPackets([]), []);

    return { packets, progress, isPaused, togglePause, clearPackets };
};