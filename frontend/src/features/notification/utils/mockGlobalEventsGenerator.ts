import type { PeerConnectedEvent, ReputationUpdateEvent } from '@/types/model/models.ts';
import { randomFloat, randomHex, randomInt, randomIp } from '@/utils/random.ts';

export type GlobalEvent = PeerConnectedEvent | ReputationUpdateEvent;
export type GlobalEventCallback = (event: GlobalEvent) => void;

const generateRandomHex = (): string => {
    return randomHex(8);
};

export class MockGlobalEventsGenerator {
    private intervalId: number | null = null;
    private readonly callback: GlobalEventCallback;

    constructor(callback: GlobalEventCallback) {
        this.callback = callback;
    }

    public start(): void {
        if (!this.intervalId) {
            this.intervalId = window.setInterval(this.emitEvent.bind(this), 8000);
        }
    }

    public stop(): void {
        if (this.intervalId) {
            window.clearInterval(this.intervalId);
            this.intervalId = null;
        }
    }

    private emitEvent(): void {
        const isPeerEvent: boolean = randomFloat() > 0.5;
        if (isPeerEvent) {
            this.callback(this.createPeerEvent());
            return;
        }
        this.callback(this.createReputationEvent());
    }

    private createPeerEvent(): PeerConnectedEvent {
        const now: number = Date.now();
        return {
            event_type: 'peer:connected',
            info_hash: `INFO_${generateRandomHex()}`,
            peer: {
                peer_id: `Peer_ID_${generateRandomHex()}`,
                ip: randomIp(),
                port: randomInt(50000) + 10000,
                client_name: 'CipherNode/1.0',
                supports_ut_reputation: true,
                reputation_score: randomFloat(),
                flags: { encryption: true, utp: true, dht: true, pex: true },
                connected_since: new Date(now).toISOString(),
                last_seen: new Date(now).toISOString(),
                download_speed_bps: 0,
                upload_speed_bps: 0
            }
        };
    }

    private createReputationEvent(): ReputationUpdateEvent {
        return {
            event_type: 'reputation:update',
            from_peer_id: `Peer_TX_${generateRandomHex()}`,
            to_peer_id: `Peer_RX_${generateRandomHex()}`,
            new_reputation_score: randomFloat()
        };
    }
}