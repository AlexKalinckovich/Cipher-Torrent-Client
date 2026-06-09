import React, { useCallback, useState } from 'react';
import type { Torrent, TorrentSignature, SignatureCreateResponse } from '@/types/model/models.ts';
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

export const SignaturesTab: React.FC<TabProps> = ({ torrent }) => {
    const initialSignatures = useTorrentSignaturesMock(torrent.info_hash);
    const [addedSignatures, setAddedSignatures] = useState<TorrentSignature[]>([]);
    const [isSignedByMe, setIsSignedByMe] = useState<boolean>(torrent.is_signed_by_me ?? false);

    const allSignatures = [...addedSignatures, ...initialSignatures];

    const handleNewSignature = useCallback((response: SignatureCreateResponse): void => {
        const newSig: TorrentSignature = {
            signer_user_id: 0,
            signer_public_key: 'Ed25519_PubKey_Current_User',
            signature_bytes: `SIG_ID_${response.signature_id}`,
            is_valid: true,
            signed_at: response.signed_at
        };
        setAddedSignatures((prev: TorrentSignature[]): TorrentSignature[] => [newSig, ...prev]);
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