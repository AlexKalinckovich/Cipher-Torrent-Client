import { PacketLog } from './packetLog';

export type PacketEventType = 'packet';

export interface PacketEvent {
    event_type: PacketEventType;
    packet: PacketLog;
}