import type {UserStats} from './userStats';

export type UserRole = 'user' | 'admin';

export interface UserFull {
    id: number;
    email: string;
    public_key: string;
    nickname: string;
    created_at: string;
    stats?: UserStats;
    role?: UserRole;
}