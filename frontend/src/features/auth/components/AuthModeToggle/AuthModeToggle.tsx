import React, { useCallback } from 'react';
import type { AuthMode } from '../../types/authTypes';
import styles from './AuthModeToggle.module.css';

interface AuthModeToggleProps {
    mode: AuthMode;
    onSetMode: (mode: AuthMode) => void;
}

export const AuthModeToggle: React.FC<AuthModeToggleProps> = ({ mode, onSetMode }) => {
    const handleLoginClick = useCallback(() => onSetMode('login'), [onSetMode]);
    const handleRegisterClick = useCallback(() => onSetMode('register'), [onSetMode]);

    return (
        <div className={styles.toggleWrapper}>
            <button
                type="button"
                className={mode === 'login' ? `${styles.toggleBtn} ${styles.toggleBtnActive}` : styles.toggleBtn}
                onClick={handleLoginClick}
            >
                Login
            </button>
            <button
                type="button"
                className={mode === 'register' ? `${styles.toggleBtn} ${styles.toggleBtnActive}` : styles.toggleBtn}
                onClick={handleRegisterClick}
            >
                Register
            </button>
        </div>
    );
};
