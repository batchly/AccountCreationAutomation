package model

type JsonRequestBody struct {
	Name                 string
	BaseRegion           string
	RoleARN              string
	MaxRIPercentageToUse float64
}
type CsvWriter struct {
	Accesskey           string
	Secretkey           string
	Region              string
	Account             string
	Domain              string
	RoleArn             string
	RoleArnStatus       string
	CmputeAccountStatus string
}
