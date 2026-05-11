import { useMutation, type UseMutationResult } from '@tanstack/react-query';
import { useNavigate } from 'react-router-dom';
import { AxiosError } from 'axios';
import type {
    AuthLoginRequest,
    AuthRegisterRequest,
    AuthResponse,
    ApiError
} from '../types/model/models.ts';
import { authService } from '../api/authService';
import { TokenRepository } from '../api/tokenRepository';

export const useLogin = (): UseMutationResult<AuthResponse, AxiosError<ApiError>, AuthLoginRequest> => {
    const navigate = useNavigate();
    const tokens = new TokenRepository();

    return useMutation<AuthResponse, AxiosError<ApiError>, AuthLoginRequest>({
        mutationFn: (request: AuthLoginRequest) => authService.login(request),
        onSuccess: (data: AuthResponse) => {
            tokens.setAccessToken(data.access_token);
            tokens.setRefreshToken(data.refresh_token);
            navigate('/dashboard');
        }
    });
};

export const useRegister = (): UseMutationResult<AuthResponse, AxiosError<ApiError>, AuthRegisterRequest> => {
    const navigate = useNavigate();
    const tokens = new TokenRepository();

    return useMutation<AuthResponse, AxiosError<ApiError>, AuthRegisterRequest>({
        mutationFn: (request: AuthRegisterRequest) => authService.register(request),
        onSuccess: (data: AuthResponse) => {
            tokens.setAccessToken(data.access_token);
            tokens.setRefreshToken(data.refresh_token);
            navigate('/dashboard');
        }
    });
};