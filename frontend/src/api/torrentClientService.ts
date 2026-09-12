import WebTorrent, { type Torrent, type TorrentFile, type TorrentOptions } from './webtorrent';
import { attachTorrentPacketListener, type WirePacketListener } from './torrentPacketCapture';

export interface TorrentProgressSnapshot {
    infoHash: string;
    name: string;
    progress: number;          // 0..1
    downloadSpeed: number;     // bytes/sec
    uploadSpeed: number;       // bytes/sec
    downloaded: number;
    length: number;
    numPeers: number;
    done: boolean;
    files: TorrentFile[];
}

export interface TorrentClientCallbacks {
    onProgress?: (snapshot: TorrentProgressSnapshot) => void;
    onDone?: (snapshot: TorrentProgressSnapshot) => void;
    onError?: (infoHash: string, error: Error | string) => void;
}

const DEFAULT_ANNOUNCE = [
    'wss://tracker.openwebtorrent.com',
    'wss://tracker.btorrent.xyz',
    'wss://tracker.fastcast.nz',
];

/**
 * Thin client-side wrapper around WebTorrent.
 *
 * Responsibilities:
 *  - create a single WebTorrent instance (browser WebRTC + WebSocket trackers)
 *  - add a torrent from a parsed .torrent file or a magnet URI
 *  - track per-torrent progress/peers/speeds and notify subscribers
 *  - expose the torrent files for saving locally (File System Access API)
 */
export class TorrentClientService {
    private readonly client: InstanceType<typeof WebTorrent>;
    private readonly torrents: Map<string, Torrent>;
    private readonly callbacks: Map<string, TorrentClientCallbacks>;
    private readonly announce: string[];

    constructor(announce: string[] = DEFAULT_ANNOUNCE) {
        this.client = new WebTorrent();
        this.torrents = new Map();
        this.callbacks = new Map();
        this.announce = announce;
    }

    private buildSnapshot(t: Torrent): TorrentProgressSnapshot {
        return {
            infoHash: t.infoHash,
            name: t.name,
            progress: t.progress,
            downloadSpeed: t.downloadSpeed,
            uploadSpeed: t.uploadSpeed,
            downloaded: t.downloaded,
            length: t.length,
            numPeers: t.numPeers,
            done: t.done,
            files: t.files,
        };
    }

    private bindTorrent(t: Torrent): void {
        const infoHash = t.infoHash;
        this.torrents.set(infoHash, t);

        const notify = () => {
            this.callbacks.get(infoHash)?.onProgress?.(this.buildSnapshot(t));
        };

        t.on('download', notify);
        t.on('upload', notify);
        t.on('wire', notify);
        t.on('done', () => {
            this.callbacks.get(infoHash)?.onDone?.(this.buildSnapshot(t));
            this.callbacks.get(infoHash)?.onProgress?.(this.buildSnapshot(t));
        });
        t.on('error', (err: Error | string) => {
            this.callbacks.get(infoHash)?.onError?.(infoHash, err);
        });
    }

    /**
     * Add a torrent from a magnet URI, a raw .torrent byte buffer,
     * a File/Blob, or an already-parsed torrent instance.
     */
    public add(
        torrentId: string | Uint8Array | File | Blob,
        callbacks?: TorrentClientCallbacks,
    ): Promise<Torrent> {
        return new Promise<Torrent>((resolve, reject) => {
            const opts: TorrentOptions = { announce: this.announce };

            const t = this.client.add(torrentId, opts, (added: Torrent) => {
                this.bindTorrent(added);
                if (callbacks) {
                    this.callbacks.set(added.infoHash, callbacks);
                }
                resolve(added);
            });

            // If add throws synchronously (bad torrent), reject.
            if (t === undefined) {
                reject(new Error('Failed to add torrent'));
            }
        });
    }

    public async get(infoHash: string): Promise<Torrent | null> {
        const cached = this.torrents.get(infoHash);
        if (cached) {
            return cached;
        }
        const t = await this.client.get(infoHash);
        if (t) {
            this.bindTorrent(t);
        }
        return t;
    }

    public async remove(infoHash: string): Promise<void> {
        const t = this.torrents.get(infoHash);
        if (t) {
            await this.client.remove(t);
            this.torrents.delete(infoHash);
            this.callbacks.delete(infoHash);
        }
    }

    public subscribe(infoHash: string, callbacks: TorrentClientCallbacks): void {
        this.callbacks.set(infoHash, callbacks);
    }

    public unsubscribe(infoHash: string): void {
        this.callbacks.delete(infoHash);
    }

    /**
     * Subscribe to progress snapshots for a torrent. Returns an unsubscribe fn.
     */
    public subscribeProgress(
        infoHash: string,
        onProgress: (snapshot: TorrentProgressSnapshot) => void,
    ): () => void {
        const existing = this.snapshot(infoHash);
        if (existing) {
            onProgress(existing);
        }
        const prev = this.callbacks.get(infoHash);
        this.callbacks.set(infoHash, {
            ...prev,
            onProgress,
        });
        return () => {
            const current = this.callbacks.get(infoHash);
            if (current && current.onProgress === onProgress) {
                this.callbacks.delete(infoHash);
            }
        };
    }

    public snapshot(infoHash: string): TorrentProgressSnapshot | null {
        const t = this.torrents.get(infoHash);
        return t ? this.buildSnapshot(t) : null;
    }

    /**
     * Attach a packet listener to a torrent (if it is currently active).
     * Returns a cleanup function, or undefined when the torrent is not loaded.
     */
    public async attachPacketListener(
        infoHash: string,
        onPacket: WirePacketListener,
    ): Promise<(() => void) | undefined> {
        const t = await this.get(infoHash);
        if (!t) {
            return undefined;
        }
        return attachTorrentPacketListener(t, onPacket);
    }

    public destroy(): void {
        this.client.destroy();
        this.torrents.clear();
        this.callbacks.clear();
    }
}

// Singleton instance for the app.
export const torrentClientService = new TorrentClientService();