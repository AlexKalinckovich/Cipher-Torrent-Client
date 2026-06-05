import type { PacketEvent, PacketLog, PacketDirection, ReputationReceipt } from '@/types/model/models.ts';

export type PacketCallback = (event: PacketEvent) => void;

export class MockPacketGenerator {
    private intervalId: number | null = null;
    private readonly callback: PacketCallback;
    private counter: number = 0;

    constructor(callback: PacketCallback) {
        this.callback = callback;
    }

    public start(): void {
        if (!this.intervalId) {
            this.intervalId = window.setInterval(this.emitPacket.bind(this), 600);
        }
    }

    public stop(): void {
        if (this.intervalId) {
            window.clearInterval(this.intervalId);
            this.intervalId = null;
        }
    }

    private emitPacket(): void {
        this.counter += 1;
        const event: PacketEvent = this.createRandomEvent(this.counter);
        this.callback(event);
    }

    private createRandomEvent(id: number): PacketEvent {
        return {
            event_type: 'packet',
            packet: this.createRandomPacket(id)
        };
    }

    private createRandomPacket(id: number): PacketLog {
        const messageType: string = this.getRandomMessageType();
        return {
            id: id,
            timestamp: new Date().toISOString(),
            direction: this.getRandomDirection(),
            message_type: messageType,
            peer_ip: this.getRandomIp(),
            raw_payload_base64: 'SGVsbG8gQ3liZXJwdW5r',
            size_bytes: this.getRandomSize(),
            has_reputation_extension: messageType === 'ut_reputation',
            reputation_payload: this.checkAndAttachReputation(messageType)
        };
    }

    private checkAndAttachReputation(messageType: string): ReputationReceipt | undefined {
        if (messageType === 'ut_reputation') {
            return this.createRandomReputationReceipt();
        }
        return undefined;
    }

    private createRandomReputationReceipt(): ReputationReceipt {
        const randomHex: string = Math.random().toString(16).substring(2, 10).toUpperCase();
        return {
            from_peer_id: `Peer_TX_${randomHex}`,
            to_peer_id: `Peer_RX_${randomHex}`,
            piece_index: Math.floor(Math.random() * 2000),
            byte_count: Math.floor(Math.random() * 65536) + 1024,
            signature: `ED25519_SIG_AB${Math.random().toString(36).substring(2, 15).toUpperCase()}F9`,
            timestamp: new Date().toISOString()
        };
    }

    private getRandomDirection(): PacketDirection {
        if (Math.random() > 0.5) {
            return 'incoming';
        }
        return 'outgoing';
    }

    private getRandomMessageType(): string {
        const types: string[] = ['handshake', 'choke', 'unchoke', 'interested', 'request', 'piece', 'ut_reputation'];
        const index: number = Math.floor(Math.random() * types.length);
        return types[index];
    }

    private getRandomIp(): string {
        const p1: number = Math.floor(Math.random() * 255);
        const p2: number = Math.floor(Math.random() * 255);
        const p3: number = Math.floor(Math.random() * 255);
        const p4: number = Math.floor(Math.random() * 255);
        return `${p1}.${p2}.${p3}.${p4}`;
    }

    private getRandomSize(): number {
        return Math.floor(Math.random() * 1500) + 40;
    }
}