import type { Peer } from './peer';

export type PeerConnectedEventType = 'peer:connected';

export interface PeerConnectedEvent {
    event_type: PeerConnectedEventType;
    info_hash: string;
    peer: Peer;
}