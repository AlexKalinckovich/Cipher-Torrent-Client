import { api } from './axiosClient.ts';
import type { AxiosResponse } from 'axios';
import type {
    Torrent,
    TorrentAddRequest,
    TorrentAddResponse,
    SignatureCreateResponse
} from '../types/model/models.ts';

export class TorrentService {
    public async getTorrents(): Promise<Torrent[]> {
        return this.tryGetTorrents();
    }

    public async addTorrent(request: TorrentAddRequest): Promise<TorrentAddResponse> {
        return this.tryAddTorrent(request);
    }

    public async signTorrent(infoHash: string): Promise<SignatureCreateResponse> {
        return this.trySignTorrent(infoHash);
    }

    private async tryGetTorrents(): Promise<Torrent[]> {
        try {
            return await this.executeGetTorrents();
        } catch (error) {
            return this.handleError(error as Error);
        }
    }

    private async tryAddTorrent(request: TorrentAddRequest): Promise<TorrentAddResponse> {
        try {
            return await this.executeAddTorrent(request);
        } catch (error) {
            return this.handleError(error as Error);
        }
    }

    private async trySignTorrent(infoHash: string): Promise<SignatureCreateResponse> {
        try {
            return await this.executeSignTorrent(infoHash);
        } catch (error) {
            return this.handleError(error as Error);
        }
    }

    private async executeGetTorrents(): Promise<Torrent[]> {
        const response: AxiosResponse<Torrent[]> = await api.get('/torrents');
        return response.data;
    }

    private async executeAddTorrent(request: TorrentAddRequest): Promise<TorrentAddResponse> {
        const response: AxiosResponse<TorrentAddResponse> = await api.post('/torrents', request);
        return response.data;
    }

    private async executeSignTorrent(infoHash: string): Promise<SignatureCreateResponse> {
        const response: AxiosResponse<SignatureCreateResponse> = await api.post(`/torrents/${infoHash}/sign`);
        return response.data;
    }

    private handleError(error: Error): never {
        throw error;
    }
}

export const torrentService = new TorrentService();