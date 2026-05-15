export interface ReputationReceipt {
    from_peer_id: string;
    to_peer_id: string;
    piece_index: number;
    byte_count: number;
    signature: string;
    timestamp: string;
}

