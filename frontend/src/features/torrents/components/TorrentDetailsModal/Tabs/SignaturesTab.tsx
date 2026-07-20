import React, { useCallback, useState } from 'react';
import type { TorrentDTO, TorrentIdentity } from '@/types/model/models.ts';
import { TorrentSignaturesList } from './TorrentSignaturesList/TorrentSignaturesList';
import { OptionalSignButton } from './TorrentSignaturesList/SignTorrentButton';
import { useSignTorrent } from '@/hooks/useTorrents.ts';
import { message } from 'antd';
import styles from './Tabs.module.css';

interface TabProps {
    torrent: TorrentDTO;
}

const getSignedClass = (isSigned: boolean): string => {
    if (isSigned) {
        return `${styles.statusBadge} ${styles.badgeGreen}`;
    }
    return `${styles.statusBadge} ${styles.badgeRed}`;
};

export const SignaturesTab: React.FC<TabProps> = ({ torrent }) => {
    const { mutate: signTorrent } = useSignTorrent();
    const [isSignedByMe, setIsSignedByMe] = useState<boolean>(false);

    const identity: TorrentIdentity = {
        info_hash: torrent.info_hash,
        creator_pub_key: torrent.creator_public_key // Assuming you add this to your DTO or extract it from context
    };

    const handleSignSuccess = useCallback((): void => {
        void message.success('TORRENT CRYPTOGRAPHICALLY SIGNED');
        setIsSignedByMe(true);
    }, []);

    const handleSignError = useCallback((error: Error): void => {
        void message.error(`SIGNATURE FAILED: ${error.message}`);
    }, []);

    const handleSignClick = useCallback((): void => {
        signTorrent(identity, {
            onSuccess: handleSignSuccess,
            onError: handleSignError
        });
    }, [signTorrent, identity, handleSignSuccess, handleSignError]);

    return (
        <div className={styles.tabContainer}>
            <div className={styles.sigRow}>
                <span className={styles.sigLabel}>Network Signatures</span>
                <span className={styles.sigValue}>{torrent.signatures.length}</span>
            </div>
            <div className={styles.sigRow}>
                <span className={styles.sigLabel}>Local Cryptographic Signature</span>
                <span className={getSignedClass(isSignedByMe)}>
                    {isSignedByMe ? 'VERIFIED' : 'UNSIGNED'}
                </span>
            </div>
            <OptionalSignButton
                isSignedByMe={isSignedByMe}
                onSignClick={handleSignClick}
            />
            <TorrentSignaturesList signatures={torrent.signatures} />
        </div>
    );
};