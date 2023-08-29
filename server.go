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

type UploadFileResponse struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

func InitMinio() *minio.Client {
	endpoint := "127.0.0.1:5000" // 接入点
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
	minioClient := InitMinio() // Minio Client

	r.POST("/douyin/publish/action", func(c *gin.Context) {
		// 视频标题, 作为视频名
		title := c.DefaultPostForm("title", "title")
		// 视频流
		video, err := c.FormFile("data")
		if err != nil {
			panic(err)
		}

		// 保存在服务器的目录, 此处为根目录的upload/video
		baseDir := "upload/video"
		// 视频的完整路径
		videoPath := filepath.Join(baseDir, fmt.Sprintf("%s.mp4", title))
		// 视频名, 视频 + 类型后缀
		videoName := fmt.Sprintf("%s.mp4", title)
		fmt.Printf("videoName%s", videoName)
		// 在服务器保存视频流为文件
		if err := c.SaveUploadedFile(video, videoPath); err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"msg": err,
			})
			return
		}

		bucketName := "tiktok" //  桶名
		// location := "us-east-1"
		contextType := "video/mp4" // 文件类型, 此处为视频
		// contextType := "binary/octet-stream"
		uploader, err := FileUploader(context.Background(), minioClient, bucketName, videoName, videoPath, contextType)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{
				"msg": err,
			})
			return
		}

		c.JSON(200, gin.H{
			"msg":  "OK",
			"data": uploader,
		})
	})

	r.Run(":8080")
}

func FileUploader(ctx context.Context, client *minio.Client, bucketName, objectName, filePath, contextType string) (*UploadFileResponse, error) {
	object, err := client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contextType})
	if err != nil {
		log.Println("上传失败：", err)
		return nil, err
	}
	log.Printf("Successfully uploaded %s of size %d\n", objectName, object.Size)

	return &UploadFileResponse{
		Name: objectName,
		Size: object.Size,
	}, nil
}
