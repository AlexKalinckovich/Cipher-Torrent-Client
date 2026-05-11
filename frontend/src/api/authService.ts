import { api } from './axiosClient';
import type {
    AuthLoginRequest,
    AuthRegisterRequest,
    AuthResponse
} from '../types/model/models.ts';
import type { AxiosResponse } from 'axios';

export class AuthService {
    public static readonly responsibility = "Handle authentication network requests and data extraction";

    public async login(request: AuthLoginRequest): Promise<AuthResponse> {
        return this.tryLogin(request);
    }

    private async tryLogin(request: AuthLoginRequest): Promise<AuthResponse> {
        try {
            return await this.executeLogin(request);
        } catch (error) {
            return this.handleAuthError(error as Error);
        }
    }

    private async executeLogin(request: AuthLoginRequest): Promise<AuthResponse> {
        const response: AxiosResponse<AuthResponse> = await api.post('/auth/login', request);
        return response.data;
    }

    public async register(request: AuthRegisterRequest): Promise<AuthResponse> {
        return this.tryRegister(request);
    }

    private async tryRegister(request: AuthRegisterRequest): Promise<AuthResponse> {
        try {
            return await this.executeRegister(request);
        } catch (error) {
            return this.handleAuthError(error as Error);
        }
    }

    private async executeRegister(request: AuthRegisterRequest): Promise<AuthResponse> {
        const response: AxiosResponse<AuthResponse> = await api.post('/auth/register', request);
        return response.data;
    }

    private handleAuthError(error: Error): never {
        throw error;
    }
}

export const authService = new AuthService();