package cloud

import "context"
import "os"
import "github.com/digitalocean/godo"

type DigitalOceanProvider struct {
	client *godo.Client
}

var _ CloudProvider = &DigitalOceanProvider{}

func NewDigitalOceanProvider() *DigitalOceanProvider {
	p := &DigitalOceanProvider{}
	token, _ := os.ReadFile("do_token")
	client := godo.NewFromToken(string(token))

	p.client = client

	return p
}

func (p *DigitalOceanProvider) GetStatus() WorkstationStatus {
	droplets, _, _ := p.client.Droplets.List(context.TODO(), nil)

	return WorkstationStatus{
		Name: droplets[0].Name,
	}
}
