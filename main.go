package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
)

var s3Client *s3.Client

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load env file")
	}

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", os.Getenv("R2_ACCOUNT_ID")),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(os.Getenv("R2_ACCESS_KEY_ID"), os.Getenv("R2_ACCESS_KEY_SECRET"), "")),
	)
	if err != nil {
		log.Fatal(err)
	}

	s3Client = s3.NewFromConfig(cfg)
}

func main() {
	r := gin.Default()

	r.POST("/upload", func(c *gin.Context) {
		file, _, err := c.Request.FormFile("files")
		if err != nil {
			log.Printf("failed to get file: %v\n", err)
			c.JSON(http.StatusBadRequest, "{\"message\": \"unknown\"}")
			return
		}
		buf, err := io.ReadAll(file)
		if err != nil {
			log.Printf("failed to read file: %v\n", err)
			c.JSON(http.StatusBadRequest, "{\"message\": \"failed to read file\"}")
			return
		}

		mime, ok := IsAllowedFileType(buf)
		if !ok {
			c.JSON(http.StatusForbidden, fmt.Sprintf("{\"message\": \"not allowed file type: %v\"}", mime))
			return
		}
		filename := uuid.NewString()
		log.Printf("filename: %v\n", filename)

		i250x250, err := ResizeAndConvertToWebp(buf, 250, 250)
		if err != nil {
			log.Printf("failed to process & convert image: %v\n", err)
			c.JSON(http.StatusInternalServerError, "{\"message\": \"failed to process & convert file\"}")
			return
		}

		filepath := fmt.Sprintf("250x250/%s", filename)

		_, err = s3Client.PutObject(context.Background(), &s3.PutObjectInput{
			Bucket: aws.String(os.Getenv("R2_BUCKET_NAME")),
			Key:    aws.String(filepath),
			Body:   bytes.NewBuffer(i250x250),
			//ContentType: aws.String("image/webp"),
		})

		if err != nil {
			log.Printf("failed to upload file to s3: %v", err)
			c.JSON(http.StatusInternalServerError, "{\"message\": \"failed to upload file\"}")
			return
		}

		c.JSON(http.StatusOK, fmt.Sprintf("{\"message\": \"success\", \"filename\": \"%s\"}", filename))
	})

	//r.POST("/upload/image", func(c *gin.Context) {
	//	var sess *session.Session
	//	go func() {
	//		sess = session.Must(session.NewSessionWithOptions(session.Options{
	//			Profile:           os.Getenv("AWS_PROFILE"),
	//			SharedConfigState: session.SharedConfigEnable,
	//			Config: aws.Config{
	//				Region: aws.String("ap-northeast-1"),
	//			},
	//		}))
	//	}()
	//
	//	file, _, err := c.Request.FormFile("files")
	//	if err != nil {
	//		log.Printf("failed to get file: %v\n", err)
	//		c.JSON(http.StatusBadRequest, "{\"message\": \"unknown\"}")
	//		return
	//	}
	//	buf, err := io.ReadAll(file)
	//	if err != nil {
	//		log.Printf("failed to read file: %v\n", err)
	//		c.JSON(http.StatusBadRequest, "{\"message\": \"failed to read file\"}")
	//		return
	//	}
	//
	//	mime, ok := IsAllowedFileType(buf)
	//	if !ok {
	//		c.JSON(http.StatusForbidden, fmt.Sprintf("{\"message\": \"not allowed file type: %v\"}", mime))
	//		return
	//	}
	//	filename := uuid.NewString()
	//	log.Printf("filename: %v\n", filename)
	//	uploader := s3manager.NewUploader(sess)
	//
	//	iQuo80, err := CompressNAndConvertWebp(buf, 80)
	//	if err != nil {
	//		log.Printf("failed to compress: %v\n", err)
	//		c.JSON(http.StatusInternalServerError, "{\"message\": \"failed to compress file\"}")
	//		return
	//	}
	//	quo80Filepath := fmt.Sprintf("q80/%s", filename)
	//	_, err = uploader.Upload(&s3manager.UploadInput{
	//		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
	//		Body:   aws.ReadSeekCloser(bytes.NewBuffer(iQuo80)),
	//		Key:    aws.String(quo80Filepath),
	//	})
	//	if err != nil {
	//		log.Printf("failed to upload file to s3: %v", err)
	//		c.JSON(http.StatusInternalServerError, "{\"message\": \"failed to upload file\"}")
	//		return
	//	}
	//
	//	thumbnail, err := ResizeAndConvertToWebp(buf, 640, 480)
	//	if err != nil {
	//		log.Printf("failed to resize & convert: %v\n", err)
	//		c.JSON(http.StatusInternalServerError, "{\"message\": \"failed to compress file\"}")
	//		return
	//	}
	//	thumbnailFilePath := fmt.Sprintf("thumb/%s", filename)
	//	_, err = uploader.Upload(&s3manager.UploadInput{
	//		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
	//		Body:   aws.ReadSeekCloser(bytes.NewBuffer(thumbnail)),
	//		Key:    aws.String(thumbnailFilePath),
	//	})
	//	if err != nil {
	//		log.Printf("failed to upload file to s3: %v", err)
	//		c.JSON(http.StatusInternalServerError, "{\"message\": \"failed to upload file\"}")
	//		return
	//	}
	//
	//	c.JSON(http.StatusOK, fmt.Sprintf("{\"message\": \"success\", \"filename\": \"%s\"}", filename))
	//})

	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatalf("failed to run gin: %v", err)
	}
}
