import { useState } from 'react';
import type { Torrent, TorrentFile as ClientTorrentFile } from '@/api/webtorrent';

// Minimal typing for the File System Access API (Chrome/Edge only).
interface FileSystemDirectoryHandle {
    getFileHandle(name: string, options?: { create?: boolean }): Promise<FileSystemFileHandle>;
    getDirectoryHandle(name: string, options?: { create?: boolean }): Promise<FileSystemDirectoryHandle>;
}

interface FileSystemFileHandle {
    createWritable(): Promise<{ write(data: Blob): Promise<void>; close(): Promise<void> }>;
}

interface WindowWithShowDirectoryPicker {
    showDirectoryPicker?: (options?: { id?: string; mode?: string }) => Promise<FileSystemDirectoryHandle>;
}

const supportsDirectoryPicker = (): boolean => {
    return typeof (window as WindowWithShowDirectoryPicker).showDirectoryPicker === 'function';
};

const splitPath = (path: string): string[] => {
    return path.split('/').filter((part: string): boolean => part.length > 0);
};

const writeFileToDirectory = async (dir: FileSystemDirectoryHandle, path: string, blob: Blob): Promise<void> => {
    const parts = splitPath(path);
    const fileName = parts.pop();
    if (!fileName) {
        return;
    }

    let current = dir;
    for (const part of parts) {
        current = await current.getDirectoryHandle(part, { create: true });
    }

    const fileHandle = await current.getFileHandle(fileName, { create: true });
    const writable = await fileHandle.createWritable();
    await writable.write(blob);
    await writable.close();
};

const readTorrentFile = async (file: ClientTorrentFile): Promise<Blob> => {
    const arrayBuffer = await file.arrayBuffer();
    return new Blob([arrayBuffer]);
};

const downloadFileFallback = (blob: Blob, path: string): void => {
    const parts = splitPath(path);
    const fileName = parts.pop() ?? 'download';
    const url = URL.createObjectURL(blob);
    const anchor = document.createElement('a');
    anchor.href = url;
    anchor.download = fileName;
    document.body.appendChild(anchor);
    anchor.click();
    anchor.remove();
    URL.revokeObjectURL(url);
};

/**
 * Saves already-downloaded (client-side WebTorrent) files to a user-chosen folder.
 * Does NOT fetch any bytes from the backend.
 */
export const useDownloadTorrentFile = () => {
    const [isSaving, setIsSaving] = useState<boolean>(false);
    const [error, setError] = useState<Error | null>(null);

    const saveToFolder = async (torrent: Torrent): Promise<void> => {
        setIsSaving(true);
        setError(null);
        try {
            const files: ClientTorrentFile[] = torrent.files;

            if (supportsDirectoryPicker()) {
                const rootDir = await (window as WindowWithShowDirectoryPicker).showDirectoryPicker!({
                    id: 'cipher-torrent-downloads',
                    mode: 'readwrite'
                });

                for (const file of files) {
                    const blob = await readTorrentFile(file);
                    await writeFileToDirectory(rootDir, file.path, blob);
                }
            } else {
                for (const file of files) {
                    const blob = await readTorrentFile(file);
                    downloadFileFallback(blob, file.path);
                }
            }
        } catch (err) {
            setError(err instanceof Error ? err : new Error('Failed to save torrent files'));
            throw err;
        } finally {
            setIsSaving(false);
        }
    };

    return { saveToFolder, isSaving, error };
};