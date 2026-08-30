import React, { memo } from 'react';
import { useNavigate } from 'react-router-dom';
import type { SignatureDTO } from '@/types/model/models.ts';
import styles from './TorrentSignaturesList.module.css';

interface SignatureRowProps {
    signature: SignatureDTO;
}

const formatDate = (timestamp: number): string => {
    return new Date(timestamp * 1000).toLocaleString();
};

const formatShortKey = (key: string): string => {
    if (key.length <= 20) {
        return key;
    }
    return `${key.substring(0, 12)}...${key.substring(key.length - 8)}`;
};

const SignatureRowComponent: React.FC<SignatureRowProps> = ({ signature }) => {
    const navigate = useNavigate();

    const handleViewProfile = (): void => {
        navigate(`/profile/${encodeURIComponent(signature.signer_public_key)}`);
    };

    return (
        <div className={styles.row}>
            <div className={styles.keyBlock}>
                <span className={styles.keyLabel}>SIGNER PUBLIC KEY</span>
                <span
                    className={styles.keyValue}
                    onClick={handleViewProfile}
                    style={{ cursor: 'pointer', textDecoration: 'underline' }}
                    title="View user profile"
                >
                    {formatShortKey(signature.signer_public_key)}
                </span>
            </div>
            <div className={styles.dateBlock}>
                {formatDate(signature.timestamp)}
            </div>
            <div className={styles.badge}>
                VALID
            </div>
        </div>
    );
};

const areSignatureRowEqual = (prevProps: SignatureRowProps, nextProps: SignatureRowProps): boolean => {
    return prevProps.signature.signature_bytes === nextProps.signature.signature_bytes;
};

export const SignatureRow = memo(SignatureRowComponent, areSignatureRowEqual);