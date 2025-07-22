package services

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/voxtmault/psc/interfaces"
	"github.com/voxtmault/psc/models"
)

type KasusService struct {
	con *sql.DB
	cfg any
}

var _ interfaces.KasusInterface = &KasusService{}

func NewKasusService(dbCon *sql.DB, config any) *KasusService {
	return &KasusService{
		con: dbCon,
		cfg: config,
	}
}

func (s *KasusService) Get(ctx context.Context, filter *models.KasusFilter) (*models.Response, error) {
	// init variable yang diperlukan
	var res models.Response
	var totalSize int64
	var metadata models.PaginationMetadata
	arr := []*models.Kasus{}

	// get / open koneksi ke db

	// bikin query
	query := `
	SELECT * FROM kasus
	`
	mQuery := `SELECT COUNT(id) AS sample_size FROM KASUS`

	var sqlFilters []string
	var sqlVars []any
	if filter.JenisMasalahID != 0 {
		sqlFilters = append(sqlFilters, "jenis_masalah_id = ?")
		sqlVars = append(sqlVars, filter.JenisMasalahID)
	}
	if filter.PelaporID != 0 {
		sqlFilters = append(sqlFilters, "pelapor_id = ?")
		sqlVars = append(sqlVars, filter.PelaporID)
	}

	if len(sqlFilters) > 0 {
		query += " WHERE "
		mQuery += " WHERE "
		for _, item := range sqlFilters {
			query += item
			query += " AND"
			mQuery += item
			mQuery += " AND"
		}
		query = strings.TrimSuffix(query, " AND")
		mQuery = strings.TrimSuffix(mQuery, " AND")
	}

	// exec pagination metadata query
	rows := s.con.QueryRowContext(ctx, mQuery, sqlVars...)
	if err := rows.Scan(&totalSize); err != nil {
		res.StatusCode = http.StatusInternalServerError
		res.Message = "Internal Server Error"
		res.Data = "failed to get total sample size"
		return &res, err
	}

	query += "LIMIT ? OFFSET ?"
	offset := (filter.PageNumber - 1) * filter.PageSize
	sqlVars = append(sqlVars, filter.PageSize, offset)

	// exec query
	result, err := s.con.QueryContext(ctx, query, sqlVars...)
	if err != nil {
		res.StatusCode = http.StatusInternalServerError
		res.Message = "Internal Server Error"
		res.Data = nil
		return &res, err
	}
	defer result.Close()

	// parse query result
	for result.Next() {
		var obj models.Kasus
		if err := result.Scan(&obj.ID, &obj.JenisMasalahID); err != nil {
			res.StatusCode = http.StatusInternalServerError
			res.Message = "Internal Server Error"
			res.Data = nil
			return &res, err
		}

		arr = append(arr, &obj)
	}

	// proses tambahan (opsional)
	// contoh : get link ke file bukti, ke photo kondisi medis, perhitungan pagination
	if err := models.CalculateMetadata(ctx, totalSize, &metadata, &filter.PaginationFilter); err != nil {
		res.StatusCode = http.StatusInternalServerError
		res.Message = "Internal Server Error"
		res.Data = nil
		return &res, err
	}

	// bikin return object
	res.StatusCode = http.StatusOK
	res.Message = "Success"
	res.Data = arr
	res.Data = map[string]interface{}{
		"data":                arr,
		"pagination_metadata": metadata,
	}

	return &res, nil
}

func (s *KasusService) Create(ctx context.Context, payload *models.KasusCreate) (*models.Response, error) {
	var res models.Response

	tx, err := s.con.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		res.StatusCode = http.StatusInternalServerError
		res.Message = "Failed to begin TX"
		res.Data = nil
		return &res, nil
	}

	query := `
	INSERT INTO kasus (jenis_masalah_id, pasien_id, pelapor_id) VALUES (?,?,?)
	`
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		tx.Rollback()
		res.StatusCode = http.StatusInternalServerError
		res.Message = "Failed to commit TX"
		res.Data = nil
		return &res, nil
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, payload.JenisMasalahID, payload.PasienID, payload.PelaporID)
	if err != nil {
		tx.Rollback()
		res.StatusCode = http.StatusInternalServerError
		res.Message = "Failed to commit TX"
		res.Data = nil
		return &res, nil
	}

	lastId, _ := result.LastInsertId()

	if err := tx.Commit(); err != nil {
		res.StatusCode = http.StatusInternalServerError
		res.Message = "Failed to commit TX"
		res.Data = nil
		return &res, nil
	}

	res.StatusCode = http.StatusCreated
	res.Message = "Created"
	res.Data = lastId // opsional

	return &res, nil
}
