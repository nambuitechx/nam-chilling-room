package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func GetS3Object(bucket *string, key *string) ([]byte, error) {
	// Using the SDK's default configuration, load additional config
    // and credentials values from the environment variables, shared
    // credentials, and shared configuration files
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))
    if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
    }

	s3Service := s3.NewFromConfig(cfg)

	resp, err := s3Service.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: bucket,
		Key: key,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get object on aws s3, %v", err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read object body, %v", err)
	}

	return body, nil
}

func DownloadS3Object(bucket *string, key *string, localPath string) (string, error) {
	// Using the SDK's default configuration, load additional config
    // and credentials values from the environment variables, shared
    // credentials, and shared configuration files
    cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))

    if err != nil {
		return "", fmt.Errorf("unable to load SDK config, %v", err)
    }

	s3Service := s3.NewFromConfig(cfg)

	resp, err := s3Service.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: bucket,
		Key:    key,
	})

	if err != nil {
		return "", fmt.Errorf("failed to download object: %w", err)
	}

	defer resp.Body.Close()

	file, err := os.Create(localPath)

	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}

	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return localPath, nil
}

func StreamS3Object(bucket string, key string, chunkSize int, sendChan chan<- []byte) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-southeast-1"))

	if err != nil {
		return fmt.Errorf("unable to load SDK config: %w", err)
	}

	s3Service := s3.NewFromConfig(cfg)

	resp, err := s3Service.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	if err != nil {
		return fmt.Errorf("failed to get object on aws s3: %w", err)
	}

	defer resp.Body.Close()

	buf := make([]byte, chunkSize)

	for {
		n, err := resp.Body.Read(buf)

		if err != nil && err != io.EOF {
			return fmt.Errorf("failed reading S3 object body: %w", err)
		}

		if n == 0 {
			break
		}

		// Copy into fresh slice to avoid data overwrite
		chunk := make([]byte, n)
		copy(chunk, buf[:n])
		sendChan <- chunk

		// Simulate playback speed
		time.Sleep(50 * time.Millisecond)
	}
	return nil
}
