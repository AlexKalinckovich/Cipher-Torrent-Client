import React, { useState, useCallback, memo } from 'react';
import { CloseOutlined } from '@ant-design/icons';
import { message, Input } from 'antd';
import type { TorrentIdentity } from '@/types/model/models.ts';
import styles from './AddTorrentModal.module.css';

export interface AddTorrentModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSubmit: (identity: TorrentIdentity) => void;
}

const stopPropagation = (e: React.MouseEvent): void => {
    e.stopPropagation();
};

const AddTorrentModalComponent: React.FC<AddTorrentModalProps> = ({ isOpen, onClose, onSubmit }) => {
    const [infoHash, setInfoHash] = useState<string>('');
    const [creatorPubKey, setCreatorPubKey] = useState<string>('');

    const handleInfoHashChange = useCallback((e: React.ChangeEvent<HTMLInputElement>): void => {
        setInfoHash(e.target.value.trim());
    }, []);

    const handleCreatorPubKeyChange = useCallback((e: React.ChangeEvent<HTMLInputElement>): void => {
        setCreatorPubKey(e.target.value.trim());
    }, []);

    const handleSubmit = useCallback((): void => {
        if (!infoHash || !creatorPubKey) {
            void message.error('INFO HASH AND CREATOR PUBLIC KEY ARE REQUIRED');
            return;
        }

        onSubmit({
            info_hash: infoHash,
            creator_pub_key: creatorPubKey
        });

        // Reset form state
        setInfoHash('');
        setCreatorPubKey('');
        onClose();
    }, [infoHash, creatorPubKey, onSubmit, onClose]);

    if (!isOpen) {
        return null;
    }

    return (
        <div className={styles.overlay} onClick={onClose}>
            <div className={styles.modal} onClick={stopPropagation}>
                <div className={styles.header}>
                    <h2 className={styles.title}>ADD EXISTING TORRENT</h2>
                    <button type="button" className={styles.closeBtn} onClick={onClose}>
                        <CloseOutlined />
                    </button>
                </div>
                <div className={styles.content}>
                    <div className={styles.inputGroup}>
                        <span className={styles.label}>INFO HASH (HEX)</span>
                        <Input
                            className={styles.input}
                            placeholder="e.g. dd8255ecdc7ca55fb0bbf81323d87062db1f6d1c"
                            value={infoHash}
                            onChange={handleInfoHashChange}
                        />
                    </div>
                    <div className={styles.inputGroup}>
                        <span className={styles.label}>CREATOR PUBLIC KEY (BASE64)</span>
                        <Input
                            className={styles.input}
                            placeholder="e.g. r1GlehPWmP67rMrCeBIbYKmqvW1HZ5XvgzFpePTgukU="
                            value={creatorPubKey}
                            onChange={handleCreatorPubKeyChange}
                        />
                    </div>
                    <button type="button" className={styles.submitBtn} onClick={handleSubmit}>
                        ADD TO NETWORK QUEUE
                    </button>
                </div>
            </div>
        </div>
    );
};

export const AddTorrentModal = memo(AddTorrentModalComponent);