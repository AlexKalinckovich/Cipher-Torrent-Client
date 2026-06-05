import React, { useState, useCallback } from 'react';
import { CloseOutlined } from '@ant-design/icons';
import type { Torrent } from '@/types/model/models.ts';
import { formatBytes } from '@/utils/DashboardScreenUtils/bytesFormatter.ts';
import { FilesTab } from './Tabs/FilesTab';
import { PeersTab } from './Tabs/PeersTab';
import { SignaturesTab } from './Tabs/SignaturesTab';
import styles from './TorrentDetailsModal.module.css';

export interface TorrentDetailsModalProps {
    torrent: Torrent | null;
    onClose: () => void;
}

type TabKey = 'files' | 'peers' | 'signatures';

interface TabProps {
    torrent: Torrent;
}

const TAB_COMPONENTS: Record<TabKey, React.FC<TabProps>> = {
    files: FilesTab,
    peers: PeersTab,
    signatures: SignaturesTab
};

const stopPropagation = (e: React.MouseEvent): void => {
    e.stopPropagation();
};

const getTabClass = (current: TabKey, active: TabKey): string => {
    if (current === active) {
        return `${styles.tabBtn} ${styles.tabBtnActive}`;
    }
    return styles.tabBtn;
};

const getPieceSizeText = (bytes: number | undefined): string => {
    if (!bytes) {
        return 'PIECE SIZE: UNKNOWN';
    }
    return `PIECE SIZE: ${formatBytes(bytes)}`;
};

export const TorrentDetailsModal: React.FC<TorrentDetailsModalProps> = ({ torrent, onClose }) => {
    const [activeTab, setActiveTab] = useState<TabKey>('files');

    const handleSetFiles = useCallback((): void => setActiveTab('files'), []);
    const handleSetPeers = useCallback((): void => setActiveTab('peers'), []);
    const handleSetSignatures = useCallback((): void => setActiveTab('signatures'), []);

    if (!torrent) {
        return null;
    }

    const ActiveComponent = TAB_COMPONENTS[activeTab];

    return (
        <div className={styles.overlay} onClick={onClose}>
            <div className={styles.modal} onClick={stopPropagation}>
                <div className={styles.header}>
                    <div className={styles.titleBlock}>
                        <h2 className={styles.title}>{torrent.name}</h2>
                        <span className={styles.pieceInfo}>{getPieceSizeText(torrent.piece_length)}</span>
                    </div>
                    <button type="button" className={styles.closeBtn} onClick={onClose}>
                        <CloseOutlined />
                    </button>
                </div>
                <div className={styles.tabHeader}>
                    <button type="button" className={getTabClass('files', activeTab)} onClick={handleSetFiles}>FILES</button>
                    <button type="button" className={getTabClass('peers', activeTab)} onClick={handleSetPeers}>PEERS</button>
                    <button type="button" className={getTabClass('signatures', activeTab)} onClick={handleSetSignatures}>SIGNATURES</button>
                </div>
                <div className={styles.content}>
                    <ActiveComponent torrent={torrent} />
                </div>
            </div>
        </div>
    );
};