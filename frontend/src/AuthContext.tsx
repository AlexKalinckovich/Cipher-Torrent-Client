import React, { createContext, useContext, useState, useEffect, type ReactNode } from 'react';
import { authService } from '@/api/authService.ts';
import { TokenRepository } from '@/api/tokenRepository';
import type { UserFull, AuthLoginRequest, AuthRegisterRequest, AuthResponse } from '@/types/model/models.ts';

interface AuthContextType {
    user: UserFull | null;
    isAuthenticated: boolean;
    isLoading: boolean;
    login: (req: AuthLoginRequest) => Promise<void>;
    register: (req: AuthRegisterRequest) => Promise<void>;
    logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);
const tokenRepo = new TokenRepository();

export const AuthProvider: React.FC<{ children: ReactNode }> = ({ children }) => {
    
    
    const [user, setUser] = useState<UserFull | null>(tokenRepo.getUser());

    
    const [isLoading, setIsLoading] = useState(!user);

    useEffect(() => {
        const initializeAuth = async () => {
            const token = tokenRepo.getAccessToken();
            const cachedUser = tokenRepo.getUser();

            if (token && !cachedUser) {
                try {
                    const profile: UserFull = await authService.getProfile();
                    tokenRepo.setUser(profile); 
                    setUser(profile);
                } catch (error) {
                    tokenRepo.clear(); 
                }
            }
            setIsLoading(false);
        };
        initializeAuth();
    }, []);

    const login = async (req: AuthLoginRequest) => {
        const authResponse: AuthResponse = await authService.login(req);

        
        tokenRepo.setAccessToken(authResponse.access_token);
        tokenRepo.setRefreshToken(authResponse.refresh_token);

        
        
        const profile: UserFull = authResponse.user;
        tokenRepo.setUser(profile);
        setUser(profile);
    };

    const register = async (req: AuthRegisterRequest) => {
        const authResponse: AuthResponse = await authService.register(req);

        tokenRepo.setAccessToken(authResponse.access_token);
        tokenRepo.setRefreshToken(authResponse.refresh_token);

        const profile = authResponse.user as unknown as UserFull;
        tokenRepo.setUser(profile);
        setUser(profile);
    };

    const logout = () => {
        tokenRepo.clear(); 
        setUser(null);
    };

    return (
        <AuthContext.Provider value={{
            user,
            isAuthenticated: !!user,
            isLoading,
            login,
            register,
            logout
        }}>
            {children}
        </AuthContext.Provider>
    );
};

export const useAuth = (): AuthContextType => {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
};