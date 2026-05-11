import { useQuery, type UseQueryResult } from '@tanstack/react-query';
import type {Torrent} from '../types/model/models.ts';
import { torrentService } from '../api/torrentService';

export const useTorrents = (): UseQueryResult<Torrent[], Error> => {
    return useQuery<Torrent[], Error>({
        queryKey: ['torrents'],
        queryFn: () => torrentService.getTorrents(),
        staleTime: 5000,
    });
};