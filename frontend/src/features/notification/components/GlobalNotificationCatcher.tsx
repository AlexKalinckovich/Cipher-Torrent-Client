import React, { useEffect, memo } from 'react';
import { MockGlobalEventsGenerator } from '../utils/mockGlobalEventsGenerator';
import { dispatchNotification } from '../utils/notificationDispatcher';

const GlobalNotificationCatcherComponent: React.FC = () => {
    useEffect((): (() => void) => {
        const generator = new MockGlobalEventsGenerator(dispatchNotification);
        generator.start();

        return (): void => {
            generator.stop();
        };
    }, []);

    return null;
};

export const GlobalNotificationCatcher = memo(GlobalNotificationCatcherComponent);