package storage

type (
	jobInfo struct {
		Id                  uint
		Name                string
		Location            *City
		Description         string
		BasicPaymentPerHour uint
	}
)
