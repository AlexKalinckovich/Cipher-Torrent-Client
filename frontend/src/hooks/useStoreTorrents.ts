import { useQuery, type UseQueryResult } from '@tanstack/react-query';
import type { StoreTorrent } from '../types/model/models';
import { torrentService } from '../api/torrentService';

export const useStoreTorrents = (): UseQueryResult<StoreTorrent[], Error> => {
    return useQuery<StoreTorrent[], Error>({
        queryKey: ['store-torrents'],
        queryFn: (): Promise<StoreTorrent[]> => torrentService.getStoreTorrents(),
        staleTime: 5000,
    });
};