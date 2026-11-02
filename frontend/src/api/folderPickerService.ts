import { directoryOpen, supported } from 'browser-fs-access';

const DB_NAME = 'cipher-torrent';
const DB_VERSION = 1;
const STORE_NAME = 'download-folders';

export interface PickedFolder {
    /** Stable id for this choice (used as the IndexedDB key). */
    id: string;
    /** Display name of the chosen directory (e.g. "Downloads"). */
    name: string;
    /** The directory handle; usable later to write files into the folder. */
    handle: FileSystemDirectoryHandle | null;
}

/** True when the browser can open a directory picker at all (native or legacy). */
export const supportsFolderPicker = (): boolean => supported;

function openDB(): Promise<IDBDatabase> {
    return new Promise((resolve, reject) => {
        const request = indexedDB.open(DB_NAME, DB_VERSION);
        request.onupgradeneeded = () => {
            const db = request.result;
            if (!db.objectStoreNames.contains(STORE_NAME)) {
                db.createObjectStore(STORE_NAME, { keyPath: 'id' });
            }
        };
        request.onsuccess = () => resolve(request.result);
        request.onerror = () => reject(request.error);
    });
}

/**
 * Shows the native folder picker, persists the chosen directory handle to
 * IndexedDB and returns it. Throws if the picker is unavailable or cancelled.
 */
export async function pickDownloadFolder(): Promise<PickedFolder> {
    if (!supported) {
        throw new Error(
            'This browser does not support the folder picker. Please use a recent version of Chrome, Edge or Firefox.',
        );
    }

    const files = await directoryOpen({
        recursive: true,
        mode: 'readwrite',
    });

        
    const first = files[0] as
        | { directoryHandle?: FileSystemDirectoryHandle | null }
        | undefined;
    const handle = first?.directoryHandle ?? null;
    const name = handle?.name ?? 'Downloads';

    const folder: PickedFolder = {
        id: `folder_${Date.now()}`,
        name,
        handle,
    };

    const db = await openDB();
    await new Promise<void>((resolve, reject) => {
        const tx = db.transaction(STORE_NAME, 'readwrite');
        tx.objectStore(STORE_NAME).put(folder);
        tx.oncomplete = () => resolve();
        tx.onerror = () => reject(tx.error);
    });

    return folder;
}

export async function getLastPickedFolder(): Promise<PickedFolder | null> {
    try {
        const db = await openDB();
        return await new Promise<PickedFolder | null>((resolve, reject) => {
            const tx = db.transaction(STORE_NAME, 'readonly');
            const all = tx.objectStore(STORE_NAME).getAll();
            all.onsuccess = () => {
                const folders = all.result as PickedFolder[];
                resolve(folders.length > 0 ? folders[folders.length - 1] : null);
            };
            all.onerror = () => reject(all.error);
        });
    } catch (err) {
        console.warn('[FolderPicker] Failed to read last folder:', err);
        return null;
    }
}