import React, { useCallback, useState } from 'react';
import type { Torrent, TorrentSignature } from '@/types/model/models.ts';
import { useTorrentSignaturesMock } from '@/hooks/useTorrentSignatureMock';
import { TorrentSignaturesList } from '@/features/torrents/components/TorrentDetailsModal/Tabs/TorrentSignaturesList/TorrentSignaturesList';
import { OptionalSignButton } from './TorrentSignaturesList/SignTorrentButton';
import styles from './Tabs.module.css';

interface TabProps {
    torrent: Torrent;
}

const getSignedClass = (isSigned: boolean): string => {
    if (isSigned) {
        return `${styles.statusBadge} ${styles.badgeGreen}`;
    }
    return `${styles.statusBadge} ${styles.badgeRed}`;
};

export const SignaturesTab: React.FC<TabProps> = ({ torrent } : TabProps) => {
    const initialSignatures = useTorrentSignaturesMock(torrent.info_hash);
    
    const [addedSignatures, setAddedSignatures] = useState<TorrentSignature[]>([]);
    
    const [isSignedByMe, setIsSignedByMe] = useState<boolean>(torrent.is_signed_by_me ?? false);

    
    const allSignatures = [...addedSignatures, ...initialSignatures];

    const handleNewSignature = useCallback((newSig: TorrentSignature): void => {
        setAddedSignatures((prev) => [newSig, ...prev]);
        setIsSignedByMe(true);
    }, []);

    return (
        <div className={styles.tabContainer}>
            <div className={styles.sigRow}>
                <span className={styles.sigLabel}>Network Signatures</span>
                <span className={styles.sigValue}>{allSignatures.length}</span>
            </div>
            <div className={styles.sigRow}>
                <span className={styles.sigLabel}>Local Cryptographic Signature</span>
                <span className={getSignedClass(isSignedByMe)}>
                    {isSignedByMe ? 'VERIFIED' : 'UNSIGNED'}
                </span>
            </div>

            <OptionalSignButton
                isSignedByMe={isSignedByMe}
                infoHash={torrent.info_hash}
                onSignSuccess={handleNewSignature}
            />

            <TorrentSignaturesList signatures={allSignatures} />
        </div>
    );
};