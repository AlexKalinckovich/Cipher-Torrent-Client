export type PacketDirection = 'inbound' | 'outbound';

export interface ReputationReceipt {
    info_hash: string;
    from_pub_key: string;
    to_pub_key: string;
    piece_index: number;
    byte_count: number;
    timestamp: number;
}

export interface PacketLog {
    id: number;
    timestamp: string;
    direction: PacketDirection;
    message_type: string;
    peer_ip?: string;
    piece_index?: number;
    chunk_offset?: number;
    chunk_length?: number;
    parsed_info?: string;
    reputation_payload?: ReputationReceipt;
    size_bytes?: number;
}