package sacloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	libsacloud "github.com/yamamoto-febc/libsacloud/api"
	sakura "github.com/yamamoto-febc/libsacloud/resources"
	"time"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		Create: createServer,
		Read:   readServer,
		Delete: deleteServer,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func createServer(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*libsacloud.Client)

	name := d.Get("name").(string)
	password := d.Get("password").(string)

	spec := &sakura.Server{
		Name:        name,
		Description: "",
		ServerPlan: sakura.NumberResource{
			ID: 1001, // TODO
		},
		ConnectedSwitches: []map[string]string{{"Scope": "shared"}},
		Tags:              []string{"@virtio-net-pci"},
	}

	res, err := provider.Create(spec, "")
	if err != nil {
		return err
	}
	server := res.Server

	var diskPlan int64 = 4  //TODO
	diskSpec := &sakura.Disk{
		Name: name,
		Plan: sakura.NumberResource{
			ID: diskPlan,
		},
		SizeMB:     20480,    // TODO
		Connection: "virtio", //TODO
		SourceArchive: sakura.Resource{
			ID: "112800262964", //TODO
		},
	}

	diskID, err := provider.CreateDisk(diskSpec)
	if err != nil {
		return err
	}

	//wait for disk available
	waitForDiskAvailable(provider, diskID)

	//connect disk for server
	connectSuccess, err := provider.ConnectDisk(diskID, server.ID)
	if err != nil || !connectSuccess {
		return err
	}

	diskEditspec := &sakura.DiskEditValue{
		Password: password,
		//SSHKey: sakura.SSHKey{
		//	PublicKey: publicKey,
		//},
		//DisablePWAuth: !d.serverConfig.EnablePWAuth,
		//Notes:         notes[:],
	}

	editSuccess, err := provider.EditDisk(diskID, diskEditspec)
	if err != nil || !editSuccess {
		return err
	}
	//wait for disk available
	waitForDiskAvailable(provider, diskID)

	//start
	err = provider.PowerOn(server.ID)
	if err != nil {
		return err
	}
	//wait for startup
	waitForServerByState(provider, server.ID, "up")

	d.Set("name", name)
	d.SetId(server.ID)

	return readServer(d, meta)
}

func waitForDiskAvailable(provider *libsacloud.Client, diskID string) {
	for {
		s, err := provider.DiskState(diskID)
		if err != nil {
			continue
		}

		if s == "available" {
			break
		}
		time.Sleep(5 * time.Second)
	}
}

func waitForServerByState(provider *libsacloud.Client, serverID string, waitForState string) {
	for {
		s, err := getState(provider, serverID)
		if err != nil {
			continue
		}

		if s == waitForState {
			break
		}
		time.Sleep(5 * time.Second)
	}
}

func getState(provider *libsacloud.Client, serverID string) (string, error) {
	s, err := provider.State(serverID)
	if err != nil {
		return "", err
	}
	return s, nil
}

func readServer(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*libsacloud.Client)

	name := d.Get("name").(string)

	//サーバ検索
	server, err := provider.SearchServerByName(name)
	if err != nil {
		return err
	}

	d.Set("name", name)
	d.SetId(server.ID)

	return nil
}

func deleteServer(d *schema.ResourceData, meta interface{}) error {
	provider := meta.(*libsacloud.Client)

	name := d.Get("name").(string)

	//サーバ検索
	server, err := provider.SearchServerByName(name)
	if err != nil {
		return err
	}

	var disks []string

	for _, d := range server.Disks {
		disks = append(disks, d.ID)
	}
	err = provider.Delete(server.ID, disks)
	if err != nil {
		return err
	}

	return nil
}
