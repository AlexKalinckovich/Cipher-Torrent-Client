import type {
    PeerConnectedEvent,
    ReputationUpdateEvent,
    TorrentProgressEvent,
    TorrentSignedEvent
} from '../types/model/models.ts';

export type GlobalEvent = PeerConnectedEvent | ReputationUpdateEvent | TorrentProgressEvent | TorrentSignedEvent;
export type EventCallback = (event: GlobalEvent) => void;

class WebSocketClient {
    private listeners: EventCallback[] = [];

    public subscribe(callback: EventCallback): void {
        this.listeners.push(callback);
    }

    public unsubscribe(callback: EventCallback): void {
        this.listeners = this.listeners.filter((cb: EventCallback): boolean => cb !== callback);
    }

    public emit(event: GlobalEvent): void {
        this.listeners.forEach((cb: EventCallback): void => cb(event));
    }
}

export const wsClient = new WebSocketClient();

const generateRandomHex = (): string => Math.random().toString(16).substring(2, 10).toUpperCase();

class MockWSEngine {
    private progressInterval: number | null = null;
    private eventInterval: number | null = null;

    public start(): void {
        if (!this.progressInterval) {
            this.progressInterval = window.setInterval(this.emitProgress.bind(this), 1000);
            this.eventInterval = window.setInterval(this.emitRandomEvent.bind(this), 8000);
        }
    }

    public stop(): void {
        if (this.progressInterval) {
            window.clearInterval(this.progressInterval);
            this.progressInterval = null;
        }
        if (this.eventInterval) {
            window.clearInterval(this.eventInterval);
            this.eventInterval = null;
        }
    }

    private emitProgress(): void {
        const eventA: TorrentProgressEvent = {
            event_type: 'torrent:progress',
            info_hash: 'A1B2C3D4E5F6G7H8I9J0',
            progress: Math.min(100, +((75.4 + Math.random() * 3).toFixed(2))),
            download_speed: Math.floor(Math.random() * 5000000) + 1000000,
            upload_speed: Math.floor(Math.random() * 500000)
        };
        const eventB: TorrentProgressEvent = {
            event_type: 'torrent:progress',
            info_hash: 'B2C3D4E5F6G7H8I9J0K1',
            progress: Math.min(100, +((45.2 + Math.random() * 3).toFixed(2))),
            download_speed: Math.floor(Math.random() * 8000000) + 2000000,
            upload_speed: Math.floor(Math.random() * 100000)
        };
        wsClient.emit(eventA);
        wsClient.emit(eventB);
    }

    private emitRandomEvent(): void {
        const rand: number = Math.random();
        if (rand > 0.66) {
            this.emitSignedEvent();
            return;
        }
        if (rand > 0.33) {
            this.emitReputationEvent();
            return;
        }
        this.emitPeerEvent();
    }

    private emitSignedEvent(): void {
        const event: TorrentSignedEvent = {
            event_type: 'torrent:signed',
            info_hash: `INFO_${generateRandomHex()}`,
            signature: {
                signer_user_id: Math.floor(Math.random() * 100),
                signer_public_key: `Ed25519_PubKey_Node_${generateRandomHex()}`,
                signature_bytes: `SIG_${generateRandomHex()}`,
                is_valid: true,
                signed_at: new Date().toISOString()
            }
        };
        wsClient.emit(event);
    }

    private emitReputationEvent(): void {
        const event: ReputationUpdateEvent = {
            event_type: 'reputation:update',
            from_peer_id: `Peer_TX_${generateRandomHex()}`,
            to_peer_id: `Peer_RX_${generateRandomHex()}`,
            new_reputation_score: Math.random()
        };
        wsClient.emit(event);
    }

    private emitPeerEvent(): void {
        const event: PeerConnectedEvent = {
            event_type: 'peer:connected',
            info_hash: `INFO_${generateRandomHex()}`,
            peer: {
                peer_id: `Peer_ID_${generateRandomHex()}`,
                ip: '192.168.1.100',
                port: 51413,
                client_name: 'CipherNode/1.0',
                supports_ut_reputation: true,
                reputation_score: Math.random(),
                flags: { encryption: true, utp: true, dht: true, pex: true },
                connected_since: new Date().toISOString(),
                last_seen: new Date().toISOString(),
                download_speed_bps: 0,
                upload_speed_bps: 0
            }
        };
        wsClient.emit(event);
    }
}

const isDev: boolean = import.meta.env ? import.meta.env.DEV : true;

if (isDev) {
    new MockWSEngine().start();
}