import React, { memo } from 'react';
import type { SignatureDTO } from '@/types/model/models.ts';
import { SignatureRow } from './SignatureRow.tsx';
import { EmptySignaturesList } from './EmptySignaturesList.tsx';
import styles from './TorrentSignaturesList.module.css';

interface TorrentSignaturesListProps {
    signatures: SignatureDTO[];
}

const TorrentSignaturesListComponent: React.FC<TorrentSignaturesListProps> = ({ signatures }) => {
    if (signatures.length === 0) {
        return <EmptySignaturesList />;
    }

    return (
        <div className={styles.listContainer}>
            {signatures.map((sig: SignatureDTO): React.ReactNode => (
                <SignatureRow key={sig.signature_bytes} signature={sig} />
            ))}
        </div>
    );
};

export const TorrentSignaturesList = memo(TorrentSignaturesListComponent);