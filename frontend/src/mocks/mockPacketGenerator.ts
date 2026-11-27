import type { PacketEvent, PacketLog, PacketDirection, ReputationReceipt } from '@/types/model/models.ts';
import { randomFloat, randomHex, randomInt, randomIp } from '@/utils/random.ts';

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
        const hex: string = randomHex(8);
        return {
            from_peer_id: `Peer_TX_${hex}`,
            to_peer_id: `Peer_RX_${hex}`,
            piece_index: randomInt(2000),
            byte_count: randomInt(65536) + 1024,
            signature: `ED25519_SIG_AB${randomHex(13)}F9`,
            timestamp: new Date().toISOString()
        };
    }

    private getRandomDirection(): PacketDirection {
        if (randomFloat() > 0.5) {
            return 'incoming';
        }
        return 'outgoing';
    }

    private getRandomMessageType(): string {
        const types: string[] = ['handshake', 'choke', 'unchoke', 'interested', 'request', 'piece', 'ut_reputation'];
        const index: number = randomInt(types.length);
        return types[index];
    }

    private getRandomIp(): string {
        return randomIp();
    }

    private getRandomSize(): number {
        return randomInt(1500) + 40;
    }
}