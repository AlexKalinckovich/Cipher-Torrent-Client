import axios, {
    type AxiosInstance,
    type InternalAxiosRequestConfig,
    AxiosError,
    type AxiosResponse
} from 'axios';
import { TokenRepository } from './tokenRepository';
import type { AuthResponse, AuthRefreshRequest } from '../types/model/models.ts';

export class HttpClient {
    public static readonly responsibility = "Configured Axios instance with automated authentication headers and token refresh logic";

    private readonly instance: AxiosInstance;
    private readonly tokens: TokenRepository;

    private readonly publicEndpoints: string[] = [
        '/auth/login',
        '/auth/register',
        '/auth/refresh'
    ];

    constructor() {
        this.tokens = new TokenRepository();
        this.instance = axios.create({
            baseURL: 'http://localhost:8080/api/v1',
            headers: { 'Content-Type': 'application/json' }
        });
        this.initializeInterceptors();
    }

    private initializeInterceptors(): void {
        this.instance.interceptors.request.use(this.injectToken.bind(this));
        this.instance.interceptors.response.use(
            (response) => response,
            this.handleUnauthorized.bind(this)
        );
    }

    // 2. Helper to check if the request is for a public endpoint
    private isPublicEndpoint(url: string | undefined): boolean {
        return url ? this.publicEndpoints.some((endpoint: string) => url.includes(endpoint)) : false;
    }

    private injectToken(config: InternalAxiosRequestConfig): InternalAxiosRequestConfig {
        if (this.isPublicEndpoint(config.url)) {
            console.log("public confing")
            return config;
        }

        const token: string | null = this.tokens.getAccessToken();
        if (token) {
            config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
    }

    private async handleUnauthorized(error: AxiosError): Promise<AxiosResponse | unknown> {
        if (error.response?.status === 401) {
            return this.tryTokenRefresh(error.config as InternalAxiosRequestConfig);
        }
        return Promise.reject(error);
    }

    private async tryTokenRefresh(originalConfig: InternalAxiosRequestConfig): Promise<AxiosResponse | unknown> {
        try {
            console.log("Trying token refresh")
            return await this.executeRefreshSequence(originalConfig);
        } catch (refreshError) {
            this.tokens.clear();
            return Promise.reject(refreshError);
        }
    }

    private async executeRefreshSequence(config: InternalAxiosRequestConfig): Promise<AxiosResponse> {
        const refreshRequest: AuthRefreshRequest = {
            refresh_token: this.tokens.getRefreshToken() || ""
        };

        const response = await axios.post<AuthResponse>(
            `${this.instance.defaults.baseURL}/auth/refresh`,
            refreshRequest
        );

        this.updateStoredTokens(response.data);
        return this.instance(config);
    }

    private updateStoredTokens(data: AuthResponse): void {
        this.tokens.setAccessToken(data.access_token);
        this.tokens.setRefreshToken(data.refresh_token);
    }

    public getInstance(): AxiosInstance {
        return this.instance;
    }
}

export const api: AxiosInstance = new HttpClient().getInstance();

const isDev: boolean = import.meta.env ? import.meta.env.DEV : false;
if (isDev) {
    //new MockServer(api);
}