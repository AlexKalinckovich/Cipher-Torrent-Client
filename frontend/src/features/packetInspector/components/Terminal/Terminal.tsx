import React from 'react';
import type { PacketLog } from '@/types/model/models.ts';
import type { TerminalProps } from '../../types/packetInspectorTypes';
import { PacketRow } from '../PacketRow/PacketRow';
import styles from './Terminal.module.css';

export const Terminal: React.FC<TerminalProps> = ({ packets, selectedId, onSelectPacket }) => {
    return (
        <div className={styles.terminal}>
            {packets.map((packet: PacketLog): React.ReactNode => (
                <PacketRow
                    key={packet.id}
                    packet={packet}
                    isSelected={packet.id === selectedId}
                    onSelect={onSelectPacket}
                />
            ))}
        </div>
    );
};