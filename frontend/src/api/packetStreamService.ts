import type { PacketEvent } from '@/types/model/models.ts';
import { MockPacketGenerator, type PacketCallback } from '@/mocks/mockPacketGenerator.ts';

export class PacketStreamService {
    private readonly mockGenerator: MockPacketGenerator;
    private subscriber: PacketCallback | null = null;

    constructor() {
        this.mockGenerator = new MockPacketGenerator(this.handleIncomingEvent.bind(this));
    }

    public connect(callback: PacketCallback): void {
        this.subscriber = callback;
        this.mockGenerator.start();
    }

    public disconnect(): void {
        this.mockGenerator.stop();
        this.subscriber = null;
    }

    private handleIncomingEvent(event: PacketEvent): void {
        if (this.subscriber) {
            this.subscriber(event);
        }
    }
}

export const packetStreamService = new PacketStreamService();