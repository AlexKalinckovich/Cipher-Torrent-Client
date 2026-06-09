import { useQuery, useMutation, useQueryClient, type UseQueryResult, type UseMutationResult } from '@tanstack/react-query';
import type { Torrent, TorrentAddRequest, TorrentAddResponse, SignatureCreateResponse } from '../types/model/models';
import { torrentService } from '../api/torrentService';

export const useTorrents = (): UseQueryResult<Torrent[], Error> => {
    return useQuery<Torrent[], Error>({
        queryKey: ['torrents'],
        queryFn: (): Promise<Torrent[]> => torrentService.getTorrents(),
        staleTime: 5000,
    });
};

export const useAddTorrent = (): UseMutationResult<TorrentAddResponse, Error, TorrentAddRequest> => {
    const queryClient = useQueryClient();

    return useMutation<TorrentAddResponse, Error, TorrentAddRequest>({
        mutationFn: (request: TorrentAddRequest): Promise<TorrentAddResponse> => {
            return torrentService.addTorrent(request);
        },
        onSuccess: (): void => {
            void queryClient.invalidateQueries({ queryKey: ['torrents'] });
        }
    });
};

export const useSignTorrent = (): UseMutationResult<SignatureCreateResponse, Error, string> => {
    const queryClient = useQueryClient();

    return useMutation<SignatureCreateResponse, Error, string>({
        mutationFn: (infoHash: string): Promise<SignatureCreateResponse> => {
            return torrentService.signTorrent(infoHash);
        },
        onSuccess: (): void => {
            void queryClient.invalidateQueries({ queryKey: ['torrents'] });
        }
    });
};