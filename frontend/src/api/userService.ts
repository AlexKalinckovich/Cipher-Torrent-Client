import { api } from './axiosClient';
import type { UserFull } from '../types/model/models.ts';
import type { AxiosResponse } from 'axios';

export class UserService {
    public static readonly responsibility = "Handle user-related network requests and data extraction";

    public async getProfile(): Promise<UserFull> {
        const response: AxiosResponse<UserFull> = await api.get('/users/me');
        return response.data;
    }

    public async getUserById(id: number): Promise<UserFull> {
        const response: AxiosResponse<UserFull> = await api.get(`/users/${id}`);
        return response.data;
    }

    public async getUserByPublicKey(publicKey: string): Promise<UserFull> {
        const encodedKey = encodeURIComponent(publicKey);
        const response: AxiosResponse<UserFull> = await api.get(`/users/public-key?pub_key=${encodedKey}`);
        return response.data;
    }
}

export const userService = new UserService();
