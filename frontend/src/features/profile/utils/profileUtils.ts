import type { StatItem } from '@/features/profile/types/profileTypes';
import type { UserStats } from '@/types/model/models.ts';

const TRAFFIC_COLORS: string[] = ['#3b82f6', '#ef4444'];
const ACTIVITY_COLORS: string[] = ['#10b981', '#f59e0b', '#8b5cf6'];
const KILOBYTE: number = 1024;
const BYTE_SIZES: string[] = ['B', 'KB', 'MB', 'GB', 'TB'];

const getByteIndex = (bytes: number): number => {
    return Math.floor(Math.log(bytes) / Math.log(KILOBYTE));
};

const calculateByteValue = (bytes: number, index: number): string => {
    const value: number = parseFloat((bytes / Math.pow(KILOBYTE, index)).toFixed(2));
    return `${value} ${BYTE_SIZES[index]}`;
};

export const formatBytes = (bytes: number): string => {
    if (bytes === 0) {
        return '0 B';
    }
    return calculateByteValue(bytes, getByteIndex(bytes));
};

const isBytesField = (name: string): boolean => {
    const lowerName: string = name.toLowerCase();
    return lowerName.includes('downloaded') || lowerName.includes('uploaded');
};

export const formatStatValue = (name: string, value: number): string => {
    if (isBytesField(name)) {
        return formatBytes(value);
    }
    return value.toString();
};

export const buildTrafficData = (stats?: UserStats): StatItem[] => {
    if (!stats) {
        return [];
    }
    return [
        { name: 'Downloaded', value: stats.total_downloaded_bytes, color: TRAFFIC_COLORS[0] },
        { name: 'Uploaded', value: stats.total_uploaded_bytes, color: TRAFFIC_COLORS[1] }
    ];
};

export const buildActivityData = (stats?: UserStats): StatItem[] => {
    if (!stats) {
        return [];
    }
    return [
        { name: 'Signed Torrents', value: stats.signed_torrents_count, color: ACTIVITY_COLORS[0] },
        { name: 'Active Torrents', value: stats.active_torrents_count, color: ACTIVITY_COLORS[1] },
        { name: 'Trusted Peers', value: stats.peers_trusted_count, color: ACTIVITY_COLORS[2] }
    ];
};

export const formatDate = (isoString: string): string => {
    const date = new Date(isoString);
    return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'short',
        day: '2-digit'
    });
};