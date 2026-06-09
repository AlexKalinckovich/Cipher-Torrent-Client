import MockAdapter from 'axios-mock-adapter';
import type { AxiosInstance } from 'axios';
import { mockTorrents } from './torrents.mock.ts';
import { AuthMockHandlers } from './auth.mock.ts';
import type { TorrentAddResponse, SignatureCreateResponse } from '../types/model/models.ts';

export class MockServer {
    private readonly mock: MockAdapter;
    private readonly handler: AuthMockHandlers;

    constructor(axiosInstance: AxiosInstance) {
        this.mock = new MockAdapter(axiosInstance, { delayResponse: 800 });
        this.handler = new AuthMockHandlers();
        this.initializeRoutes();
    }

    private initializeRoutes(): void {
        this.mock.onGet('/torrents').reply(200, mockTorrents);
        this.mock.onPost('/torrents').reply(this.handleAddTorrent.bind(this));
        this.mock.onPost(/\/torrents\/[A-Za-z0-9]+\/sign/).reply(this.handleSignTorrent.bind(this));

        this.mock.onPost('/auth/login').reply((config) => this.handler.handleLogin(config));
        this.mock.onPost('/auth/register').reply((config) => this.handler.handleRegister(config));
    }

    private handleAddTorrent(): [number, TorrentAddResponse] {
        const response: TorrentAddResponse = {
            status: 'success',
            info_hash: `NEW_${Math.random().toString(16).substring(2, 10).toUpperCase()}`
        };
        return [200, response];
    }

    private handleSignTorrent(): [number, SignatureCreateResponse] {
        const response: SignatureCreateResponse = {
            signature_id: Math.floor(Math.random() * 1000000),
            signed_at: new Date().toISOString()
        };
        return [200, response];
    }
}