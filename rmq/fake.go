package main

type FakeRMQClient struct {
	URI   string
	value string
}

func FakeFactory() RMQClient {
	return &FakeRMQClient{}
}

func (f *FakeRMQClient) Connect(uri, _channelName string) error {
	f.URI = uri
	return nil
}

func (f *FakeRMQClient) Send(value string) error {
	f.value = value
	return nil
}

func (f *FakeRMQClient) Receive() (string, error) {
	return f.value, nil
}

func (f *FakeRMQClient) Close() {}
