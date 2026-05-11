import MockAdapter from 'axios-mock-adapter';
import type {AxiosInstance} from 'axios';
import { mockTorrents } from './torrents.mock';

export class MockServer {
    public static readonly responsibility = "Intercept HTTP requests and return mocked responses for development";
    private readonly mock: MockAdapter;

    constructor(axiosInstance: AxiosInstance) {
        this.mock = new MockAdapter(axiosInstance, { delayResponse: 800 });
        this.initializeRoutes();
    }

    private initializeRoutes(): void {
        this.mock.onGet('/torrents').reply(200, mockTorrents);
    }
}