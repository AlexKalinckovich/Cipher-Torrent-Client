import React, { memo } from 'react';
import type { ActiveDetailsPanelProps } from '../../types/packetInspectorTypes';
import { OptionalReputationBlock } from './ReputationBlock';
import styles from './DetailsPanel.module.css';

const ActiveDetailsPanelComponent: React.FC<ActiveDetailsPanelProps> = ({ packet }) => {
    return (
        <div className={styles.detailsPanel}>
            <div className={styles.detailsTitle}>Packet Details [{packet.id}]</div>

            <div className={styles.dataBlock}>
                <span className={styles.dataLabel}>Peer Node IP</span>
                <div className={styles.dataContent}>{packet.peer_ip ?? 'Internal System Packet'}</div>
            </div>

            <div className={styles.dataBlock}>
                <span className={styles.dataLabel}>Parsed Protocol Info</span>
                <div className={styles.dataContent}>
                    {`Type: ${packet.message_type.toUpperCase()}\nSize: ${packet.size_bytes} Bytes\nDirection: ${packet.direction}`}
                </div>
            </div>

            <div className={styles.dataBlock}>
                <span className={styles.dataLabel}>Raw Payload (Base64)</span>
                <div className={styles.dataContent}>{packet.raw_payload_base64}</div>
            </div>

            <OptionalReputationBlock payload={packet.reputation_payload} />
        </div>
    );
};

const areActiveDetailsEqual = (prevProps: ActiveDetailsPanelProps, nextProps: ActiveDetailsPanelProps): boolean => {
    return prevProps.packet.id === nextProps.packet.id;
};

export const ActiveDetailsPanel = memo(ActiveDetailsPanelComponent, areActiveDetailsEqual);