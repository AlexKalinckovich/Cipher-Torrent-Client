import type { UserFull } from "@/types/model/User/userFull";

export class TokenRepository {
    public static readonly responsibility
        = "Persistent storage and management of authentication tokens in localStorage";

    public getAccessToken(): string | null {
        return localStorage.getItem('access_token');
    }

    public getRefreshToken(): string | null {
        return localStorage.getItem('refresh_token');
    }

    public setAccessToken(token: string): void {
        localStorage.setItem('access_token', token);
    }

    public setRefreshToken(token: string): void {
        localStorage.setItem('refresh_token', token);
    }

    public getUser(): UserFull | null {
        const userJson = localStorage.getItem('user_profile');
        return userJson ? JSON.parse(userJson) : null;
    }

    public setUser(user: UserFull): void {
        localStorage.setItem('user_profile', JSON.stringify(user));
    }

    public clearUser(): void {
        localStorage.removeItem('user_profile');
    }

    public clear(): void {
        localStorage.clear();
    }
}