export interface TorrentSignature {
    signer_user_id: number;
    signer_public_key: string;
    signature_bytes: string;
    is_valid: boolean;
    signed_at: string;
}

export interface SignatureDTO {
    signer_public_key: string; // Base64 or Hex string depending on your mapper
    signature_bytes: string;
    timestamp: number;
}
