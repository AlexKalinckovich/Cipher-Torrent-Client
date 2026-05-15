import React, { useState } from 'react';
import type { AuthMode } from './types/authTypes';
import { AuthHeader } from './components/AuthHeader/AuthHeader';
import { AuthModeToggle } from './components/AuthModeToggle/AuthModeToggle';
import { LoginForm } from './components/LoginForm/LoginForm';
import { RegisterForm } from './components/RegisterForm/RegisterForm';
import styles from './common.module.css';

const renderForm = (mode: AuthMode): React.ReactElement => {
    if (mode === 'login') {
        return <LoginForm />;
    }
    return <RegisterForm />;
};

export const Authentication: React.FC = () => {
    const [mode, setMode] = useState<AuthMode>('login');

    return (
        <div className={styles.pageContainer}>
            <AuthHeader />
            <div className={styles.authCard}>
                <AuthModeToggle mode={mode} onSetMode={setMode} />
                {renderForm(mode)}
            </div>
        </div>
    );
};
