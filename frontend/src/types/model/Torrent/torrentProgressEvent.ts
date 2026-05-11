export type TorrentProgressEventType = 'torrent:progress';

export interface TorrentProgressEvent {
    event_type: TorrentProgressEventType;
    info_hash: string;
    progress: number;
    download_speed: number;
    upload_speed: number;
}