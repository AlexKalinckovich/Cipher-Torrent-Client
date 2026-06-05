import React from 'react';
import type { NotificationArgsProps } from 'antd';
import type { PeerConnectedEvent, ReputationUpdateEvent } from '@/types/model/models.ts';
import styles from '../components/GlobalNotifications.module.css';

const renderPeerContent: (event: PeerConnectedEvent) => React.ReactNode = (event: PeerConnectedEvent): React.ReactNode => (
    <div className={styles.message}>
        <div className={styles.row}>
            <span className={styles.label}>ENDPOINT:</span>
            <span className={styles.value}>{event.peer.ip}:{event.peer.port}</span>
        </div>
        <div className={styles.row}>
            <span className={styles.label}>CLIENT:</span>
            <span className={styles.value}>{event.peer.client_name}</span>
        </div>
    </div>
);

const renderReputationContent: (event: ReputationUpdateEvent) => React.ReactNode = (event: ReputationUpdateEvent): React.ReactNode => (
    <div className={styles.message}>
        <div className={styles.row}>
            <span className={styles.label}>TARGET NODE:</span>
            <span className={styles.value}>{event.to_peer_id.substring(0, 8)}...</span>
        </div>
        <div className={styles.row}>
            <span className={styles.label}>NEW SCORE:</span>
            <span className={styles.value}>{event.new_reputation_score.toFixed(2)}</span>
        </div>
    </div>
);

export const getPeerNotificationConfig = (event: PeerConnectedEvent): NotificationArgsProps => ({
    title: <span className={styles.titlePeer}>SYSTEM ALARM: PEER CONNECTED</span>,
    description: renderPeerContent(event),
    className: styles.toastPeer,
    placement: 'bottomRight',
    duration: 5,
});

export const getReputationNotificationConfig = (event: ReputationUpdateEvent): NotificationArgsProps => ({
    title: <span className={styles.titleReputation}>REPUTATION SHIFT DETECTED</span>,
    description: renderReputationContent(event),
    className: styles.toastReputation,
    placement: 'bottomRight',
    duration: 5,
});