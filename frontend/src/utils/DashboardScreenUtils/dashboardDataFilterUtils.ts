import type {Torrent} from "../../types/model/Torrent/torrent.ts";

const filterByStatus = (torrents: Torrent[], status: string | null): Torrent[] => {
    if (!status) return torrents;
    return torrents.filter((t: Torrent) => t.status === status);
};

const filterBySearch = (torrents: Torrent[], search: string): Torrent[] => {
    if (!search) return torrents;
    return torrents.filter((t: Torrent) => t.name.toLowerCase().includes(search.toLowerCase()));
};

export const getFilteredTorrents = (torrents: Torrent[], status: string | null, search: string): Torrent[] => {
    const statusFiltered = filterByStatus(torrents, status);
    return filterBySearch(statusFiltered, search);
};