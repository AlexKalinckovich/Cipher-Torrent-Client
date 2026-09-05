// WebTorrent browser entry.
//
// `webtorrent`'s raw ESM source (index.js -> lib/*.js) is Node-oriented and
// references globals like `process`, `events`, `os`, `stream` that don't exist
// in the browser. WebTorrent ships an officially browser-targeted bundle at
// `dist/webtorrent.min.js` that inlines all of these shims (process, events,
// os, stream, ...) and exposes a single ESM default export.
//
// We import that prebuilt browser bundle here so the app does not need any
// Node polyfill aliasing. TypeScript types come from `@types/webtorrent`.
import WebTorrent from 'webtorrent/dist/webtorrent.min.js';
import type { Torrent, TorrentFile, TorrentOptions } from 'webtorrent';

export type { Torrent, TorrentFile, TorrentOptions };
export default WebTorrent;