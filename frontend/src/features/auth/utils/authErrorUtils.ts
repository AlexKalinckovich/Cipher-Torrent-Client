import type { FormInstance } from 'antd/es/form';
import { message } from 'antd';
import type { AuthFieldError, AuthErrorPayload } from '../types/authTypes';

const mapDetailsToFields = (details: Record<string, string[]>): AuthFieldError[] => {
    return Object.keys(details).map((key: string): AuthFieldError => ({
        name: key,
        errors: details[key]
    }));
};

const processAuthErrorDetails = (details: Record<string, string[]>, form: FormInstance): void => {
    form.setFields(mapDetailsToFields(details));
};

const handleDetailsIfExists = (details: Record<string, string[]> | undefined, form: FormInstance): void => {
    if (details) {
        processAuthErrorDetails(details, form);
    }
};

const extractErrorData = (error: unknown): AuthErrorPayload | undefined => {
    const typedError = error as { response?: { data?: AuthErrorPayload } };
    return typedError?.response?.data;
};

export const applyServerErrors = (error: unknown, form: FormInstance): void => {
    const errorData: AuthErrorPayload | undefined = extractErrorData(error);
    if (!errorData) {
        return;
    }
    void message.error(errorData.message);
    handleDetailsIfExists(errorData.details, form);
};
