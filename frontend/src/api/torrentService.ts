import { api } from './axiosClient.ts';
import type {
    TorrentDTO,
    TorrentEntity,
    TorrentIdentity,
    ProgressUpdateRequest
} from '../types/model/models.ts'; // Adjust import path to your models file

export class TorrentService {

    /**
     * GET /torrents/
     * Fetches all torrents for the authenticated user.
     */
    public async getTorrents(): Promise<TorrentDTO[]> {
        const response = await api.get<TorrentDTO[]>('/torrents/');
        return response.data;
    }

    /**
     * POST /torrents/create
     * Uploads a new torrent file (multipart/form-data).
     */
    public async createTorrent(file: File): Promise<TorrentDTO> {
        const formData = new FormData();
        formData.append('torrent_file', file);

        const response = await api.post<TorrentDTO>('/torrents/create', formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        });
        return response.data;
    }

    /**
     * POST /torrents/add
     * Adds an existing torrent to the user's list by identity (JSON body).
     */
    public async addTorrent(identity: TorrentIdentity): Promise<void> {
        await api.post('/torrents/add', identity);
    }

    /**
     * GET /torrents/info?info_hash=...&creator_pub_key=...
     * Fetches a specific torrent entity by identity (Query params).
     */
    public async getTorrentByIdentity(identity: TorrentIdentity): Promise<TorrentEntity> {
        const params = new URLSearchParams({
            info_hash: identity.info_hash,
            creator_pub_key: identity.creator_pub_key
        });
        const response = await api.get<TorrentEntity>(`/torrents/info?${params.toString()}`);
        return response.data;
    }

    /**
     * POST /torrents/pause
     * Pauses a torrent (JSON body).
     */
    public async pauseTorrent(identity: TorrentIdentity): Promise<void> {
        await api.post('/torrents/pause', identity);
    }

    /**
     * POST /torrents/resume
     * Resumes a torrent (JSON body).
     */
    public async resumeTorrent(identity: TorrentIdentity): Promise<void> {
        await api.post('/torrents/resume', identity);
    }

    /**
     * POST /torrents/progress
     * Updates the download progress of a torrent (JSON body).
     */
    public async updateProgress(identity: TorrentIdentity, progress: number): Promise<void> {
        const req: ProgressUpdateRequest = { ...identity, progress };
        await api.post('/torrents/progress', req);
    }

    /**
     * DELETE /torrents/
     * Deletes a torrent from storage and database (JSON body).
     * Note: Axios requires the `data` property to send a body in a DELETE request.
     */
    public async deleteTorrent(identity: TorrentIdentity): Promise<void> {
        await api.delete('/torrents/', { data: identity });
    }

    /**
     * POST /torrents/sign
     * Signs a torrent using the user's DPKI keys (JSON body).
     */
    public async signTorrent(identity: TorrentIdentity): Promise<TorrentDTO> {
        const response = await api.post<TorrentDTO>('/torrents/sign', identity);
        return response.data;
    }
}

export const torrentService = new TorrentService();