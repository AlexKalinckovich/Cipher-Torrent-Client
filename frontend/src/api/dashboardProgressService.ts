// api/dashboardProgressService.ts
import type { ProgressLog } from '@/types/model/models';

type ProgressCallback = (update: ProgressLog) => void;

class DashboardProgressService {
    private ws: WebSocket | null = null;
    private subscribers = new Set<ProgressCallback>();
    private reconnectTimeout: number | null = null;
    private intentionallyClosed = false;

    public connect(): void {
        if (this.ws?.readyState === WebSocket.OPEN) return;
        if (this.ws?.readyState === WebSocket.CONNECTING) return;

        this.intentionallyClosed = false;
        this.ws = new WebSocket('ws://localhost:8080/api/v1/ws/dashboard');

        this.ws.onopen = () => {
            console.log('[DashboardWS] Connected');
        };

        this.ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                // Expecting { event_type: "torrent:progress", progress: ProgressLog }
                if (data.event_type === 'torrent:progress' && data.progress) {
                    this.subscribers.forEach((cb) => cb(data.progress));
                }
            } catch (err) {
                console.error('[DashboardWS] Parse error:', err);
            }
        };

        this.ws.onclose = () => {
            console.log('[DashboardWS] Disconnected');
            this.ws = null;
            if (!this.intentionallyClosed) {
                this.scheduleReconnect();
            }
        };

        this.ws.onerror = (err) => {
            console.error('[DashboardWS] Error:', err);
            this.ws?.close(); // trigger reconnect
        };
    }

    public subscribe(callback: ProgressCallback): () => void {
        this.subscribers.add(callback);
        if (!this.ws || this.ws.readyState === WebSocket.CLOSED) {
            this.connect();
        }
        return () => {
            this.subscribers.delete(callback);
            if (this.subscribers.size === 0) {
                this.disconnect();
            }
        };
    }

    private scheduleReconnect(): void {
        if (this.reconnectTimeout) return;
        this.reconnectTimeout = window.setTimeout(() => {
            this.reconnectTimeout = null;
            this.connect();
        }, 3000);
    }

    public disconnect(): void {
        this.intentionallyClosed = true;
        if (this.reconnectTimeout) {
            clearTimeout(this.reconnectTimeout);
            this.reconnectTimeout = null;
        }
        this.ws?.close();
        this.ws = null;
        this.subscribers.clear();
    }
}

export const dashboardProgressService = new DashboardProgressService();