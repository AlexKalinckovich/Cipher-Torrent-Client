import React, { useState, useCallback, memo } from 'react';
import { CloseOutlined } from '@ant-design/icons';
import { message } from 'antd';
import type { TorrentAddRequest } from '@/types/model/models.ts';
import { MagnetTab } from './Tabs/MagnetTab';
import { FileTab } from './Tabs/FileTab';
import styles from './AddTorrentModal.module.css';

export interface AddTorrentModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSubmit: (req: TorrentAddRequest) => void;
}

type TabKey = 'magnet' | 'file';

interface FormProps {
    magnetUri: string;
    setMagnetUri: (val: string) => void;
    selectedFile: File | null;
    setSelectedFile: (file: File) => void;
}

const MagnetWrapper: React.FC<FormProps> = ({ magnetUri, setMagnetUri }) => (
    <MagnetTab value={magnetUri} onChange={setMagnetUri} />
);

const FileWrapper: React.FC<FormProps> = ({ selectedFile, setSelectedFile }) => (
    <FileTab fileName={selectedFile?.name ?? null} onFileSelect={setSelectedFile} />
);

const TAB_COMPONENTS: Record<TabKey, React.FC<FormProps>> = {
    magnet: MagnetWrapper,
    file: FileWrapper
};

const getTabClass = (current: TabKey, active: TabKey): string => {
    if (current === active) {
        return `${styles.tabBtn} ${styles.tabBtnActive}`;
    }
    return styles.tabBtn;
};

const stopPropagation = (e: React.MouseEvent): void => {
    e.stopPropagation();
};

const convertFileToBase64 = (file: File, callback: (base64: string) => void): void => {
    const reader = new FileReader();
    reader.onloadend = (): void => {
        const result: string = reader.result as string;
        const base64Data: string = result.split(',')[1];
        callback(base64Data);
    };
    reader.readAsDataURL(file);
};

const AddTorrentModalComponent: React.FC<AddTorrentModalProps> = ({ isOpen, onClose, onSubmit }) => {
    const [activeTab, setActiveTab] = useState<TabKey>('magnet');
    const [magnetUri, setMagnetUri] = useState<string>('');
    const [savePath, setSavePath] = useState<string>('/downloads/cipher');
    const [selectedFile, setSelectedFile] = useState<File | null>(null);

    const handleSetMagnet = useCallback((): void => setActiveTab('magnet'), []);
    const handleSetFile = useCallback((): void => setActiveTab('file'), []);

    const handlePathChange = useCallback((e: React.ChangeEvent<HTMLInputElement>): void => {
        setSavePath(e.target.value);
    }, []);

    const handleSubmitMagnet = useCallback((): void => {
        onSubmit({ magnet_uri: magnetUri, save_path: savePath });
        onClose();
    }, [magnetUri, savePath, onSubmit, onClose]);

    const handleSubmitFile = useCallback((): void => {
        if (!selectedFile) {
            void message.error('PLEASE SELECT A .TORRENT FILE');
            return;
        }
        convertFileToBase64(selectedFile, (base64: string): void => {
            onSubmit({ file_base64: base64, save_path: savePath });
            onClose();
        });
    }, [selectedFile, savePath, onSubmit, onClose]);

    const executeSubmit = useCallback((): void => {
        if (activeTab === 'magnet') {
            handleSubmitMagnet();
            return;
        }
        handleSubmitFile();
    }, [activeTab, handleSubmitMagnet, handleSubmitFile]);

    if (!isOpen) {
        return null;
    }

    const ActiveForm = TAB_COMPONENTS[activeTab];

    return (
        <div className={styles.overlay} onClick={onClose}>
            <div className={styles.modal} onClick={stopPropagation}>
                <div className={styles.header}>
                    <h2 className={styles.title}>ADD NEW TORRENT</h2>
                    <button type="button" className={styles.closeBtn} onClick={onClose}>
                        <CloseOutlined />
                    </button>
                </div>
                <div className={styles.tabHeader}>
                    <button type="button" className={getTabClass('magnet', activeTab)} onClick={handleSetMagnet}>MAGNET LINK</button>
                    <button type="button" className={getTabClass('file', activeTab)} onClick={handleSetFile}>TORRENT FILE</button>
                </div>
                <div className={styles.content}>
                    <ActiveForm
                        magnetUri={magnetUri}
                        setMagnetUri={setMagnetUri}
                        selectedFile={selectedFile}
                        setSelectedFile={setSelectedFile}
                    />

                    <div className={styles.inputGroup}>
                        <span className={styles.label}>SAVE PATH</span>
                        <input
                            type="text"
                            className={styles.input}
                            value={savePath}
                            onChange={handlePathChange}
                        />
                    </div>

                    <button type="button" className={styles.submitBtn} onClick={executeSubmit}>
                        INITIALIZE DOWNLOAD
                    </button>
                </div>
            </div>
        </div>
    );
};

export const AddTorrentModal = memo(AddTorrentModalComponent);