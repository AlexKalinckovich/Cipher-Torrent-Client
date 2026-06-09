import React, { memo } from 'react';
import type { DetailsPanelProps } from '../../types/packetInspectorTypes';
import { EmptyDetailsPanel } from './EmptyDetailsPanel';
import { ActiveDetailsPanel } from './ActiveDetailsPanel';

const DetailsPanelComponent: React.FC<DetailsPanelProps> = ({ packet }: DetailsPanelProps) => {

    if (!packet) {
        return <EmptyDetailsPanel />;
    }
    return <ActiveDetailsPanel packet={packet} />;
};

const areDetailsEqual = (prevProps: DetailsPanelProps, nextProps: DetailsPanelProps): boolean => {

    const isBothNull: boolean = prevProps.packet === null && nextProps.packet === null;
    if (isBothNull) {
        return true;
    }

    const isOneNull: boolean = prevProps.packet === null || nextProps.packet === null;
    if (isOneNull) {
        return false;
    }

    return prevProps.packet?.id === nextProps.packet?.id;
};

export const DetailsPanel = memo(DetailsPanelComponent, areDetailsEqual);