import type { ReputationReceipt } from '../Reputation/reputationReceipt';

export type PacketDirection = 'inbound' | 'outbound';

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