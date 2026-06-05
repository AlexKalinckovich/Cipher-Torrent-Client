import type { Torrent } from '@/types/model/models.ts';
import type { ColumnsType } from 'antd/es/table';
import React from 'react';

export interface DashboardStats {
    total: number;
    downloading: number;
    seeding: number;
    paused: number;
    error: number;
    totalSize: number;
}

export interface BadgeConfig {
    icon: React.ReactElement;
    text: string;
    className: string;
}

export interface FilterButtonProps {
    label: string;
    isActive: boolean;
    onClick: () => void;
    baseClassName: string;
    activeClassName: string;
}

export interface StatsSectionProps {
    stats: DashboardStats;
}

export interface ErrorStateProps {
    error: Error;
}

export interface FilterBarProps {
    currentFilter: string | null;
    searchText: string;
    onFilterChange: (status: string | null) => void;
    onSearchChange: (value: string) => void;
    baseClassName: string;
    activeClassName: string;
}

export interface TorrentTableProps {
    data: Torrent[];
    loading: boolean;
    columns: ColumnsType<Torrent>;
    onRowClick: (torrent: Torrent) => void;
}

export interface ProgressBarProps {
    progress: number;
}

export interface ActionButtonsProps {
    onPlay: () => void;
    onPause: () => void;
    onRemove: () => void;
    playClassName: string;
    pauseClassName: string;
    removeClassName: string;
}
