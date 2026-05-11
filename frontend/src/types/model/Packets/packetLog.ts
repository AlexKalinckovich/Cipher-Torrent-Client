import type { ReputationReceipt } from '../Reputation/reputationReceipt';

export type PacketDirection = 'incoming' | 'outgoing';

export interface PacketLog {
    id: number;
    timestamp: string;
    direction: PacketDirection;
    message_type: string;
    peer_ip?: string;
    raw_payload_base64: string;
    parsed_info?: string;
    reputation_payload?: ReputationReceipt;
    size_bytes?: number;
}