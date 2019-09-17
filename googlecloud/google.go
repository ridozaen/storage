package googlecloud

import (
	"context"
	"errors"
	"fmt"
	"strings"

	cfg "git.bluebird.id/bluebird/util/config"

	"google.golang.org/api/option"

	log "github.com/sirupsen/logrus"

	gstorage "cloud.google.com/go/storage"
)

type GoogleCloudStorage struct {
	Storage
	Options    []option.ClientOption
	client     *gstorage.Client
	bucket     *gstorage.BucketHandle
	bucketName string
}

func Setup() (Storage, error) {
	bktName := cfg.Get("GOOGLE_STORAGE_BUCKET", "")
	if bktName == "" {
		return nil, errors.New("GOOGLE_STORAGE_BUCKET env must be set")
	}
	// GoogleCloudStorage.bucketName = bktName

	projectID := cfg.Get("GOOGLE_STORAGE_PROJECT_ID", "")
	if projectID == "" {
		return nil, errors.New("GOOGLE_STORAGE_PROJECT_ID env must be set")
	}

	location := cfg.Get("GOOGLE_STORAGE_LOCATION", "")
	if location == "" {
		return nil, errors.New("GOOGLE_STORAGE_LOCATION env must be set")
	}

	keyfile := cfg.Get("keyfile", "keyfile.json")

	ctx := context.Background()
	client, err := gstorage.NewClient(ctx, option.WithCredentialsFile(keyfile))
	if err != nil {
		fmt.Println("error: ", err)
		return nil, err
	}

	// GoogleCloudStorage.client = client

	bkt := client.Bucket(bktName)

	attrs := &gstorage.BucketAttrs{Location: location}
	err = bkt.Create(ctx, projectID, attrs)
	fmt.Println("err: ", err)
	if err == nil {
		log.Printf("Created Google Cloud Storage bucket %s in %s",
			bktName, location)
	}

	if err != nil {
		if !strings.Contains(err.Error(), "You already own this bucket") {
			return nil, err
		}

		log.Printf("Using existing Google Cloud Storage bucket %v", bktName)
	}

	// GoogleCloudStorage.bucket = bkt

	return &GoogleCloudStorage{
		bucket:     bkt,
		client:     client,
		bucketName: bktName,
	}, nil
}

func (gcs *GoogleCloudStorage) PublicURL(filename string) string {
	return "https://storage.googleapis.com/" + gcs.bucketName + "/" + filename
}

func (gcs *GoogleCloudStorage) Store(ctx context.Context, filename string, data []byte, metadata map[string]string) error {
	fmt.Println("file name: ", data)
	o := gcs.bucket.Object(filename)
	w := o.NewWriter(ctx)

	w.ObjectAttrs = gstorage.ObjectAttrs{
		Name:     filename,
		Metadata: metadata,
	}

	_, err := w.Write(data)
	if err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

func (gcs *GoogleCloudStorage) Delete(ctx context.Context, filename string) error {
	o := gcs.bucket.Object(filename)
	return o.Delete(ctx)
}
