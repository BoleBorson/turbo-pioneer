package utils

func Rate(amount float64, time int) float64 {
	return (amount / float64(time)) * 60
}
