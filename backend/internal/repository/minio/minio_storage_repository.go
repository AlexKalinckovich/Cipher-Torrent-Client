package minio

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"log"

	storage_ports "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/storage/ports"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorageRepository struct {
	client     *minio.Client
	bucketName string
}

func NewMinioStorageRepository(
	endpoint, accessKeyID, secretAccessKey, bucketName string,
	useSSL bool,
) (*MinioStorageRepository, error) {
	client, err := minio.New(
		endpoint,
		&minio.Options{
			Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
			Secure: useSSL,
		},
	)
	if err != nil {
		log.Printf("Failed to create a new Minio client: %v", err)
		return nil, err
	}

	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		log.Printf("Error checking if bucket exists: %v", err)
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("Error creating bucket %s: %v", bucketName, err)
			return nil, err
		}
		log.Printf("Successfully created bucket: %s", bucketName)
	} else {
		log.Printf("Bucket %s already exists", bucketName)
	}

	return &MinioStorageRepository{client: client, bucketName: bucketName}, nil
}

func (r *MinioStorageRepository) UploadBaseTorrent(ctx context.Context, req storage_ports.StorageUploadRequest) error {
	key := r.buildObjectKey(req.InfoHash, req.CreatorPubKey)
	return r.putObject(ctx, key, req.FileBytes)
}

func (r *MinioStorageRepository) putObject(ctx context.Context, key string, data []byte) error {
	reader := bytes.NewReader(data)
	size := int64(len(data))
	_, err := r.client.PutObject(ctx, r.bucketName, key, reader, size, minio.PutObjectOptions{ContentType: "application/x-bittorrent"})
	return err
}

func (r *MinioStorageRepository) DownloadBaseTorrent(ctx context.Context, req storage_ports.StorageIdentityRequest) ([]byte, error) {
	key := r.buildObjectKey(req.InfoHash, req.CreatorPubKey)
	return r.getObject(ctx, key)
}

func (r *MinioStorageRepository) getObject(ctx context.Context, key string) ([]byte, error) {
	obj, err := r.client.GetObject(ctx, r.bucketName, key, minio.GetObjectOptions{})
	if err != nil {
		log.Println("Error getting object", err.Error())
		return nil, err
	}
	return r.readAndCloseObject(obj)
}

func (r *MinioStorageRepository) readAndCloseObject(obj *minio.Object) ([]byte, error) {
	defer obj.Close()
	return io.ReadAll(obj)
}

func (r *MinioStorageRepository) DeleteTorrent(ctx context.Context, req storage_ports.StorageIdentityRequest) error {
	key := r.buildObjectKey(req.InfoHash, req.CreatorPubKey)
	return r.client.RemoveObject(ctx, r.bucketName, key, minio.RemoveObjectOptions{})
}

func (r *MinioStorageRepository) buildObjectKey(infoHash []byte, creatorPubKey []byte) string {
	return fmt.Sprintf("torrents/%s/%s.torrent", hex.EncodeToString(infoHash), hex.EncodeToString(creatorPubKey))
}
