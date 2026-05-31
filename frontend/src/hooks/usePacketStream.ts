import { useState, useEffect, useCallback } from 'react';
import type { PacketLog, PacketEvent } from '../types/model/models.ts';
import { packetStreamService } from '@/api/packetStreamService.ts';

const MAX_BUFFER_SIZE = 100;

interface UsePacketStreamResult {
    packets: PacketLog[];
    isPaused: boolean;
    togglePause: () => void;
    clearPackets: () => void;
}

const enforceBufferLimit = (list: PacketLog[]): PacketLog[] => {
    if (list.length > MAX_BUFFER_SIZE) {
        return list.slice(0, MAX_BUFFER_SIZE);
    }
    return list;
};

const appendPacket = (currentList: PacketLog[], newPacket: PacketLog): PacketLog[] => {
    const updatedList: PacketLog[] = [newPacket, ...currentList];
    return enforceBufferLimit(updatedList);
};

const disconnectStream = (): void => {
    packetStreamService.disconnect();
};

export const usePacketStream = (): UsePacketStreamResult => {
    const [packets, setPackets] = useState<PacketLog[]>([]);
    const [isPaused, setIsPaused] = useState<boolean>(false);

    const handleNewPacket = useCallback((event: PacketEvent): void => {
        if (!isPaused) {
            setPackets((prev: PacketLog[]) : PacketLog[] => {
                return appendPacket(prev, event.packet);
            });
        }
    }, [isPaused]);

    useEffect(() => {
        packetStreamService.connect(handleNewPacket);
        return disconnectStream;
    }, [handleNewPacket]);

    const togglePause = useCallback((): void => {
        setIsPaused((prev: boolean) : boolean => {
            return !prev;
        });
    }, []);

    const clearPackets = useCallback((): void => {
        setPackets([]);
    }, []);

    return { packets, isPaused, togglePause, clearPackets };
};