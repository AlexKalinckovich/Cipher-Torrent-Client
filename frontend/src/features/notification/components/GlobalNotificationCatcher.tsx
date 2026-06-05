import React, { useEffect, memo, useCallback } from 'react';
import { notification } from 'antd';
import type { NotificationInstance } from 'antd/es/notification/interface';
import { MockGlobalEventsGenerator } from '../utils/mockGlobalEventsGenerator';
import { getPeerNotificationConfig, getReputationNotificationConfig } from '../utils/notificationConfig';
import type { GlobalEvent } from '../utils/mockGlobalEventsGenerator';
import type { PeerConnectedEvent, ReputationUpdateEvent } from '@/types/model/models.ts';

type EventDispatcher = (api: NotificationInstance, event: GlobalEvent) => void;

const dispatchPeer = (api: NotificationInstance, event: GlobalEvent): void => {
    api.open(getPeerNotificationConfig(event as PeerConnectedEvent));
};

const dispatchReputation = (api: NotificationInstance, event: GlobalEvent): void => {
    api.open(getReputationNotificationConfig(event as ReputationUpdateEvent));
};

const DISPATCH_STRATEGIES: Record<string, EventDispatcher> = {
    'peer:connected': dispatchPeer,
    'reputation:update': dispatchReputation
};

const GlobalNotificationCatcherComponent: React.FC = () => {
    const [api, contextHolder] = notification.useNotification();

    const handleEvent: (event: GlobalEvent) => void = useCallback((event: GlobalEvent): void => {
        const dispatcher: EventDispatcher | undefined = DISPATCH_STRATEGIES[event.event_type];
        if (dispatcher) {
            dispatcher(api, event);
        }
    }, [api]);

    useEffect((): (() => void) => {
        const generator = new MockGlobalEventsGenerator(handleEvent);
        generator.start();

        return (): void => {
            generator.stop();
        };
    }, [handleEvent]);

    return <>{contextHolder}</>;
};

export const GlobalNotificationCatcher = memo(GlobalNotificationCatcherComponent);