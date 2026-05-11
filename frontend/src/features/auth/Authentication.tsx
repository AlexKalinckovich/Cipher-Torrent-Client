import React, { useState, useCallback, useEffect } from 'react';
import { Form, Input, message } from 'antd';
import type { FormInstance } from 'antd/es/form';
import { UserOutlined, LockOutlined, MailOutlined } from '@ant-design/icons';
import { useLogin, useRegister } from '../../hooks/useAuth';
import type { AuthLoginRequest, AuthRegisterRequest } from '../../types/model/models.ts';
import styles from './Authentication.module.css';

type AuthMode = 'login' | 'register';

interface AuthFieldError {
    name: string;
    errors: string[];
}

interface AuthErrorPayload {
    details?: Record<string, string[]>;
    message: string;
}

interface AuthErrorResponse {
    response?: { data?: AuthErrorPayload };
}

interface AuthModeToggleProps {
    mode: AuthMode;
    onSetMode: (mode: AuthMode) => void;
}

interface AuthFormFieldsProps {
    mode: AuthMode;
}

interface AuthSubmitButtonProps {
    isPending: boolean;
    mode: AuthMode;
}

interface UseAuthErrorsProps {
    form: FormInstance;
    mode: AuthMode;
    loginError: unknown;
    registerError: unknown;
}

const mapDetailsToFields = (details: Record<string, string[]>): AuthFieldError[] => {
    return Object.keys(details).map((key: string): AuthFieldError => ({
        name: key,
        errors: details[key]
    }));
};

const applyFieldErrors = (form: FormInstance, fields: AuthFieldError[]): void => {
    form.setFields(fields);
};

const showErrorMessage = (messageText: string): void => {
    void message.error(messageText);
};

const processAuthErrorDetails = (details: Record<string, string[]>, form: FormInstance): void => {
    applyFieldErrors(form, mapDetailsToFields(details));
};

const handleDetailsIfExists = (details: Record<string, string[]> | undefined, form: FormInstance): void => {
    if (details === undefined) return;
    processAuthErrorDetails(details, form);
};

const processValidAuthError = (data: AuthErrorPayload, form: FormInstance): void => {
    showErrorMessage(data.message);
    handleDetailsIfExists(data.details, form);
};

const processAuthError = (errorData: AuthErrorPayload | undefined, form: FormInstance): void => {
    if (errorData === undefined) return;
    processValidAuthError(errorData, form);
};

const extractErrorData = (error: AuthErrorResponse | undefined): AuthErrorPayload | undefined => {
    return error?.response?.data;
};

const getActiveError = (mode: AuthMode, loginError: unknown, registerError: unknown): AuthErrorResponse | undefined => {
    if (mode === 'login') return loginError as AuthErrorResponse | undefined;
    return registerError as AuthErrorResponse | undefined;
};

const useAuthErrorHandling = ({ form, mode, loginError, registerError }: UseAuthErrorsProps): void => {
    useEffect(() => {
        const activeError = getActiveError(mode, loginError, registerError);
        const errorData = extractErrorData(activeError);
        processAuthError(errorData, form);
    }, [form, mode, loginError, registerError]);
};

const executeLogin = (values: AuthLoginRequest, mutate: (vars: AuthLoginRequest) => void): void => {
    mutate(values);
};

const executeRegister = (values: AuthRegisterRequest, mutate: (vars: AuthRegisterRequest) => void): void => {
    mutate(values);
};

const handleAuthSubmit = (values: AuthLoginRequest | AuthRegisterRequest, mode: AuthMode, loginMutate: (vars: AuthLoginRequest) => void, registerMutate: (vars: AuthRegisterRequest) => void): void => {
    if (mode === 'login') {
        executeLogin(values as AuthLoginRequest, loginMutate);
    }else {
        executeRegister(values as AuthRegisterRequest, registerMutate);
    }
};

const AuthModeToggle: React.FC<AuthModeToggleProps> = ({ mode, onSetMode }) => {
    const handleLoginClick = useCallback(() => onSetMode('login'), [onSetMode]);
    const handleRegisterClick = useCallback(() => onSetMode('register'), [onSetMode]);

    return (
        <div className={styles.toggleWrapper}>
            <button type="button" className={mode === 'login' ? `${styles.toggleBtn} ${styles.toggleBtnActive}` : styles.toggleBtn} onClick={handleLoginClick}>Login</button>
            <button type="button" className={mode === 'register' ? `${styles.toggleBtn} ${styles.toggleBtnActive}` : styles.toggleBtn} onClick={handleRegisterClick}>Register</button>
        </div>
    );
};

const RegisterField: React.FC = () => {
    return (
        <Form.Item name="nickname" rules={[{ required: true }]}>
            <Input placeholder="Nickname" prefix={<UserOutlined />} className={styles.cyberInput} />
        </Form.Item>
    );
};

const AuthFormFields: React.FC<AuthFormFieldsProps> = ({ mode }) => {
    return (
        <>
            {mode === 'register' ? <RegisterField /> : null}
            <Form.Item name="email" rules={[{ required: true, type: 'email' }]}>
                <Input placeholder="Email" prefix={<MailOutlined />} className={styles.cyberInput} />
            </Form.Item>
            <Form.Item name="password" rules={[{ required: true }]}>
                <Input.Password placeholder="Password" prefix={<LockOutlined />} className={styles.cyberInput} />
            </Form.Item>
        </>
    );
};

const AuthSubmitButton: React.FC<AuthSubmitButtonProps> = ({ isPending, mode }) => {
    const buttonText: string = mode === 'login' ? 'Access System' : 'Create Identity';
    return (
        <button type="submit" className={styles.submitBtn} disabled={isPending}>
            {buttonText}
        </button>
    );
};

export const Authentication: React.FC = () => {
    const [mode, setMode] = useState<AuthMode>('login');
    const [form] = Form.useForm();
    const loginMutation = useLogin();
    const registerMutation = useRegister();

    useAuthErrorHandling({ form, mode, loginError: loginMutation.error, registerError: registerMutation.error });

    const handleSubmit = useCallback((values: AuthLoginRequest | AuthRegisterRequest) => {
        handleAuthSubmit(values, mode, loginMutation.mutate, registerMutation.mutate);
    }, [mode, loginMutation.mutate, registerMutation.mutate]);

    const handleModeChange = useCallback((newMode: AuthMode) => {
        setMode(newMode);
    }, []);

    const isPending: boolean = loginMutation.isPending || registerMutation.isPending;

    return (
        <div className={styles.pageContainer}>
            <div className={`${styles.backgroundBlob} ${styles.blobLeft}`} />
            <div className={`${styles.backgroundBlob} ${styles.blobRight}`} />
            <h1 className={styles.cyberTitle}>Welcome to <br /> Cipher-Torrent-Client!</h1>
            <div className={styles.authCard}>
                <AuthModeToggle mode={mode} onSetMode={handleModeChange} />
                <Form form={form} onFinish={handleSubmit} layout="vertical" requiredMark={false}>
                    <AuthFormFields mode={mode} />
                    <AuthSubmitButton isPending={isPending} mode={mode} />
                </Form>
            </div>
        </div>
    );
};