export type AuthMode = 'login' | 'register';

export interface AuthFieldError {
    name: string;
    errors: string[];
}

export interface AuthErrorPayload {
    details?: Record<string, string[]>;
    message: string;
}

export interface AuthErrorResponse {
    response?: { data?: AuthErrorPayload };
}
