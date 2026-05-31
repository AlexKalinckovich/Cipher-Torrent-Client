import React from 'react';
import type { DetailsPanelProps } from '../../types/packetInspectorTypes';
import { EmptyDetailsPanel } from './EmptyDetailsPanel';
import { ActiveDetailsPanel } from './ActiveDetailsPanel';

export const DetailsPanel: React.FC<DetailsPanelProps> = ({ packet }) => {
    if (!packet) {
        return <EmptyDetailsPanel />;
    }
    return <ActiveDetailsPanel packet={packet} />;
};