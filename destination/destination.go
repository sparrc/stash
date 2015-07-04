package destination

type Destination interface {
	Name() string
	// Key() string
	// KeyID() string
}

type AmazonDestination struct {
	DestinationName string
	AWSAccessKey    string
	AWSAccessKeyID  string
}

func (self *AmazonDestination) Name() string {
	return self.DestinationName
}

type GoogleDestination struct {
	DestinationName string
	Key             string
	KeyID           string
}

func (self *GoogleDestination) Name() string {
	return self.DestinationName
}
