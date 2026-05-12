import type {AxiosRequestConfig} from 'axios';
import {type AuthLoginRequest, type AuthRegisterRequest, type AuthResponse, type ApiError } from '../types/model/models.ts';

export class AuthMockHandlers {
    public static readonly responsibility = "Intercept authentication requests and provide scenario-based mocked ApiError or AuthResponse";

    public handleLogin(config: AxiosRequestConfig): [number, AuthResponse | ApiError] {
        const body: AuthLoginRequest = JSON.parse(config.data);
        console.log(body)
        return this.validateLoginCredentials(body);
    }

    private validateLoginCredentials(body: AuthLoginRequest): [number, AuthResponse | ApiError] {
        console.log(`received request ${body.email}, ${body.password}`)
        if (body.email === 'admin@cipher.com' && body.password === 'Cipher123!') {
            return this.createSuccessResponse();
        }
        return this.createErrorResponse(401, 'Invalid credentials or user does not exist', 'AUTH_INVALID_CREDENTIALS');
    }

    public handleRegister(config: AxiosRequestConfig): [number, AuthResponse | ApiError] {
        const body: AuthRegisterRequest = JSON.parse(config.data);
        return this.validateRegistrationEmail(body);
    }

    private validateRegistrationEmail(body: AuthRegisterRequest): [number, AuthResponse | ApiError] {
        if (body.email === 'admin@cipher.com') {
            return this.createErrorResponse(400, 'User with this email already exists', 'AUTH_EMAIL_EXISTS');
        }
        return this.validateRegistrationPassword(body);
    }

    private validateRegistrationPassword(body: AuthRegisterRequest): [number, AuthResponse | ApiError] {
        if (body.password.length < 8) {
            return this.createPasswordErrorResponse();
        }
        return this.createSuccessResponse();
    }

    private createPasswordErrorResponse(): [number, ApiError] {
        const details: Record<string, string[]> = {
            password: ['Password must be at least 8 characters long', 'Must contain numbers and symbols']
        };
        return this.createErrorResponse(400, 'Validation failed', 'VALIDATION_ERROR', details);
    }

    private createSuccessResponse(): [number, AuthResponse] {
        return [200, {
            access_token: 'mock_jwt_access_token_777',
            refresh_token: 'mock_jwt_refresh_token_999',
            user: {
                id: 1,
                email: 'admin@cipher.com',
                public_key: 'Ed25519_PubKey_Cipher',
                nickname: 'CipherAdmin',
                created_at: '2026-05-11T10:00:00Z'
            }
        }];
    }

    private createErrorResponse(status: number, message: string, errorCode: string, details?: Record<string, string[]>): [number, ApiError] {
        return [status, { status, message, errorCode, details }];
    }
}