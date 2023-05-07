package postgresql

import (
	"context"
	"fmt"

	"github.com/buivuanh/elotusteam-hackathon/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FileRepository struct{}

func NewFileRepository() *FileRepository {
	return &FileRepository{}
}

func (r *FileRepository) InsertImage(ctx context.Context, db *pgxpool.Pool, file *domain.Image) (*domain.Image, error) {
	query := `
		INSERT INTO images (file_path, original_name, content_type, byte_size, owner_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	err := db.QueryRow(ctx, query, file.FilePath, file.OriginalName, file.ContentType, file.ByteSize, file.OwnerID).Scan(&file.ID, &file.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert file: %w", err)
	}

	return file, nil
}
