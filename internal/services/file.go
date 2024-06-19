package service

import (
	"cloudStorage/internal/dto"
	"cloudStorage/internal/middleware"
	"cloudStorage/internal/transport/rest/response"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/minio/minio-go/v7"
)

type FileLink struct {
	Name string `json:"name"`
	Link string `json:"link"`
	Type string `json:"type"`
}

func UploadFile(w http.ResponseWriter, r *http.Request, redis *redis.Client, minioStorage *minio.Client) {
	customRequest, isAuth := middleware.AuthMiddleware(w, r, redis)
	if !isAuth || customRequest == nil {
		response.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userId, ok := customRequest.Context().Value("userId").(string)
	if !ok {
		http.Error(w, "UserId not found in context", http.StatusInternalServerError)
		return
	}

	dirPath := r.FormValue("path")

	file, handler, err := r.FormFile("file")
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "Failed to get file from request")
		return
	}
	defer file.Close()

	objectName := handler.Filename
	filePath := "user-" + userId + "-files/"
	if dirPath != "" {
		filePath = filePath + dirPath + objectName
	} else {
		filePath = filePath + objectName
	}

	contentType := handler.Header.Get("Content-Type")

	_, err = minioStorage.PutObject(context.Background(), os.Getenv("BUCKET_NAME"), filePath, file, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, "Failed to upload file to storage")
		return
	}

	response.SendData(w, http.StatusCreated, "File uploaded successfully")
}

func FindFiles(w http.ResponseWriter, r *http.Request, redis *redis.Client, minioStorage *minio.Client) {
	customRequest, isAuth := middleware.AuthMiddleware(w, r, redis)
	if !isAuth || customRequest == nil {
		response.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userId, ok := customRequest.Context().Value("userId").(string)
	if !ok {
		response.SendError(w, http.StatusInternalServerError, "UserId not found in context")
		return
	}

	ctx := context.Background()
	prefix := "user-" + userId + "-files/"
	searchQuery := r.URL.Query().Get("search")

	objectCh := minioStorage.ListObjects(ctx, os.Getenv("BUCKET_NAME"), minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	})

	var fileLinks []FileLink

	for object := range objectCh {
		if object.Err != nil {
			log.Fatalln(object.Err)
		}

		if searchQuery != "" && !strings.Contains(object.Key, searchQuery) {
			continue
		}

		reqParams := make(url.Values)
		presignedURL, err := minioStorage.PresignedGetObject(ctx, os.Getenv("BUCKET_NAME"), object.Key, time.Hour*24, reqParams)
		if err != nil {
			response.SendError(w, http.StatusInternalServerError, "Failed to get presigned URL")
			return
		}

		fileLink := strings.Replace(presignedURL.String(), "minio:9000", os.Getenv("MINIO_PUBLIC_HOST"), 1)
		fileType := "file"
		if object.Size == 0 && strings.HasSuffix(object.Key, "/") {
			fileType = "folder"
		}
		fileLinks = append(fileLinks, FileLink{
			Name: strings.TrimPrefix(object.Key, prefix),
			Link: fileLink,
			Type: fileType,
		})
	}

	response.SendData(w, http.StatusOK, fileLinks)
}

func UpdateFile(w http.ResponseWriter, r *http.Request, redis *redis.Client, minioStorage *minio.Client) {
	customRequest, isAuth := middleware.AuthMiddleware(w, r, redis)
	if !isAuth || customRequest == nil {
		response.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userId, ok := customRequest.Context().Value("userId").(string)
	if !ok {
		response.SendError(w, http.StatusInternalServerError, "UserId not found in context")
		return
	}

	var body dto.UpdateFileDto
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		response.SendError(w, http.StatusBadRequest, "Invalid request body")
	}
	defer r.Body.Close()
	bucketName := os.Getenv("BUCKET_NAME")
	prefix := "/user-" + userId + "-files/"

	ctx := context.Background()

	originalPath := body.Path
	originalPathArr := strings.Split(originalPath, "/")
	// Имя файла с новым именем
	originalPathArr[len(originalPathArr)-1] = body.Name
	newPath := strings.Join(originalPathArr, "/")

	src := minio.CopySrcOptions{
		Bucket: bucketName,
		Object: prefix + originalPath,
	}
	dst := minio.CopyDestOptions{
		Bucket: bucketName,
		Object: prefix + newPath,
	}
	_, err = minioStorage.CopyObject(ctx, dst, src)
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, "Failed to copy object")
		return
	}

	err = minioStorage.RemoveObject(ctx, bucketName, prefix+originalPath, minio.RemoveObjectOptions{})
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, "Failed to remove original object")
	}

	response.SendData(w, http.StatusOK, newPath)
}

func CreateFolder(w http.ResponseWriter, r *http.Request, redis *redis.Client, minioStorage *minio.Client) {
	customRequest, isAuth := middleware.AuthMiddleware(w, r, redis)
	if !isAuth || customRequest == nil {
		response.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userId, ok := customRequest.Context().Value("userId").(string)
	if !ok {
		response.SendError(w, http.StatusInternalServerError, "UserId not found in context")
		return
	}

	dirPath := r.FormValue("path")
	pathFolder := "user-" + userId + "-files/" + dirPath

	if !strings.HasSuffix(pathFolder, "/") {
		pathFolder += "/"
	}

	uploadInfo, err := minioStorage.PutObject(context.Background(), os.Getenv("BUCKET_NAME"), pathFolder, nil, 0, minio.PutObjectOptions{})
	if err != nil {
		response.SendError(w, http.StatusInternalServerError, "Failed to create folder")
		return
	}

	response.SendData(w, http.StatusCreated, uploadInfo.Key)
}

func GetAllByPath(w http.ResponseWriter, r *http.Request, redis *redis.Client, minioStorage *minio.Client) {
	customRequest, isAuth := middleware.AuthMiddleware(w, r, redis)
	if !isAuth || customRequest == nil {
		response.SendError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	userId, ok := customRequest.Context().Value("userId").(string)
	if !ok {
		response.SendError(w, http.StatusInternalServerError, "UserId not found in context")
		return
	}
	dirPath := r.FormValue("path")
	pathFolder := "user-" + userId + "-files/" + dirPath

	ctx := context.Background()
	objectCh := minioStorage.ListObjects(ctx, os.Getenv("BUCKET_NAME"), minio.ListObjectsOptions{
		Prefix:    pathFolder,
		Recursive: false,
	})
	var files []FileLink
	for object := range objectCh {
		if object.Err != nil {
			log.Println("Error listing objects:", object.Err)
			http.Error(w, "Error listing objects", http.StatusInternalServerError)
			return
		}

		fileType := "file"
		if object.Size == 0 && strings.HasSuffix(object.Key, "/") {
			fileType = "folder"
		}

		files = append(files, FileLink{
			Name: strings.TrimPrefix(object.Key, pathFolder),
			Type: fileType,
		})
	}

	response.SendData(w, http.StatusOK, files)
}
