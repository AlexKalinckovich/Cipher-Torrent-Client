import type { PacketLog } from '@/types/model/models.ts';

export interface InspectorToolbarProps {
    isPaused: boolean;
    onTogglePause: () => void;
    onClear: () => void;
    searchText: string;
    onSearchChange: (value: string) => void;
}

export interface TerminalProps {
    packets: PacketLog[];
    selectedId: number | null;
    onSelectPacket: (packet: PacketLog) => void;
}

export interface PacketRowProps {
    packet: PacketLog;
    isSelected: boolean;
    onSelect: (packet: PacketLog) => void;
}

export interface DetailsPanelProps {
    packet: PacketLog | null;
}

export interface ActiveDetailsPanelProps {
    packet: PacketLog;
}