package models

type PenggunaRes struct {
	IDPengguna  string `form:"id" json:"id"`
	Nama        string `form:"nama" json:"nama"`
	Alamat      string `form:"alamat" json:"alamat"`
	Email       string `form:"email" json:"email"`
	Username    string `form:"username" json:"username"`
	Password    string `form:"password" json:"password"`
	DBname      string `form:"db_name" json:"db_name"`
	DBuser      string `form:"db_username" json:"db_username"`
	DBpass      string `form:"db_pass" json:"db_pass"`
	ConfPath    string `form:"config_path" json:"config_path"`
	StorageSize string `form:"storage_size" json:"storage_size"`
}
