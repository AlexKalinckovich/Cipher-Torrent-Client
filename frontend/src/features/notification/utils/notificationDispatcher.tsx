import React from 'react';
import { notification } from 'antd';
import type { PeerConnectedEvent, ReputationUpdateEvent } from '@/types/model/models.ts';
import type { GlobalEvent } from './mockGlobalEventsGenerator';
import styles from '../components/GlobalNotifications.module.css';

type EventHandler = (event: GlobalEvent) => void;

const renderPeerContent = (event: PeerConnectedEvent): React.ReactNode => (
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

const renderReputationContent = (event: ReputationUpdateEvent): React.ReactNode => (
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


const displayPeerNotification = (event: PeerConnectedEvent): void => {
    notification.open({
        title: <span className={styles.titlePeer}>SYSTEM ALARM: PEER CONNECTED</span>,
        description: renderPeerContent(event),
        className: styles.toastPeer,
        placement: 'topLeft',
        duration: 5,
    });
};

const displayReputationNotification = (event: ReputationUpdateEvent): void => {
    notification.open({
        title: <span className={styles.titleReputation}>REPUTATION SHIFT DETECTED</span>,
        description: renderReputationContent(event),
        className: styles.toastReputation,
        placement: 'topLeft',
        duration: 5,
    });
};

const EVENT_HANDLERS: Record<string, EventHandler> = {
    'peer:connected': (e: GlobalEvent): void => displayPeerNotification(e as PeerConnectedEvent),
    'reputation:update': (e: GlobalEvent): void => displayReputationNotification(e as ReputationUpdateEvent),
};

export const dispatchNotification = (event: GlobalEvent): void => {
    const handler: EventHandler | undefined = EVENT_HANDLERS[event.event_type];
    if (handler) {
        handler(event);
    }
};