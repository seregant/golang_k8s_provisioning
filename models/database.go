package models

var prefix = "dbs_"

type Pengguna struct {
	IDPengguna  int    `gorm:"column:pengguna_id;type:int(8);primary_key:yes;auto_increment"`
	Nama        string `gorm:"column:pengguna_nama;type:char(30)"`
	Alamat      string `gorm:"column:pengguna_alamat;type:varchar(100)"`
	Email       string `gorm:"column:pengguna_email;type:varchar(50)"`
	Username    string `gorm:"column:pengguna_username;type:char(12)"`
	Password    string `gorm:"column:pengguna_password;type:varchar(12)"`
	DBname      string `gorm:"column:db_name;type:char(12)"`
	DBuser      string `gorm:"column:db_username;type:char(12)"`
	DBpass      string `gorm:"column:db_pass;type:char(12)"`
	StorageSize int    `gorm:"column:storage_size;type:int(3)"`
	OcUrl       string `gorm:"column:oc_url;type:varchar(100)"`
	IsAdmin     int    `gorm:"column:is_admin;type:tinyint(1)"`
}

func (Pengguna) TableName() string {
	return prefix + "pengguna"
}
