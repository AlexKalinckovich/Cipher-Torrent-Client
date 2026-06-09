import React, { useCallback, memo } from 'react';
import { EditOutlined, LoadingOutlined } from '@ant-design/icons';
import { message } from 'antd';
import { useSignTorrent } from '@/hooks/useTorrents.ts';
import type { SignatureCreateResponse } from '@/types/model/models.ts';
import styles from './SignTorrentButton.module.css';

interface SignTorrentButtonProps {
    infoHash: string;
    onSignSuccess: (res: SignatureCreateResponse) => void;
}

interface OptionalSignButtonProps extends SignTorrentButtonProps {
    isSignedByMe: boolean;
}

const SignTorrentButtonComponent: React.FC<SignTorrentButtonProps> = ({ infoHash, onSignSuccess }) => {
    const { mutate: signTorrent, isPending } = useSignTorrent();

    const handleSign = useCallback((): void => {
        signTorrent(infoHash, {
            onSuccess: (response: SignatureCreateResponse): void => {
                void message.success('TORRENT CRYPTOGRAPHICALLY SIGNED');
                onSignSuccess(response);
            },
            onError: (error: Error): void => {
                void message.error(`SIGNATURE FAILED: ${error.message}`);
            }
        });
    }, [infoHash, signTorrent, onSignSuccess]);

    return (
        <div className={styles.actionContainer}>
            <button type="button" className={styles.signBtn} onClick={handleSign} disabled={isPending}>
                {isPending ? <LoadingOutlined /> : <EditOutlined />}
                {isPending ? 'GENERATING SIGNATURE...' : 'SIGN THIS TORRENT'}
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