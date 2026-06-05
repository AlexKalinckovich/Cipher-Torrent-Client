import React, { useState, useCallback, memo } from 'react';
import { EditOutlined, LoadingOutlined } from '@ant-design/icons';
import { message } from 'antd';
import type { TorrentSignature } from '@/types/model/models.ts';
import styles from './SignTorrentButton.module.css';

interface SignTorrentButtonProps {
    infoHash: string;
    onSignSuccess: (sig: TorrentSignature) => void;
}

interface OptionalSignButtonProps extends SignTorrentButtonProps {
    isSignedByMe: boolean;
}

const generateNewSignature = (): TorrentSignature => {
    return {
        signer_user_id: 777,
        signer_public_key: 'Ed25519_PubKey_Phantom_X9A4B2K8M3L1P0Q5R',
        signature_bytes: `SIG_BYTES_NEW_${Date.now()}_A9F8C7`,
        is_valid: true,
        signed_at: new Date().toISOString()
    };
};

const SignTorrentButtonComponent: React.FC<SignTorrentButtonProps> = ({ onSignSuccess }) => {
    const [isSigning, setIsSigning] = useState<boolean>(false);

    const executeSign = useCallback((): void => {
        const newSig: TorrentSignature = generateNewSignature();
        onSignSuccess(newSig);
        setIsSigning(false);
        void message.success('TORRENT CRYPTOGRAPHICALLY SIGNED');
    }, [onSignSuccess]);

    const handleSign = useCallback((): void => {
        setIsSigning(true);
        window.setTimeout(executeSign, 1200);
    }, [executeSign]);

    return (
        <div className={styles.actionContainer}>
            <button type="button" className={styles.signBtn} onClick={handleSign} disabled={isSigning}>
                {isSigning ? <LoadingOutlined /> : <EditOutlined />}
                {isSigning ? 'GENERATING SIGNATURE...' : 'SIGN THIS TORRENT'}
            </button>
        </div>
    );
};

export const SignTorrentButton = memo(SignTorrentButtonComponent);

export const OptionalSignButton: React.FC<OptionalSignButtonProps> = ({ isSignedByMe, infoHash, onSignSuccess }) => {
    if (isSignedByMe) {
        return null;
    }
    return <SignTorrentButton infoHash={infoHash} onSignSuccess={onSignSuccess} />;
};