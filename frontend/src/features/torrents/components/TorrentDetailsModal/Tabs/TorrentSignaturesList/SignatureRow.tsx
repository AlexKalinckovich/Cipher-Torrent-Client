import React, { memo } from 'react';
import type { TorrentSignature } from '@/types/model/models.ts';
import styles from './TorrentSignaturesList.module.css';

interface SignatureRowProps {
    signature: TorrentSignature;
}

const getRowClass = (isValid: boolean): string => {
    if (isValid) {
        return styles.row;
    }
    return `${styles.row} ${styles.rowInvalid}`;
};

const getBadgeClass = (isValid: boolean): string => {
    if (isValid) {
        return `${styles.badge} ${styles.badgeValid}`;
    }
    return `${styles.badge} ${styles.badgeInvalid}`;
};

const formatDate = (isoString: string): string => {
    return new Date(isoString).toLocaleString();
};

const formatShortKey = (key: string): string => {
    if (key.length <= 20) {
        return key;
    }
    return `${key.substring(0, 12)}...${key.substring(key.length - 8)}`;
};

const SignatureRowComponent: React.FC<SignatureRowProps> = ({ signature }) => {
    return (
        <div className={getRowClass(signature.is_valid)}>
            <div className={styles.keyBlock}>
                <span className={styles.keyLabel}>SIGNER PUBLIC KEY</span>
                <span className={styles.keyValue}>{formatShortKey(signature.signer_public_key)}</span>
            </div>
            <div className={styles.dateBlock}>
                {formatDate(signature.signed_at)}
            </div>
            <div className={getBadgeClass(signature.is_valid)}>
                {signature.is_valid ? 'VALID' : 'INVALID'}
            </div>
        </div>
    );
};

const areSignatureRowEqual = (prevProps: SignatureRowProps, nextProps: SignatureRowProps): boolean => {
    return prevProps.signature.signature_bytes === nextProps.signature.signature_bytes;
};

export const SignatureRow = memo(SignatureRowComponent, areSignatureRowEqual);