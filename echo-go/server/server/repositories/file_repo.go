package repositories

import (
	"context"
	"fmt"
	"time"

	"echo-react-serve/config"
	"echo-react-serve/helpers/minio"
	"echo-react-serve/server/models/dto"
	"echo-react-serve/server/models/entity"
)

// type FileRepo interface {
// 	DeleteFileByPath(ctx context.Context, path string) error
// }
// type fileRepo struct {
// 	db *gorm.DB
// }

// func NewFileRepo(db *gorm.DB) FileRepo {
// 	return &fileRepo{db: db}
// }

// func (r *fileRepo) DeleteFileByPath(ctx context.Context, path string) error {
// 	if err := r.db.Where("path = ?", path).Delete(&entity.File{}).Error; err != nil {
// 		return fmt.Errorf("failed to delete file from database: %v", err)
// 	}
// 	if err := minio.DeleteObject(ctx, config.MinioClient, path); err != nil {
// 		return fmt.Errorf("failed to delete file from storage: %v", err)
// 	}
// 	return nil
// }

// func updateFilePath(db *gorm.DB, bucket string, id int) (file entity.File, err error) {
// 	_ = db.First(&file, id).Error
// 	date := file.CreatedAt.Format("2006-01-02")
// 	path := fmt.Sprintf("%s/%s/%d/%s", bucket, date, file.RefID, file.Name)
// 	file.Path = path
// 	err = db.Save(&file).Error
// 	return
// }

func clearFiles(ctx context.Context, files []entity.File) error {
	for _, f := range files {
		if err := minio.DeleteObject(ctx, config.MinioClient, f.Path); err != nil {
			return fmt.Errorf("failed to delete file from storage: %v", err)
		}
		// has been done from each repositories that call this function
		// if err := db.Delete(&f).Error; err != nil {
		// 	return fmt.Errorf("failed to delete file from database: %v", err)
		// }
	}
	return nil
}

// save file to storage and return the file object
func saveFiles(ctx context.Context, parentSlug, bucket, dir string, files []dto.File) (entityFiles []entity.File, err error) {
	entityFiles = make([]entity.File, 0, len(files))
	for _, f := range files {
		if f.Size > 0 {
			file := entity.File{
				Name:      f.Filename,
				Size:      int(f.Size),
				Type:      f.Header.Get("Content-Type"),
				Path:      fmt.Sprintf("%s/%s/%s/%s", bucket, parentSlug, dir, f.Filename),
				CreatedAt: time.Now(),
			}
			if err = minio.PutObject(ctx, config.MinioClient, file.Path, f); err != nil {
				return
			}
			entityFiles = append(entityFiles, file)
		}
	}
	return
}
