import type {TorrentSignature} from './torrentSignature';

export type TorrentSignedEventType = 'torrent:signed';

export interface TorrentSignedEvent {
    event_type: TorrentSignedEventType;
    info_hash: string;
    signature: TorrentSignature;
}