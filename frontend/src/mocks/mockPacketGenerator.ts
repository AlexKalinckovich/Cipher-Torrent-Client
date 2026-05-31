import type { PacketEvent, PacketLog, PacketDirection } from '@/types/model/models.ts';

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
        return {
            id: id,
            timestamp: new Date().toISOString(),
            direction: this.getRandomDirection(),
            message_type: this.getRandomMessageType(),
            peer_ip: this.getRandomIp(),
            raw_payload_base64: 'SGVsbG8gQ3liZXJwdW5r',
            size_bytes: this.getRandomSize()
        };
    }

    private getRandomDirection(): PacketDirection {
        return Math.random() > 0.5 ? 'incoming' : 'outgoing';
    }

    private getRandomMessageType(): string {
        const types: string[] = ['handshake', 'choke', 'unchoke', 'interested', 'bitfield', 'request', 'piece', 'cancel'];
        const index: number = Math.floor(Math.random() * types.length);
        return types[index];
    }

    private getRandomIp(): string {
        const p1: number = Math.floor(Math.random() * 255);
        const p2: number = Math.floor(Math.random() * 255);
        return `192.168.${p1}.${p2}`;
    }

    private getRandomSize(): number {
        return Math.floor(Math.random() * 16384) + 16;
    }
}