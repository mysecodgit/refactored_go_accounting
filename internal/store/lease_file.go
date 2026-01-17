package store

import (
	"context"
	"database/sql"
)

type LeaseFile struct {
	ID int64 `json:"id"`

	LeaseID int64 `json:"lease_id"`

	Filename     string `json:"filename"`
	OriginalName string `json:"original_name"`
	FilePath     string `json:"file_path"`
	FileType     string `json:"file_type"`
	FileSize     int64  `json:"file_size"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type LeaseFileStore struct {
	db *sql.DB
}

func NewLeaseFileStore(db *sql.DB) *LeaseFileStore {
	return &LeaseFileStore{db: db}
}

func (s *LeaseFileStore) GetByLeaseID(
	ctx context.Context,
	leaseID int64,
) ([]LeaseFile, error) {

	query := `
		SELECT
			id, lease_id,
			filename, original_name, file_path,
			file_type, file_size,
			created_at, updated_at
		FROM lease_files
		WHERE lease_id = ?
		ORDER BY created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, query, leaseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []LeaseFile

	for rows.Next() {
		var f LeaseFile
		if err := rows.Scan(
			&f.ID,
			&f.LeaseID,
			&f.Filename,
			&f.OriginalName,
			&f.FilePath,
			&f.FileType,
			&f.FileSize,
			&f.CreatedAt,
			&f.UpdatedAt,
		); err != nil {
			return nil, err
		}
		files = append(files, f)
	}

	return files, nil
}

func (s *LeaseFileStore) GetByID(ctx context.Context, id int64) (*LeaseFile, error) {
	query := `
		SELECT
			id, lease_id,
			filename, original_name, file_path,
			file_type, file_size,
			created_at, updated_at
		FROM lease_files
		WHERE id = ?
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	var f LeaseFile
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&f.ID,
		&f.LeaseID,
		&f.Filename,
		&f.OriginalName,
		&f.FilePath,
		&f.FileType,
		&f.FileSize,
		&f.CreatedAt,
		&f.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &f, nil
}

func (s *LeaseFileStore) Create(
	ctx context.Context,
	tx *sql.Tx,
	f *LeaseFile,
) (*int64, error) {

	query := `
		INSERT INTO lease_files
		(lease_id, filename, original_name, file_path, file_type, file_size)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := tx.ExecContext(
		ctx,
		query,
		f.LeaseID,
		f.Filename,
		f.OriginalName,
		f.FilePath,
		f.FileType,
		f.FileSize,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	f.ID = id
	return &id, nil
}

func (s *LeaseFileStore) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM lease_files WHERE id = ?`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeOutDuration)
	defer cancel()

	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}
