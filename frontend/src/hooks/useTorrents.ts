import { useQuery, useMutation, useQueryClient, type UseQueryResult, type UseMutationResult } from '@tanstack/react-query';
import type { TorrentDTO, TorrentIdentity } from '../types/model/models';
import { torrentService } from '../api/torrentService';

export const useTorrents = (): UseQueryResult<TorrentDTO[], Error> => {
    return useQuery<TorrentDTO[], Error>({
        queryKey: ['torrents'],
        queryFn: (): Promise<TorrentDTO[]> => torrentService.getTorrents(),
        staleTime: 5000,
    });
};

export const useCreateTorrent = (): UseMutationResult<TorrentDTO, Error, File> => {
    const queryClient = useQueryClient();
    return useMutation<TorrentDTO, Error, File>({
        mutationFn: (file: File): Promise<TorrentDTO> => torrentService.createTorrent(file),
        onSuccess: (): void => {
            void queryClient.invalidateQueries({ queryKey: ['torrents'] });
        }
    });
};

export const useAddTorrent = (): UseMutationResult<void, Error, TorrentIdentity> => {
    const queryClient = useQueryClient();
    return useMutation<void, Error, TorrentIdentity>({
        mutationFn: (identity: TorrentIdentity): Promise<void> => {
            return torrentService.addTorrent(identity);
        },
        onSuccess: (): void => {
            void queryClient.invalidateQueries({ queryKey: ['torrents'] });
        }
    });
};

export const useSignTorrent = (): UseMutationResult<TorrentDTO, Error, TorrentIdentity> => {
    const queryClient = useQueryClient();
    return useMutation<TorrentDTO, Error, TorrentIdentity>({
        mutationFn: (identity: TorrentIdentity): Promise<TorrentDTO> => {
            return torrentService.signTorrent(identity);
        },
        onSuccess: (): void => {
            void queryClient.invalidateQueries({ queryKey: ['torrents'] });
        }
    });
};

export const usePauseTorrent = (): UseMutationResult<void, Error, TorrentIdentity> => {
    const queryClient = useQueryClient();
    return useMutation<void, Error, TorrentIdentity>({
        mutationFn: (identity: TorrentIdentity): Promise<void> => {
            return torrentService.pauseTorrent(identity);
        },
        onSuccess: (): void => {
            void queryClient.invalidateQueries({ queryKey: ['torrents'] });
        }
    });
};

export const useResumeTorrent = (): UseMutationResult<void, Error, TorrentIdentity> => {
    const queryClient = useQueryClient();
    return useMutation<void, Error, TorrentIdentity>({
        mutationFn: (identity: TorrentIdentity): Promise<void> => {
            return torrentService.resumeTorrent(identity);
        },
        onSuccess: (): void => {
            void queryClient.invalidateQueries({ queryKey: ['torrents'] });
        }
    });
};

export const useDeleteTorrent = (): UseMutationResult<void, Error, TorrentIdentity> => {
    const queryClient = useQueryClient();
    return useMutation<void, Error, TorrentIdentity>({
        mutationFn: (identity: TorrentIdentity): Promise<void> => {
            return torrentService.deleteTorrent(identity);
        },
        onSuccess: (): void => {
            void queryClient.invalidateQueries({ queryKey: ['torrents'] });
        }
    });
};

export const useDeletePublishedTorrent = (): UseMutationResult<void, Error, TorrentIdentity> => {
    const queryClient = useQueryClient();
    return useMutation<void, Error, TorrentIdentity>({
        mutationFn: (identity: TorrentIdentity): Promise<void> => {
            return torrentService.deletePublishedTorrent(identity);
        },
        onSuccess: (): void => {
            void queryClient.invalidateQueries({ queryKey: ['store-torrents'] });
            void queryClient.invalidateQueries({ queryKey: ['torrents'] });
        }
    });
};
