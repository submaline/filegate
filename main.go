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

type JsonResponse struct {
	Message       string `json:"message"`
	FileName      string `json:"file_name"`
	FilePath      string `json:"file_path"`
	ThumbnailPath string `json:"thumbnail_path"`
}

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

	r.POST("/upload/icon", func(c *gin.Context) {
		file, _, err := c.Request.FormFile("files")
		if err != nil {
			log.Printf("failed to get file: %v\n", err)

			c.JSON(http.StatusBadRequest, JsonResponse{Message: "unknown"})
			return
		}
		buf, err := io.ReadAll(file)
		if err != nil {
			log.Printf("failed to read file: %v\n", err)
			c.JSON(http.StatusBadRequest, JsonResponse{Message: "failed to read file"})
			return
		}

		mime, ok := IsAllowedFileType(buf)
		if !ok {
			c.JSON(http.StatusForbidden, JsonResponse{Message: fmt.Sprintf("not allowed file type: %v", mime)})
			return
		}
		filename := uuid.NewString()
		log.Printf("filename: %v\n", filename)

		i250x250, err := ResizeAndConvertToWebp(buf, 250, 250)
		if err != nil {
			log.Printf("failed to process & convert image: %v\n", err)
			c.JSON(http.StatusInternalServerError, JsonResponse{Message: "failed to process & convert file"})
			return
		}

		filepath := fmt.Sprintf("250x250/%s", filename)

		_, err = s3Client.PutObject(context.Background(), &s3.PutObjectInput{
			Bucket: aws.String(os.Getenv("R2_BUCKET_NAME")),
			Key:    aws.String(filepath),
			Body:   bytes.NewBuffer(i250x250),
		})

		if err != nil {
			log.Printf("failed to upload file to s3: %v", err)
			c.JSON(http.StatusInternalServerError, JsonResponse{Message: "failed to upload file"})
			return
		}

		c.JSON(http.StatusOK, JsonResponse{
			Message:  "success",
			FileName: filename,
			FilePath: filepath,
		})
	})

	r.POST("/upload/image", func(c *gin.Context) {
		file, _, err := c.Request.FormFile("files")
		if err != nil {
			log.Printf("failed to get file: %v\n", err)
			c.JSON(http.StatusBadRequest, JsonResponse{Message: "unknown"})
			return
		}
		buf, err := io.ReadAll(file)
		if err != nil {
			log.Printf("failed to read file: %v\n", err)
			c.JSON(http.StatusBadRequest, JsonResponse{Message: "failed to read file"})
			return
		}

		mime, ok := IsAllowedFileType(buf)
		if !ok {
			c.JSON(http.StatusForbidden, JsonResponse{Message: fmt.Sprintf("not allowed file type: %v", mime)})
			return
		}
		filename := uuid.NewString()
		log.Printf("filename: %v\n", filename)

		var filepath string
		go func() { // 80%クオリティ
			compressedImg, err := CompressNAndConvertWebp(buf, 80)
			if err != nil {
				log.Printf("failed to process & convert image: %v\n", err)
				c.JSON(http.StatusInternalServerError, JsonResponse{Message: "failed to process & convert file"})
				return
			}
			filepath = fmt.Sprintf("q80/%s", filename)
			//
			_, err = s3Client.PutObject(context.Background(), &s3.PutObjectInput{
				Bucket: aws.String(os.Getenv("R2_BUCKET_NAME")),
				Key:    aws.String(filepath),
				Body:   bytes.NewBuffer(compressedImg),
			})
			//
			if err != nil {
				log.Printf("failed to upload file to s3: %v", err)
				c.JSON(http.StatusInternalServerError, JsonResponse{Message: "failed to upload file"})
				return
			}
		}()

		// サムネイル
		thumb, err := ThumbnailWebp(buf)
		if err != nil {
			log.Printf("failed to process & convert image: %v\n", err)
			c.JSON(http.StatusInternalServerError, JsonResponse{Message: "failed to process & convert file"})
			return
		}
		_, err = s3Client.PutObject(context.Background(), &s3.PutObjectInput{
			Bucket: aws.String(os.Getenv("R2_BUCKET_NAME")),
			Key:    aws.String("thumbnail/" + filename),
			Body:   bytes.NewBuffer(thumb),
		})
		if err != nil {
			log.Printf("failed to upload file to s3: %v", err)
			c.JSON(http.StatusInternalServerError, JsonResponse{Message: "failed to upload file"})
			return
		}

		c.JSON(http.StatusOK, JsonResponse{
			Message:       "success",
			FileName:      filename,
			FilePath:      filepath,
			ThumbnailPath: "thumbnail/" + filename,
		})
	})

	if err := r.Run("0.0.0.0:8080"); err != nil {
		log.Fatalf("failed to run gin: %v", err)
	}
}
