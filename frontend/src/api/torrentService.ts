import { api } from './axiosClient';
import type {Torrent} from '../types/model/models.ts';
import type {AxiosResponse} from 'axios';

export class TorrentService {
    public static readonly responsibility = "Provide encapsulated methods to interact with Torrent related REST endpoints";

    public async getTorrents(): Promise<Torrent[]> {
        return this.tryGetTorrents();
    }

    private async tryGetTorrents(): Promise<Torrent[]> {
        try {
            return await this.executeGetTorrents();
        } catch (error) {
            return this.handleError(error as Error);
        }
    }

    private async executeGetTorrents(): Promise<Torrent[]> {
        const response: AxiosResponse<Torrent[]> = await api.get('/torrents');
        return response.data;
    }

    private handleError(error: Error): never {
        throw error;
    }
}

export const torrentService = new TorrentService();