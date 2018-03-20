package downloadService

import (
	"crypto/rand"
	"encoding/hex"
	"os"

	"github.com/ONSdigital/go-ns/log"
	"github.com/ONSdigital/go-ns/vault"
	"github.com/ONSdigital/s3crypto"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func sendV4FileToAWS(region, bucket, filename string) error {
	config := aws.NewConfig().WithRegion(region)

	sess, err := session.NewSession(config)
	if err != nil {
		log.ErrorC("failed to create session", err, nil)
		return err
	}

	v4File, err := os.Open(filename)
	if err != nil {
		log.ErrorC("failed to open file", err, nil)
		return err
	}
	log.Info("successfully retrieved file", nil)
	defer v4File.Close()

	putObject := &s3.PutObjectInput{
		Key:    &filename,
		Bucket: &bucket,
		Body:   v4File,
	}

	vaultAddress := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")
	client := s3crypto.New(sess, &s3crypto.Config{HasUserDefinedPSK: true})

	vc, err := vault.CreateVaultClient(token, vaultAddress, 3)
	if err != nil {
		return err
	}

	psk := createPSK()
	pskStr := hex.EncodeToString(psk)

	if err := vc.WriteKey("secret/shared/psk", filename, pskStr); err != nil {
		return err
	}

	_, err = client.PutObjectWithPSK(putObject, psk)
	if err != nil {
		return err
	}

	return nil
}

func deleteS3File(region, bucket, filename string) error {

	config := aws.NewConfig().WithRegion(region)

	session, err := session.NewSession(config)
	if err != nil {
		log.ErrorC("failed to create session", err, nil)
		return err
	}

	svc := s3.New(session)

	input := &s3.DeleteObjectInput{
		Key:    aws.String(filename),
		Bucket: aws.String(bucket),
	}

	if _, err = svc.DeleteObject(input); err != nil {
		log.ErrorC("failed to remove file", err, nil)
		return err
	}

	return nil
}

func createPSK() []byte {
	key := make([]byte, 16)
	rand.Read(key)

	return key
}
