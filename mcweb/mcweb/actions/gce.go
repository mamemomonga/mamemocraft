package actions

import (
	"context"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

// GCE GCE Control utility
type GCE struct {
	ctx      context.Context
	compute  *compute.Service
	instance *compute.InstancesService
	c *GCEConfig
}

type GCEConfig struct {
	CredentialsFile string
	Project  string
	Zone     string
	Instance string
}

// NewGCE New GCE Control utility
func NewGCE(c *GCEConfig) *GCE {
	t := new(GCE)
	t.c = c
	t.ctx = context.Background()
	return t
}

// LoadCredentialsFile Load Service Account Credentials File
func (t *GCE) LoadCredentialsFile() (err error) {
	t.compute, err = compute.NewService(t.ctx, option.WithCredentialsFile(t.c.CredentialsFile))
	if err != nil {
		return
	}
	t.instance = compute.NewInstancesService(t.compute)
	return
}

// List Instance List
func (t *GCE) List() (*compute.InstanceList, error) {
	return t.instance.List(t.c.Project, t.c.Zone).Do()
}

// Get Get Instance List
func (t *GCE) Get() (*compute.Instance, error) {
	return t.instance.Get(t.c.Project, t.c.Zone, t.c.Instance).Do()
}

// Start Instance
func (t *GCE) Start() (*compute.Operation, error) {
	return t.instance.Start(t.c.Project, t.c.Zone, t.c.Instance).Do()
}

// Stop Instance
func (t *GCE) Stop() (*compute.Operation, error) {
	return t.instance.Stop(t.c.Project, t.c.Zone, t.c.Instance).Do()
}

