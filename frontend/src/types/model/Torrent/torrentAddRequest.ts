export interface TorrentAddRequest {
    magnet_uri?: string;
    file_base64?: string;
    save_path?: string;
}
export interface TorrentIdentity {
    info_hash: string;       // Hex string
    creator_pub_key: string; // Base64 string
}

export interface ProgressUpdateRequest extends TorrentIdentity {
    progress: number;
}

export type TorrentProgressEventType = 'torrent:progress';

export interface TorrentProgressLog {
    info_hash: string;
    progress: number;
    download_speed_bps: number;
    upload_speed_bps: number;
}

export interface TorrentProgressEvent {
    event_type: TorrentProgressEventType;
    progress: TorrentProgressLog;
}