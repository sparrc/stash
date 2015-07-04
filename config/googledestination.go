package config

type GoogleDestination struct {
	DestinationName string
	Key             string
	KeyID           string
}

func (self *GoogleDestination) Name() string {
	return self.DestinationName
}
