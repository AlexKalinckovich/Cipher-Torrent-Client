export interface TorrentSignature {
    signer_user_id: number;
    signer_public_key: string;
    signature_bytes: string;
    is_valid: boolean;
    signed_at: string;
}

