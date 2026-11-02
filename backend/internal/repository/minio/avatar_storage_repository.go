package minio

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"

	storage_ports "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/storage/ports"
	"github.com/minio/minio-go/v7"
)

// AvatarStorageRepository stores user profile images in MinIO under a unique
// per-user path derived from their public key, e.g.
//
//	avatars/<base64url(public_key)>/avatar.<ext>
//
// The object key is persisted in users.avatar_path so it can be served/refreshed.
type AvatarStorageRepository struct {
	client     *minio.Client
	bucketName string
}

func NewAvatarStorageRepository(client *minio.Client, bucketName string) *AvatarStorageRepository {
	return &AvatarStorageRepository{
		client:     client,
		bucketName: bucketName,
	}
}

func (r *AvatarStorageRepository) UploadAvatar(ctx context.Context, req storage_ports.AvatarUploadRequest) (string, error) {
	key := r.buildObjectKey(req.PublicKey, req.Ext)
	reader := bytes.NewReader(req.FileBytes)
	size := int64(len(req.FileBytes))

	_, err := r.client.PutObject(ctx, r.bucketName, key, reader, size, minio.PutObjectOptions{
		ContentType: contentTypeForExt(req.Ext),
	})
	if err != nil {
		log.Printf("Failed to upload avatar: %v", err)
		return "", err
	}
	return key, nil
}

func (r *AvatarStorageRepository) DeleteAvatar(ctx context.Context, publicKey []byte) error {
	// Remove any existing avatar object regardless of extension by listing the
	// per-user directory prefix.
	prefix := fmt.Sprintf("avatars/%s/", encodePublicKey(publicKey))
	objectCh := r.client.ListObjects(ctx, r.bucketName, minio.ListObjectsOptions{Prefix: prefix, Recursive: true})
	for obj := range objectCh {
		if err := r.client.RemoveObject(ctx, r.bucketName, obj.Key, minio.RemoveObjectOptions{}); err != nil {
			log.Printf("Failed to remove avatar object %s: %v", obj.Key, err)
			return err
		}
	}
	return nil
}

func (r *AvatarStorageRepository) buildObjectKey(publicKey []byte, ext string) string {
	return fmt.Sprintf("avatars/%s/avatar.%s", encodePublicKey(publicKey), ext)
}

func encodePublicKey(publicKey []byte) string {
	return base64.RawURLEncoding.EncodeToString(publicKey)
}

func contentTypeForExt(ext string) string {
	switch ext {
	case "jpg", "jpeg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "gif":
		return "image/gif"
	case "webp":
		return "image/webp"
	default:
		return "application/octet-stream"
	}
}