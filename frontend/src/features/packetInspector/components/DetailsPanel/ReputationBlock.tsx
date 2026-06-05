import React from 'react';
import type { ReputationReceipt } from '@/types/model/Reputation/reputationReceipt';
import styles from './ReputationBlock.module.css';

interface ReputationBlockProps {
    receipt: ReputationReceipt;
}

interface OptionalReputationBlockProps {
    payload?: ReputationReceipt;
}

const formatShortId = (id: string): string => {
    return `${id.substring(0, 8)}...${id.substring(id.length - 8)}`;
};

export const ReputationBlock: React.FC<ReputationBlockProps> = ({ receipt }) => {
    return (
        <div className={styles.reputationWrapper}>
            <div className={styles.glitchHeader}>REPUTATION RECEIPT</div>
            <div className={styles.grid}>
                <div className={styles.cell}>
                    <span className={styles.label}>FROM PEER ID</span>
                    <span className={styles.value}>{formatShortId(receipt.from_peer_id)}</span>
                </div>
                <div className={styles.cell}>
                    <span className={styles.label}>TO PEER ID</span>
                    <span className={styles.value}>{formatShortId(receipt.to_peer_id)}</span>
                </div>
                <div className={styles.cell}>
                    <span className={styles.label}>PIECE INDEX</span>
                    <span className={styles.value}>{receipt.piece_index}</span>
                </div>
                <div className={styles.cell}>
                    <span className={styles.label}>TRANSMITTED BYTES</span>
                    <span className={styles.value}>{receipt.byte_count}</span>
                </div>
            </div>
            <div className={styles.sigBlock}>
                <span className={styles.sigLabel}>ED25519 CRYPTOGRAPHIC SIGNATURE</span>
                <span className={styles.sigValue}>{receipt.signature}</span>
            </div>
        </div>
    );
};

export const OptionalReputationBlock: React.FC<OptionalReputationBlockProps> = ({ payload }) => {
    if (!payload) {
        return null;
    }
    return <ReputationBlock receipt={payload} />;
};