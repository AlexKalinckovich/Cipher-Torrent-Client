import React from 'react';
import styles from './AuthHeader.module.css';

export const AuthHeader: React.FC = () => {
    return (
        <>
            <div className={`${styles.backgroundBlob} ${styles.blobLeft}`} />
            <div className={`${styles.backgroundBlob} ${styles.blobRight}`} />
            <h1 className={styles.cyberTitle}>
                Welcome to <br /> Cipher-Torrent-Client!
            </h1>
        </>
    );
};
