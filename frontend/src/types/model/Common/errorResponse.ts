export interface ApiError {
    status: number;
    message: string;
    errorCode: string;
    details?: Record<string, string[]>;
}