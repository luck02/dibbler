package logs

type logs interface {
	LogError(interface{})
	LogRequestToBid(guid string)
}
