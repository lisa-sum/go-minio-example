package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"path/filepath"
)

func InitMinio() *minio.Client {
	// endpoint := "47.120.5.83:5001/api/v1/service-account-credentials"
	// endpoint := "47.120.5.83:5001"
	// endpoint := "127.0.0.1:9000" // 接入点
	// endpoint := "47.120.5.83:9001" // 接入点
	// endpoint := "47.120.5.83:5000/api/v1/service-account-credentials" // 接入点
	endpoint := "47.120.5.83:5000" // 接入点
	accessKeyID := "iKnOlz40VQPXVmL7Ranr"
	secretAccessKey := "x5y8VF0A2TlLUlAOj2yTNRyX90cm667zyg1jC9xu"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln("InitMinio err", err)
	}

	return minioClient
}

func main() {
	r := gin.Default()
	minioClient := InitMinio()

	r.POST("/douyin/publish/action", func(c *gin.Context) {
		title := c.DefaultPostForm("title", "title")
		video, err := c.FormFile("data")
		if err != nil {
			panic(err)
		}

		baseDir := "upload/video"
		videoPath := filepath.Join(baseDir, fmt.Sprintf("%s.mp4", title))
		videoName := fmt.Sprintf("%s.mp4", title)
		fmt.Printf("videoName%s", videoName)
		if err := c.SaveUploadedFile(video, videoPath); err != nil {
			c.JSON(500, gin.H{
				"msg": err,
			})
		}

		//  桶名
		bucketName := "tiktok"
		// location := "us-east-1"
		contextType := "video/mp4"
		// filePath := filepath.Join("video", fmt.Sprintf("%s.mp4", title))
		// contextType := "binary/octet-stream"
		FileUploader(context.Background(), minioClient, bucketName, videoName, videoPath, contextType)

	})

	r.Run(":8080")
}

func FileUploader(ctx context.Context, client *minio.Client, bucketName, objectName, filePath, contextType string) {
	// bucketName := "mymusic"
	// objectName := "audit.log"
	// filePath := "./audit.log"
	// contextType := "application/text"

	object, err := client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contextType})
	if err != nil {
		log.Println("上传失败：", err)
	}
	log.Printf("Successfully uploaded %s of size %d\n", objectName, object.Size)
}
