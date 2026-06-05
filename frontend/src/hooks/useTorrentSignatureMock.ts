import { useMemo } from 'react';
import type { TorrentSignature } from '@/types/model/models.ts';

const generateMockSignatures = (): TorrentSignature[] => {
    return [
        {
            signer_user_id: 1,
            signer_public_key: 'Ed25519_PubKey_Admin_X9A4B2K8',
            signature_bytes: 'SIG_BYTES_8F7A6B5C4D3E2F1A0B9C8D7E6F5A4B3C2D1E0F',
            is_valid: true,
            signed_at: new Date(Date.now() - 86400000).toISOString()
        },
        {
            signer_user_id: 42,
            signer_public_key: 'Ed25519_PubKey_User_M3L1P0Q5',
            signature_bytes: 'SIG_BYTES_1A2B3C4D5E6F7A8B9C0D1E2F3A4B5C6D7E8F9A',
            is_valid: true,
            signed_at: new Date(Date.now() - 172800000).toISOString()
        },
        {
            signer_user_id: 88,
            signer_public_key: 'Ed25519_PubKey_Guest_Z1X2C3V4',
            signature_bytes: 'SIG_BYTES_INVALID_CORRUPTED_PAYLOAD_9999999999',
            is_valid: false,
            signed_at: new Date(Date.now() - 259200000).toISOString()
        }
    ];
};

export const useTorrentSignaturesMock = (infoHash: string): TorrentSignature[] => {
    return useMemo((): TorrentSignature[] => {
        if (!infoHash) {
            return [];
        }
        return generateMockSignatures();
    }, [infoHash]);
};