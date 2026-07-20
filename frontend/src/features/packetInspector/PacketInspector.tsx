import React, { useState, useMemo, useCallback } from 'react';
import { useSearchParams } from 'react-router-dom';
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
    const [searchParams, setSearchParams] = useSearchParams();
    const infoHash = searchParams.get('infoHash');

    const { packets, isPaused, togglePause, clearPackets } = usePacketStream(infoHash);

    const [searchText, setSearchText] = useState<string>('');
    const [selectedPacket, setSelectedPacket] = useState<PacketLog | null>(null);

    const filteredPackets: PacketLog[] = useMemo((): PacketLog[] => {
        return filterPackets(packets, searchText);
    }, [packets, searchText]);

    const handleSelectPacket = useCallback((packet: PacketLog): void => {
        setSelectedPacket(packet);
        checkAndTogglePause(isPaused, togglePause);
    }, [isPaused, togglePause]);

    const handleInfoHashChange = useCallback((e: React.ChangeEvent<HTMLInputElement>): void => {
        setSearchParams({ infoHash: e.target.value });
    }, [setSearchParams]);

    return (
        <div className={styles.container}>
            <InspectorToolbar
                isPaused={isPaused}
                onTogglePause={togglePause}
                onClear={clearPackets}
                searchText={searchText}
                onSearchChange={setSearchText}
            />
            {!infoHash && (
                <div className={styles.connectPrompt}>
                    <span>Enter Torrent Info Hash to connect: </span>
                    <input
                        type="text"
                        placeholder="e.g. dd8255ecdc7ca55fb0bbf81323d87062db1f6d1c"
                        onChange={handleInfoHashChange}
                        className={styles.hashInput}
                    />
                </div>
            )}
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