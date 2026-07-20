import type { PacketEvent, ProgressEvent } from '@/types/model/models.ts';

export type StreamEvent = PacketEvent | ProgressEvent;
export type StreamCallback = (event: StreamEvent) => void;

class PacketStreamService {
    private ws: WebSocket | null = null;
    private subscriber: StreamCallback | null = null;
    private reconnectTimeout: number | null = null;
    private infoHash: string | null = null;
    private intentionallyClosed = false;

    public connect(infoHash: string, callback: StreamCallback): void {
        if (
            this.ws &&
            this.infoHash === infoHash &&
            (this.ws.readyState === WebSocket.OPEN || this.ws.readyState === WebSocket.CONNECTING)
        ) {
            this.subscriber = callback;
            return;
        }

        if (this.ws) {
            this.disconnect();
        }

        this.intentionallyClosed = false;
        this.infoHash = infoHash;
        this.subscriber = callback;

        const wsUrl = `ws://localhost:8080/api/v1/ws/torrent/${infoHash}`;
        const socket = new WebSocket(wsUrl);
        this.ws = socket;

        socket.onopen = () => {
            console.log(`[WS] ✅ Connected to ${infoHash}`);
        };

        socket.onmessage = (event) => {
            if (this.ws !== socket) return; // stale socket
            try {
                const parsedEvent: StreamEvent = JSON.parse(event.data);
                this.subscriber?.(parsedEvent);
            } catch (err) {
                console.error('[WS] ❌ Failed to parse message:', err);
            }
        };

        socket.onclose = () => {
            if (this.ws !== socket) return;
            console.log('[WS] ⚠️ Disconnected');
            this.ws = null;
            if (!this.intentionallyClosed) {
                this.scheduleReconnect();
            }
            this.intentionallyClosed = false;
        };

        socket.onerror = (error) => {
            console.error('[WS] ❌ Error:', error);
            socket.close();
        };
    }

    private scheduleReconnect(): void {
        if (this.reconnectTimeout || !this.infoHash || !this.subscriber) return;

        this.reconnectTimeout = window.setTimeout(() => {
            this.reconnectTimeout = null;
            if (this.infoHash && this.subscriber) {
                this.connect(this.infoHash, this.subscriber);
            }
        }, 3000);
    }

    public disconnect(): void {
        this.intentionallyClosed = true;
        if (this.reconnectTimeout) {
            clearTimeout(this.reconnectTimeout);
            this.reconnectTimeout = null;
        }
        if (this.ws) {
            this.ws.close();
            this.ws = null;
        }
        this.subscriber = null;
        this.infoHash = null;
    }
}

export const packetStreamService = new PacketStreamService();