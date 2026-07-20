import type { PacketLog } from './packetLog';

export type PacketEventType = 'packet';

export interface PacketEvent {
    event_type: PacketEventType;
    packet: PacketLog;
}

export type ProgressEventType = 'torrent:progress';

export interface ProgressLog {
    info_hash: string;
    progress: number;
    creator_public_key: string;
    download_speed_bps: number;
    upload_speed_bps: number;
}

export interface ProgressEvent {
    event_type: ProgressEventType;
    progress: ProgressLog;
}