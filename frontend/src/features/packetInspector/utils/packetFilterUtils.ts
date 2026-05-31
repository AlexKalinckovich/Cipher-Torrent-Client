import type { PacketLog } from '@/types/model/models.ts';

const safeSearchString = (value: string | undefined): string => {
    return value ?? '';
};

const checkMatch = (target: string, query: string): boolean => {
    return target.toLowerCase().includes(query);
};

const isPacketMatch = (packet: PacketLog, search: string): boolean => {
    const query: string = search.toLowerCase();
    const safeIp: string = safeSearchString(packet.peer_ip);
    const isIpMatch: boolean = checkMatch(safeIp, query);
    const isTypeMatch: boolean = checkMatch(packet.message_type, query);
    return isIpMatch || isTypeMatch;
};

export const filterPackets = (packets: PacketLog[], search: string): PacketLog[] => {
    if (!search) {
        return packets;
    }
    return packets.filter((p: PacketLog): boolean => isPacketMatch(p, search));
};