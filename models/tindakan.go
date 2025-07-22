package models

type Tindakan struct {
	ID        uint   `json:"id"`
	KasusID   uint   `json:"kasus_id"`
	PetugasID uint   `json:"petugas_id"`
	ObatID    uint   `json:"obat_id"`
	Deskripsi string `json:"deskripsi"`
}
