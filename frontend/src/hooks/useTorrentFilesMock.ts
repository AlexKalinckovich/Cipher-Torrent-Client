import { useMemo } from 'react';
import type { Torrent, TorrentFile } from '@/types/model/models.ts';

const generateMockFiles = (): TorrentFile[] => {
    return [
        { path: 'Cyberpunk_Data_Core/kernel_exploit.bin', size_bytes: 45000000 },
        { path: 'Cyberpunk_Data_Core/payload_stage_1.dat', size_bytes: 2147483648 },
        { path: 'Cyberpunk_Data_Core/payload_stage_2.dat', size_bytes: 2147483648 },
        { path: 'Cyberpunk_Data_Core/readme_nfo.txt', size_bytes: 1024 }
    ];
};

const getFiles = (torrent: Torrent): TorrentFile[] => {
    const hasRealFiles: boolean = (torrent.files?.length ?? 0) > 0;
    if (hasRealFiles) {
        return torrent.files as TorrentFile[];
    }
    return generateMockFiles();
};

export const useTorrentFilesMock = (torrent: Torrent): TorrentFile[] => {
    return useMemo((): TorrentFile[] => {
        return getFiles(torrent);
    }, [torrent]);
};