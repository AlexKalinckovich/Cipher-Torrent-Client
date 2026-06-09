import React, { useEffect, memo, useCallback } from 'react';
import { notification } from 'antd';
import type { NotificationInstance } from 'antd/es/notification/interface';
import { wsClient } from '@/api/webSocketClient.ts';
import { getPeerNotificationConfig, getReputationNotificationConfig, getTorrentSignedNotificationConfig } from '../utils/notificationConfig';
import type { GlobalEvent } from '@/api/webSocketClient.ts';
import type { PeerConnectedEvent, ReputationUpdateEvent, TorrentSignedEvent } from '@/types/model/models.ts';

type EventDispatcher = (api: NotificationInstance, event: GlobalEvent) => void;

const dispatchPeer = (api: NotificationInstance, event: GlobalEvent): void => {
    api.open(getPeerNotificationConfig(event as PeerConnectedEvent));
};

const dispatchReputation = (api: NotificationInstance, event: GlobalEvent): void => {
    api.open(getReputationNotificationConfig(event as ReputationUpdateEvent));
};

const dispatchSigned = (api: NotificationInstance, event: GlobalEvent): void => {
    api.open(getTorrentSignedNotificationConfig(event as TorrentSignedEvent));
};

const DISPATCH_STRATEGIES: Record<string, EventDispatcher> = {
    'peer:connected': dispatchPeer,
    'reputation:update': dispatchReputation,
    'torrent:signed': dispatchSigned
};

const GlobalNotificationCatcherComponent: React.FC = () => {

    const [api, contextHolder] = notification.useNotification();

    const handleEvent = useCallback((event: GlobalEvent): void => {
        const dispatcher: EventDispatcher | undefined = DISPATCH_STRATEGIES[event.event_type];
        if (dispatcher) {
            dispatcher(api, event);
        }
    }, [api]);

    useEffect((): (() => void) => {
        wsClient.subscribe(handleEvent);
        return (): void => {
            wsClient.unsubscribe(handleEvent);
        };
    }, [handleEvent]);

    return <>{contextHolder}</>;
};

export const GlobalNotificationCatcher = memo(GlobalNotificationCatcherComponent);