import React, { memo } from 'react';
import { EditOutlined, LoadingOutlined } from '@ant-design/icons';
import styles from './SignTorrentButton.module.css';

interface SignTorrentButtonProps {
    onSignClick: () => void;
    isPending?: boolean;
}

interface OptionalSignButtonProps extends SignTorrentButtonProps {
    isSignedByMe: boolean;
}

const SignTorrentButtonComponent: React.FC<SignTorrentButtonProps> = ({ onSignClick, isPending }) => {
    return (
        <div className={styles.actionContainer}>
            <button type="button" className={styles.signBtn} onClick={onSignClick} disabled={isPending}>
                {isPending ? <LoadingOutlined /> : <EditOutlined />}
                {isPending ? 'GENERATING SIGNATURE...' : 'SIGN THIS TORRENT'}
            </button>
        </div>
    );
};

export const SignTorrentButton = memo(SignTorrentButtonComponent);

export const OptionalSignButton: React.FC<OptionalSignButtonProps> = ({ isSignedByMe, onSignClick, isPending }) => {
    if (isSignedByMe) {
        return null;
    }
    return <SignTorrentButton onSignClick={onSignClick} isPending={isPending} />;
};