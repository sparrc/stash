package config

type AmazonDestination struct {
	DestinationName string
	AWSAccessKey    string
	AWSAccessKeyID  string
}

func (self *AmazonDestination) Name() string {
	return self.DestinationName
}

func (self *AmazonDestination) Folders() []string {
	return []string{"/tmp/foo", "/tmp/bar"}
}

func (self *AmazonDestination) Type() string {
	return "AmazonDestination"
}

func (self *AmazonDestination) Credentials() map[string]string {
	return map[string]string{
		"key":   self.AWSAccessKey,
		"keyID": self.AWSAccessKeyID,
	}
}
