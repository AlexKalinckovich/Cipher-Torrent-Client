import type { TorrentFile } from './torrentFile';
import type {SignatureDTO} from "@/types/model/Torrent/torrentSignature.ts";

export type TorrentStatus = 'idle' | 'downloading' | 'seeding' | 'paused';

export interface TorrentEntity {
    name: string;
    size_bytes: number;
    piece_length: number;
    is_private: boolean;
    storage_path: string;
    added_at: string;
}

export interface TorrentDTO {
    info_hash: string;
    creator_public_key : string
    name: string;
    size_bytes: number;
    storage_path: string;
    status: TorrentStatus;
    progress: number;
    added_at: string;
    files: TorrentFile[];
    signatures: SignatureDTO[];
}