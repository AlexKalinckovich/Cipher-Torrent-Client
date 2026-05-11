import React, { useState, useEffect } from 'react';
import { Input, Form, message } from 'antd';
import { UserOutlined, LockOutlined, MailOutlined } from '@ant-design/icons';
import { useLogin, useRegister } from '../../hooks/useAuth';
import type { AuthLoginRequest, AuthRegisterRequest } from '../../types/model/models.ts';
import styles from './Authentication.module.css';

type AuthMode = 'login' | 'register';

const mapDetailsToFields = (details: Record<string, string[]>) => {
    return Object.keys(details).map(key => ({
        name: key,
        errors: details[key]
    }));
};

export const Authentication: React.FC = () => {
    const [mode, setMode] = useState<AuthMode>('login');
    const [form] = Form.useForm();

    const loginMutation = useLogin();
    const registerMutation = useRegister();

    const handleAction = (values: AuthLoginRequest | AuthRegisterRequest): void => {
        if (mode === 'login') {
            loginMutation.mutate(values as AuthLoginRequest);
        } else {
            registerMutation.mutate(values as AuthRegisterRequest);
        }
    };

    useEffect(() => {
        const error = (loginMutation.error || registerMutation.error)?.response?.data;
        if (error) {
            if (error.details) {
                form.setFields(mapDetailsToFields(error.details));
            }
            void message.error(error.message);
        }
    }, [form, loginMutation.error, registerMutation.error]);

    return (
        <div className={styles.pageContainer}>
            <div className={`${styles.backgroundBlob} ${styles.blobLeft}`} />
            <div className={`${styles.backgroundBlob} ${styles.blobRight}`} />

            <h1 className={styles.cyberTitle}>
                Welcome to <br /> Cipher-Torrent-Client!
            </h1>

            <div className={styles.authCard}>
                <div className={styles.toggleWrapper}>
                    <button
                        type="button"
                        className={mode === 'login' ? `${styles.toggleBtn} ${styles.toggleBtnActive}` : styles.toggleBtn}
                        onClick={() => setMode('login')}
                    >
                        Login
                    </button>
                    <button
                        type="button"
                        className={mode === 'register' ? `${styles.toggleBtn} ${styles.toggleBtnActive}` : styles.toggleBtn}
                        onClick={() => setMode('register')}
                    >
                        Register
                    </button>
                </div>

                <Form
                    form={form}
                    onFinish={handleAction}
                    layout="vertical"
                    requiredMark={false}
                >
                    {mode === 'register' && (
                        <Form.Item name="nickname" rules={[{ required: true }]}>
                            <Input placeholder="Nickname" prefix={<UserOutlined />} className={styles.cyberInput} />
                        </Form.Item>
                    )}

                    <Form.Item name="email" rules={[{ required: true, type: 'email' }]}>
                        <Input placeholder="Email" prefix={<MailOutlined />} className={styles.cyberInput} />
                    </Form.Item>

                    <Form.Item name="password" rules={[{ required: true }]}>
                        <Input.Password placeholder="Password" prefix={<LockOutlined />} className={styles.cyberInput} />
                    </Form.Item>

                    <button
                        type="submit"
                        className={styles.submitBtn}
                        disabled={loginMutation.isPending || registerMutation.isPending}
                    >
                        {mode === 'login' ? 'Access System' : 'Create Identity'}
                    </button>
                </Form>
            </div>
        </div>
    );
};