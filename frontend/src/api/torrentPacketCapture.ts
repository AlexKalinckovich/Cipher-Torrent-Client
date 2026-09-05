import type { PacketLog, PacketDirection, ReputationReceipt } from '@/types/model/models.ts';
import type { Torrent } from './webtorrent';
import type { Wire } from 'bittorrent-protocol';

/**
 * Maps bittorrent-protocol wire events into the app's PacketLog shape.
 *
 * WebTorrent's Torrent emits a "wire" event for every connected peer. The wire
 * is a bittorrent-protocol duplex stream exposing per-message events:
 *   request, piece, cancel, have, bitfield, choke, unchoke, interested,
 *   uninterested, extended, keep-alive, handshake, unknownmessage.
 *
 * We subscribe to these and produce the same PacketLog the backend WebSocket
 * used to emit, so the PacketInspector UI keeps working unchanged.
 */

export interface WirePacketListener {
    (packet: PacketLog): void;
}

let packetSequence = 0;

const nextPacketId = (): number => {
    packetSequence += 1;
    return packetSequence;
};

const nowIso = (): string => new Date().toISOString();

const baseLog = (
    direction: PacketDirection,
    messageType: string,
    wire: Wire,
    size: number,
): PacketLog => ({
    id: nextPacketId(),
    timestamp: nowIso(),
    direction,
    message_type: messageType,
    peer_ip: wire.peerId || 'unknown-peer',
    size_bytes: size,
});

type PacketBuilder = (...args: unknown[]) => PacketLog;

const sizeOf = (value: unknown): number => {
    if (value && typeof value === 'object' && 'length' in value) {
        return Number((value as { length: unknown }).length) || 0;
    }
    return 0;
};

const asNumber = (value: unknown): number => Number(value) || 0;

/**
 * Attach a packet listener to a single peer wire.
 * Returns a cleanup function.
 */
export const attachWirePacketListener = (wire: Wire, onPacket: WirePacketListener): (() => void) => {
    const cleanupFns: Array<() => void> = [];

    const listen = (event: string, build: PacketBuilder): void => {
        const handler = (...args: unknown[]): void => {
            onPacket(build(...args));
        };
        wire.on(event, handler);
        cleanupFns.push(() => {
            wire.off(event, handler);
        });
    };

    // --- Inbound message events ---
    listen('request', (index, offset, length) => {
        const log = baseLog('inbound', 'request', wire, asNumber(length));
        log.piece_index = asNumber(index);
        log.chunk_offset = asNumber(offset);
        log.chunk_length = asNumber(length);
        log.parsed_info = `request piece #${asNumber(index)} @${asNumber(offset)} (${asNumber(length)} B)`;
        return log;
    });

    listen('piece', (index, offset, buffer) => {
        const size = sizeOf(buffer);
        const log = baseLog('inbound', 'piece', wire, size);
        log.piece_index = asNumber(index);
        log.chunk_offset = asNumber(offset);
        log.chunk_length = size;
        log.parsed_info = `piece #${asNumber(index)} @${asNumber(offset)} (${size} B)`;
        return log;
    });

    listen('cancel', (index, offset, length) => {
        const log = baseLog('inbound', 'cancel', wire, asNumber(length));
        log.piece_index = asNumber(index);
        log.chunk_offset = asNumber(offset);
        log.chunk_length = asNumber(length);
        return log;
    });

    listen('have', (index) => {
        const log = baseLog('inbound', 'have', wire, 0);
        log.piece_index = asNumber(index);
        log.parsed_info = `have piece #${asNumber(index)}`;
        return log;
    });

    listen('bitfield', (bitfield) => baseLog('inbound', 'bitfield', wire, sizeOf(bitfield)));

    listen('choke', () => baseLog('inbound', 'choke', wire, 0));
    listen('unchoke', () => baseLog('inbound', 'unchoke', wire, 0));
    listen('interested', () => baseLog('inbound', 'interested', wire, 0));
    listen('uninterested', () => baseLog('inbound', 'uninterested', wire, 0));

    listen('extended', (ext, buf) => {
        const log = baseLog('inbound', 'extended', wire, sizeOf(buf));
        log.parsed_info = `extended message: ${String(ext)}`;
        return log;
    });

    listen('handshake', (infoHash) => {
        const log = baseLog('inbound', 'handshake', wire, 68);
        log.parsed_info = `handshake infoHash=${String(infoHash)}`;
        return log;
    });

    listen('unknownmessage', (buffer) => baseLog('inbound', 'unknownmessage', wire, sizeOf(buffer)));

    // --- Outbound message events ---
    listen('request', (index, offset, length) => {
        const log = baseLog('outbound', 'request', wire, asNumber(length));
        log.piece_index = asNumber(index);
        log.chunk_offset = asNumber(offset);
        log.chunk_length = asNumber(length);
        return log;
    });

    listen('have', (index) => {
        const log = baseLog('outbound', 'have', wire, 0);
        log.piece_index = asNumber(index);
        return log;
    });

    // Outbound piece uploads are observed via the 'upload' counter.
    const onUpload = (bytes: number): void => {
        const log = baseLog('outbound', 'piece', wire, bytes);
        log.parsed_info = `uploaded ${bytes} B`;
        onPacket(log);
    };
    wire.on('upload', onUpload);
    cleanupFns.push(() => {
        wire.off('upload', onUpload);
    });

    return () => {
        cleanupFns.forEach((cleanup) => cleanup());
    };
};

/**
 * Attach a packet listener to a WebTorrent torrent (all current + future wires).
 * Returns a cleanup function.
 */
export const attachTorrentPacketListener = (
    torrent: Torrent,
    onPacket: WirePacketListener,
): (() => void) => {
    const cleanupFns = new Set<() => void>();

    const attachToWire = (wire: Wire): void => {
        const cleanup = attachWirePacketListener(wire, onPacket);
        cleanupFns.add(cleanup);
    };

    torrent.on('wire', attachToWire);

    // Attach to wires already present.
    const existingWires = (torrent as unknown as { wires?: Wire[] }).wires ?? [];
    existingWires.forEach(attachToWire);

    return () => {
        torrent.off('wire', attachToWire);
        cleanupFns.forEach((cleanup) => cleanup());
        cleanupFns.clear();
    };
};

export const hasReputationExtension = (packet: PacketLog): boolean =>
    packet.message_type === 'extended' && (packet.parsed_info?.includes('reputation') ?? false);

export const extractReputationPayload = (packet: PacketLog): ReputationReceipt | undefined => {
    if (!hasReputationExtension(packet)) {
        return undefined;
    }
    try {
        const payload = JSON.parse(packet.parsed_info?.replace(/^extended message: /, '') ?? '{}');
        return payload as ReputationReceipt;
    } catch {
        return undefined;
    }
};