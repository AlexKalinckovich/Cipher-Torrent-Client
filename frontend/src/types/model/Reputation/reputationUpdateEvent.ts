export type ReputationUpdateEventType = 'reputation:update';

export interface ReputationUpdateEvent {
    event_type: ReputationUpdateEventType;
    from_peer_id: string;
    to_peer_id: string;
    new_reputation_score: number;
}