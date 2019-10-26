package models

type PenggunaRes struct {
	IDPengguna  int    `form:"id" json:"id"`
	Nama        string `form:"nama" json:"nama"`
	Alamat      string `form:"alamat" json:"alamat"`
	Email       string `form:"email" json:"email"`
	Username    string `form:"username" json:"username"`
	Password    string `form:"password" json:"password"`
	DBname      string `form:"db_name" json:"db_name"`
	DBuser      string `form:"db_username" json:"db_username"`
	DBpass      string `form:"db_pass" json:"db_pass"`
	StorageSize int    `form:"storage_size" json:"storage_size"`
	OcUrl       string `form:"oc_url" json:"oc_url"`
	IsAdmin     int    `form:"is_admin" json:"is_admin"`
}
