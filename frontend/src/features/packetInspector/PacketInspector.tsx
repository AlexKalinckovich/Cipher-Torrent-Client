import React, { useState, useMemo, useCallback } from 'react';
import { usePacketStream } from '@/hooks/usePacketStream.ts';
import type { PacketLog } from '@/types/model/models.ts';
import { InspectorToolbar } from './components/InspectorToolbar/InspectorToolbar';
import { Terminal } from './components/Terminal/Terminal';
import { DetailsPanel } from './components/DetailsPanel/DetailsPanel';
import { filterPackets } from './utils/packetFilterUtils';
import styles from './PacketInspector.module.css';

const checkAndTogglePause = (isPaused: boolean, togglePause: () => void): void => {
    if (!isPaused) {
        togglePause();
    }
};

export const PacketInspector: React.FC = () => {
    const { packets, isPaused, togglePause, clearPackets } = usePacketStream();
    const [searchText, setSearchText] = useState<string>('');
    const [selectedPacket, setSelectedPacket] = useState<PacketLog | null>(null);

    const filteredPackets: PacketLog[] = useMemo((): PacketLog[] => {
        return filterPackets(packets, searchText);
    }, [packets, searchText]);

    const handleSelectPacket = useCallback((packet: PacketLog): void => {
        setSelectedPacket(packet);
        checkAndTogglePause(isPaused, togglePause);
    }, [isPaused, togglePause]);

    return (
        <div className={styles.container}>
            <InspectorToolbar
                isPaused={isPaused}
                onTogglePause={togglePause}
                onClear={clearPackets}
                searchText={searchText}
                onSearchChange={setSearchText}
            />
            <div className={styles.mainArea}>
                <Terminal
                    packets={filteredPackets}
                    selectedId={selectedPacket?.id ?? null}
                    onSelectPacket={handleSelectPacket}
                />
                <DetailsPanel packet={selectedPacket} />
            </div>
        </div>
    );
};