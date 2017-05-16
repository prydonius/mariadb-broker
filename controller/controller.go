package controller

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
	"github.com/kubernetes-incubator/service-catalog/contrib/pkg/broker/controller"
	"github.com/kubernetes-incubator/service-catalog/pkg/brokerapi"
	"github.com/prydonius/mariadb-broker/client"
)

type errNoSuchInstance struct {
	instanceID string
}

func (e errNoSuchInstance) Error() string {
	return fmt.Sprintf("no such instance with ID %s", e.instanceID)
}

type mariadbController struct {
}

// CreateController creates an instance of a User Provided service broker controller.
func CreateController() controller.Controller {
	return &mariadbController{}
}

func (c *mariadbController) Catalog() (*brokerapi.Catalog, error) {
	return &brokerapi.Catalog{
		Services: []*brokerapi.Service{
			{
				Name:        "mariadb",
				ID:          "3533e2f0-6335-4a4e-9d15-d7c0b90b75b5",
				Description: "MariaDB database",
				Plans: []brokerapi.ServicePlan{
					{
						Name:        "default",
						ID:          "b9600ecb-d511-4621-b450-a0fa1738e632",
						Description: "MariaDB database",
						Free:        true,
					},
				},
				Bindable: true,
			},
		},
	}, nil
}

func (c *mariadbController) CreateServiceInstance(id string, req *brokerapi.CreateServiceInstanceRequest) (*brokerapi.CreateServiceInstanceResponse, error) {
	if err := client.Install(releaseName(id), id); err != nil {
		return nil, err
	}
	glog.Infof("Created MariaDB Service Instance:\n%v\n", id)
	return &brokerapi.CreateServiceInstanceResponse{}, nil
}

func (c *mariadbController) GetServiceInstance(id string) (string, error) {
	return "", errors.New("Unimplemented")
}

func (c *mariadbController) RemoveServiceInstance(id string) (*brokerapi.DeleteServiceInstanceResponse, error) {
	if err := client.Delete(releaseName(id)); err != nil {
		return nil, err
	}
	return &brokerapi.DeleteServiceInstanceResponse{}, nil
}

func (c *mariadbController) Bind(instanceID, bindingID string, req *brokerapi.BindingRequest) (*brokerapi.CreateServiceBindingResponse, error) {
	host := releaseName(instanceID) + "-mariadb." + instanceID + ".svc.cluster.local"
	port := "3306"
	database := "dbname"
	username := "root"
	password, err := client.GetPassword(releaseName(instanceID), instanceID)
	if err != nil {
		return nil, err
	}
	return &brokerapi.CreateServiceBindingResponse{
		Credentials: brokerapi.Credential{
			"uri":      "mysql://" + username + ":" + password + "@" + host + ":" + port + "/" + database,
			"username": username,
			"password": password,
			"host":     host,
			"port":     port,
			"database": database,
		},
	}, nil
}

func (c *mariadbController) UnBind(instanceID string, bindingID string) error {
	// Since we don't persist the binding, there's nothing to do here.
	return nil
}

func releaseName(id string) string {
	return "i-" + id
}
