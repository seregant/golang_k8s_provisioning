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
	ClusterConf string `gorm:"column:config_path;type:varchar(100)"`
	StorageSize int    `gorm:"column:storage_size;type:int(3)"`
	OcUrl       string `gorm:"column:oc_url;type:varchar(100)"`
}

func (Pengguna) TableName() string {
	return prefix + "pengguna"
}

type ClusterNode struct {
	IDWorker int `gorm:"column:worker_id;type:int(6);primary_key:yes;auto_increment"`
	CpuCores int `gorm:"column:worker_cpu;type:int(2)"`
	Ram      int `gorm:"column:worker_ram;type:int(4)"`
	Dsik     int `gorm:"column:worker_disk;type:int(15)"`
}

func (ClusterNode) TableName() string {
	return prefix + "cluster_node"
}

type ClusterLoad struct {
	DeploymentName string `gorm:"column:deployment_name;type:varchar(20)"`
	CpuLoad        string `gorm:"column:load_cpu_percent;type:int(3)"`
	RamLoad        string `gorm:"column:load_ram_percent;type:int(3)"`
	DiskLoad       string `gorm:"column:load_disk_percent;type:int(3)"`
}

func (ClusterLoad) TableName() string {
	return prefix + "cluster_load"
}
