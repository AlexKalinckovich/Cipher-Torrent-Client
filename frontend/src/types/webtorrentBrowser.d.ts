declare module 'webtorrent/dist/webtorrent.min.js' {
    import type WebTorrent from 'webtorrent';
    const WebTorrentBrowser: typeof WebTorrent;
    export default WebTorrentBrowser;
}