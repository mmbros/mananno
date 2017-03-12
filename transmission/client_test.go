package transmission

import "testing"

var trClient *Client

func getTransmissionClient() *Client {
	if trClient == nil {
		trClient = NewClient("192.168.1.2:9091", "", "")
	}
	return trClient
}
func TestTorrentAdd(t *testing.T) {

	client := getTransmissionClient()
	res, err := client.TorrentAdd("h", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)

}
func TestSessionGet(t *testing.T) {
	client := getTransmissionClient()
	res, err := client.SessionGet()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Version: %s", res.Version)
}
func TestPing(t *testing.T) {
	client := getTransmissionClient()
	ok := client.Ping()
	if !ok {
		t.Fatal("Ping returned KO")
	}
}
