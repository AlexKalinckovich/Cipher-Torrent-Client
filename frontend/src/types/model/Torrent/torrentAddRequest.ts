export interface TorrentAddRequest {
    magnet_uri?: string;
    file_base64?: string;
    save_path?: string;
}
export interface TorrentIdentity {
    info_hash: string;       // Hex string
    creator_pub_key: string; // Base64 string
}
