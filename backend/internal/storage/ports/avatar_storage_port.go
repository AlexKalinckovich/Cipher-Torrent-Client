package ports

import "context"

type AvatarUploadRequest struct {
	PublicKey []byte
	FileBytes []byte
	Ext       string
}

type AvatarStoragePort interface {
	UploadAvatar(ctx context.Context, req AvatarUploadRequest) (string, error)
	DeleteAvatar(ctx context.Context, publicKey []byte) error
}