import React, { memo, useCallback } from 'react';
import styles from '../AddTorrentModal.module.css';

interface MagnetTabProps {
    value: string;
    onChange: (val: string) => void;
}

const MagnetTabComponent: React.FC<MagnetTabProps> = ({ value, onChange }) => {
    const handleChange = useCallback((e: React.ChangeEvent<HTMLInputElement>): void => {
        onChange(e.target.value);
    }, [onChange]);

    return (
        <div className={styles.inputGroup}>
            <span className={styles.label}>MAGNET URI</span>
            <input
                type="text"
                className={styles.input}
                placeholder="magnet:?xt=urn:btih:..."
                value={value}
                onChange={handleChange}
            />
        </div>
    );
};

export const MagnetTab = memo(MagnetTabComponent);