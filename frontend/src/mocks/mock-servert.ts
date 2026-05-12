import MockAdapter from 'axios-mock-adapter';
import type { AxiosInstance } from 'axios';
import { mockTorrents } from './torrents.mock';
import { AuthMockHandlers } from './auth.mock.ts';

export class MockServer {
    public static readonly responsibility = "Intercept HTTP requests and return mocked responses for development";
    private readonly mock: MockAdapter;
    private readonly handler: AuthMockHandlers;

    constructor(axiosInstance: AxiosInstance) {
        this.mock = new MockAdapter(axiosInstance, { delayResponse: 800 });
        this.handler = new AuthMockHandlers();
        this.initializeRoutes();
    }

    private initializeRoutes(): void {
        this.mock.onGet('/torrents').reply(200, mockTorrents);
        this.mock.onPost('/auth/login').reply((config) => this.handler.handleLogin(config));
        this.mock.onPost('/auth/register').reply((config) => this.handler.handleRegister(config));
    }
}