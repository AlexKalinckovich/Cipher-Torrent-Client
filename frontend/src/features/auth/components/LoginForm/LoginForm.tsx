import React, { useEffect } from 'react';
import { Form, Input } from 'antd';
import { LockOutlined, MailOutlined } from '@ant-design/icons';
import { useLogin } from '@/hooks/useAuth';
import type { AuthLoginRequest } from '@/types/model/models.ts';
import { applyServerErrors } from '../../utils/authErrorUtils';
import styles from '../shared/FormStyles.module.css';

export const LoginForm: React.FC = () => {
    const [form] = Form.useForm();
    const { mutate, isPending, error } = useLogin();

    useEffect(() => {
        applyServerErrors(error, form);
    }, [error, form]);

    const onFinish = (values: AuthLoginRequest): void => {
        mutate(values);
    };

    return (
        <Form form={form} onFinish={onFinish} layout="vertical" requiredMark={false}>
            <Form.Item name="email" rules={[{ required: true, type: 'email' }]}>
                <Input placeholder="Email" prefix={<MailOutlined />} className={styles.cyberInput} />
            </Form.Item>

            <Form.Item name="password" rules={[{ required: true }]}>
                <Input.Password placeholder="Password" prefix={<LockOutlined />} className={styles.cyberInput} />
            </Form.Item>

            <button type="submit" className={styles.submitBtn} disabled={isPending}>
                Access System
            </button>
        </Form>
    );
};
