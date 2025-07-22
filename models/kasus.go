package models

import "time"

type Kasus struct {
	ID                  uint        `json:"id"`
	JenisMasalahID      uint        `json:"jenis_masalah_id"`
	PelaporID           uint        `json:"pelapor_id"`
	PasienID            uint        `json:"pasien_id"`
	KategoriPelaporanID uint        `json:"kategori_pelaporan_id"`
	StatusKasusID       uint        `json:"status_kasus_id"`
	Lokasi              string      `json:"lokasi"`
	TanggalDibuat       time.Time   `json:"tanggal_dibuat"`
	Keterangan          string      `json:"keterangan"`
	Tindakan            []*Tindakan `json:"tindakan"`
}

type KasusCreate struct {
	JenisMasalahID      uint   `json:"jenis_masalah_id"`
	PelaporID           uint   `json:"pelapor_id"`
	PasienID            uint   `json:"pasien_id"`
	KategoriPelaporanID uint   `json:"kategori_pelaporan_id"`
	Lokasi              string `json:"lokasi"`
	Keterangan          string `json:"keterangan"`
}

type KasusUpdate struct {
	ID             uint        `json:"id"`
	JenisMasalahID uint        `json:"jenis_masalah_id"`
	PasienID       uint        `json:"pasien_id"`
	StatusKasusID  uint        `json:"status_kasus_id"`
	Lokasi         string      `json:"lokasi"`
	Keterangan     string      `json:"keterangan"`
	Tindakan       []*Tindakan `json:"tindakan"`
}

type KasusDelete struct {
	ID        uint `json:"id"`
	DeletedBy uint `json:"deleted_by"` // bisa diambil dari jwt
}

// Untuk get request
type KasusFilter struct {
	PaginationFilter

	JenisMasalahID      uint   `query:"jenis_masalah_id" validate:"omitempty,gte=1"`
	PelaporID           uint   `query:"pelapor_id" validate:"omitempty,uuid4"`
	PasienID            uint   `query:"pasien_id"`
	KategoriPelaporanID uint   `query:"kategori_pelaporan_id"`
	StatusKasusID       uint   `query:"status_kasus_id"`
	TanggalDibuat       string `query:"tanggal_dibuat" validate:"datetime=2006-01-02 15:04:05"`
}
