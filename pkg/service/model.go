package service

type TronAccountResource struct {
	FreeNetUsed          int64 `json:"freeNetUsed"`
	FreeNetLimit         int64 `json:"freeNetLimit"`
	NetUsed              int64 `json:"netUsed"`
	NetLimit             int64 `json:"netLimit"`
	AssetNetUsed         int64 `json:"assetNetUsed"`
	AssetNetLimit        int64 `json:"assetNetLimit"`
	TotalNetLimit        int64 `json:"totalNetLimit"`
	TotalNetWeight       int64 `json:"totalNetWeight"`
	TotalTronPowerWeight int64 `json:"totalTronPowerWeight"`
	TronPowerUsed        int64 `json:"tronPowerUsed"`
	TronPowerLimit       int64 `json:"tronPowerLimit"`
	EnergyUsed           int64 `json:"energyUsed"`
	EnergyLimit          int64 `json:"energyLimit"`
	TotalEnergyLimit     int64 `json:"totalEnergyLimit"`
	TotalEnergyWeight    int64 `json:"totalEnergyWeight"`
	StorageUsed          int64 `json:"storageUsed"`
	StorageLimit         int64 `json:"storageLimit"`
}

type UserInfo struct {
	Id       string
	Username string
}
