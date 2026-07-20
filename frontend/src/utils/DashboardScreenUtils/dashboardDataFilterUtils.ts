import type {TorrentDTO} from "../../types/model/Torrent/torrentEntity.ts";

const filterByStatus = (torrents: TorrentDTO[], status: string | null): TorrentDTO[] => {
    if (!status) return torrents;
    return torrents.filter((t: TorrentDTO) => t.status === status);
};

const filterBySearch = (torrents: TorrentDTO[], search: string): TorrentDTO[] => {
    if (!search) return torrents;
    return torrents.filter((t: TorrentDTO) => t.name.toLowerCase().includes(search.toLowerCase()));
};

export const getFilteredTorrents = (torrents: TorrentDTO[], status: string | null, search: string): TorrentDTO[] => {
    const statusFiltered = filterByStatus(torrents, status);
    return filterBySearch(statusFiltered, search);
};