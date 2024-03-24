import {createContext, useState, ReactNode, useContext} from 'react';
import {redirect} from "@tanstack/react-router";

export type AuthAccessToken = {
    accessToken: string;
    expires_at: number;
}

export type AuthRefreshToken = {
    refreshToken: string;
    expires_at: number;
}

type AuthContextType = {
    saveCredentials: ({accessToken, refreshToken}: {accessToken: AuthAccessToken, refreshToken: AuthRefreshToken}) => void;
    logout: () => void;
    isAuthenticated: () => boolean;
}

const AuthContext = createContext<AuthContextType>({
    saveCredentials: () => {},
    logout: () => {},
    isAuthenticated: () => false,
});

const AuthProvider = ({ children }: { children: ReactNode }) => {
    const [accessToken, setAccessToken] = useState<AuthAccessToken>(() => {
        const storage = localStorage.getItem('auth');
        if (storage) {
            const { accessToken } = JSON.parse(storage);
            return accessToken;
        }

        return {accessToken: '', expiresAt: 0};
    });

    const [refreshToken, setRefreshToken] = useState<AuthRefreshToken>(() => {
        const storage = localStorage.getItem('auth');
        if (storage) {
            const { refreshToken } = JSON.parse(storage);
            return refreshToken;
        }

        return {refreshToken: '', expiresAt: 0};
    });

    const saveCredentials = ({accessToken, refreshToken}: {accessToken: AuthAccessToken, refreshToken: AuthRefreshToken}) => {
        setAccessToken(accessToken);
        setRefreshToken(refreshToken);

        localStorage.setItem('auth', JSON.stringify({ accessToken, refreshToken }));
    }

    const logout = () => {
        setAccessToken({accessToken: '', expires_at: 0});
        setRefreshToken({refreshToken: '', expires_at: 0});

        localStorage.removeItem('auth');
    }

    const isAuthenticated = () => {
        // todo: check the validity of the access token from the server
        // and/or refresh the token.
        const now = new Date();
        return new Date(accessToken.expires_at) > now && new Date(refreshToken.expires_at) > now;
    }

    return (
        <AuthContext.Provider value={{ saveCredentials, logout, isAuthenticated }}>
            {children}
        </AuthContext.Provider>
    );
}

const useAuth = () => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error('useAuth must be used within an AuthProvider');
    }

    return context;
}

const authMiddleware = (auth: AuthContextType) => {
    if (!auth.isAuthenticated()) {
        throw redirect({to: '/authenticate'})
    }
}

// eslint-disable-next-line react-refresh/only-export-components
export { AuthProvider, useAuth, type AuthContextType, authMiddleware }