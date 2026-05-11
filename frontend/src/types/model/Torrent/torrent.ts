import type { TorrentFile } from './torrentFile';

export type TorrentStatus = 'downloading' | 'seeding' | 'paused' | 'error';

export interface Torrent {
    id: number;
    info_hash: string;
    name: string;
    size_bytes: number;
    status: TorrentStatus;
    progress?: number;
    download_speed_bps?: number;
    upload_speed_bps?: number;
    added_at: string;
    piece_length?: number;
    files?: Array<TorrentFile>;
    signatures_count?: number;
    is_signed_by_me?: boolean;
    peers_count?: number;
}