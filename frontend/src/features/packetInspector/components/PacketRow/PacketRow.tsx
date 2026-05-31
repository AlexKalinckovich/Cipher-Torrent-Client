import React, { useCallback, memo } from 'react';
import type { PacketDirection } from '@/types/model/models.ts';
import type { PacketRowProps } from '@/features/packetInspector/types/packetInspectorTypes';
import { formatTime } from '@/features/packetInspector/utils/packetFormatUtils';
import styles from './PacketRow.module.css';

const DIRECTION_ICONS: Record<PacketDirection, string> = {
    incoming: '[<-]',
    outgoing: '[->]'
};

const DIRECTION_CLASSES: Record<PacketDirection, string> = {
    incoming: styles.incoming,
    outgoing: styles.outgoing
};

const getDirectionIcon = (direction: PacketDirection): string => {
    return DIRECTION_ICONS[direction];
};

const getDirectionClass = (direction: PacketDirection): string => {
    return DIRECTION_CLASSES[direction];
};

const getRowClass = (isSelected: boolean): string => {
    if (isSelected) {
        return `${styles.row} ${styles.rowSelected}`;
    }
    return styles.row;
};

const PacketRowComponent: React.FC<PacketRowProps> = ({ packet, isSelected, onSelect }) => {
    const handleClick = useCallback((): void => {
        onSelect(packet);
    }, [packet, onSelect]);

    return (
        <div className={getRowClass(isSelected)} onClick={handleClick}>
            <span className={styles.timestamp}>{formatTime(packet.timestamp)}</span>
            <span className={`${styles.direction} ${getDirectionClass(packet.direction)}`}>
                {getDirectionIcon(packet.direction)}
            </span>
            <span className={styles.ip}>{packet.peer_ip ?? 'SYSTEM'}</span>
            <span className={styles.type}>{packet.message_type}</span>
            <span className={styles.size}>{packet.size_bytes ?? 0} B</span>
        </div>
    );
};

const arePacketRowPropsEqual = (prevProps: PacketRowProps, nextProps: PacketRowProps): boolean => {
    return prevProps.packet.id === nextProps.packet.id && prevProps.isSelected === nextProps.isSelected;
};

export const PacketRow = memo(PacketRowComponent, arePacketRowPropsEqual);