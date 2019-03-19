package ccloud

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/sapcc/kubernikus/pkg/api/client/operations"
	"github.com/sapcc/kubernikus/pkg/api/models"
)

const (
	KlusterNameRegex = "^[a-z][-a-z0-9]{0,18}[a-z0-9]?$"
	PoolNameRegex    = "^[a-z][a-z0-9]{0,19}$"
)

func resourceCCloudKubernetesV1() *schema.Resource {
	return &schema.Resource{
		SchemaVersion: 1,

		Read:   resourceCCloudKubernetesV1Read,
		Update: resourceCCloudKubernetesV1Update,
		Create: resourceCCloudKubernetesV1Create,
		Delete: resourceCCloudKubernetesV1Delete,
		Importer: &schema.ResourceImporter{
			State: resourceCCloudKubernetesV1Import,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: kubernikusValidateClusterName,
			},

			"is_admin": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"advertise_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "1.1.1.1",
				ValidateFunc: validation.SingleIP(),
			},

			"cluster_cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "10.100.0.0/16",
				ValidateFunc: validation.CIDRNetwork(8, 16),
			},

			"service_cidr": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "198.18.128.0/17",
				ValidateFunc: validation.CIDRNetwork(8, 24),
			},

			"dns_address": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.SingleIP(),
			},

			"dns_domain": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "cluster.local",
			},

			"ssh_public_key": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"node_pools": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: kubernikusValidatePoolName,
						},
						"flavor": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"image": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "coreos-stable-amd64",
							ValidateFunc: validation.NoZeroValues,
						},
						"size": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 127),
						},
						"availability_zone": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.NoZeroValues,
						},
					},
				},
			},

			"openstack": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"lb_floating_network_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"lb_subnet_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"network_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"project_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"router_id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validation.NoZeroValues,
						},
						"security_group_name": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.NoZeroValues,
						},
					},
				},
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"phase": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"wormhole": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCCloudKubernetesV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	log.Printf("[KUBERNETES] Creating Kubernikus Kluster in project %s", config.TenantID)

	klient, err := config.kubernikusV1Client(GetRegion(d, config), d.Get("is_admin").(bool))
	if err != nil {
		return fmt.Errorf("Error creating Kubernikus client: %s", err)
	}

	cluster := &models.Kluster{
		Spec: models.KlusterSpec{
			NodePools: []models.NodePool{},
			Openstack: models.OpenstackSpec{},
		},
	}

	cluster.Name = d.Get("name").(string)
	cluster.Spec.AdvertiseAddress = d.Get("advertise_address").(string)
	cluster.Spec.ClusterCIDR = d.Get("cluster_cidr").(string)
	cluster.Spec.DNSAddress = d.Get("dns_address").(string)
	cluster.Spec.DNSDomain = d.Get("dns_domain").(string)
	cluster.Spec.SSHPublicKey = d.Get("ssh_public_key").(string)
	cluster.Spec.ServiceCIDR = d.Get("service_cidr").(string)
	cluster.Spec.NodePools, err = kubernikusExpandNodePoolsV1(d.Get("node_pools"))
	if err != nil {
		return err
	}
	if v := kubernikusExpandOpenstackSpecV1(d.Get("openstack")); v != nil {
		cluster.Spec.Openstack = *v
	}

	result, err := klient.CreateCluster(operations.NewCreateClusterParams().WithBody(cluster), klient.authFunc())
	if err != nil {
		return kubernikusHandleErrorV1("Error creating cluster", err)
	}

	pretty, _ := json.MarshalIndent(result.Payload, "", "  ")
	log.Printf("[DEBUG] Payload create: %s", string(pretty))

	d.SetId(cluster.Name)

	// waiting for Running state
	timeout := d.Timeout(schema.TimeoutCreate)
	target := "Running"
	pending := []string{"Pending", "Creating"}
	err = kubernikusWaitForClusterV1(klient, cluster.Name, target, pending, timeout)
	if err != nil {
		return kubernikusHandleErrorV1("Error waiting for running cluster state", err)
	}

	err = kubernikusWaitForNodePoolsV1(klient, cluster.Name, timeout)
	if err != nil {
		return kubernikusHandleErrorV1("Error waiting for cluster node pools active state", err)
	}

	return resourceCCloudKubernetesV1Read(d, meta)
}

func resourceCCloudKubernetesV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	log.Printf("[KUBERNETES] Reading Kubernikus Kluster in project %s", config.TenantID)

	klient, err := config.kubernikusV1Client(GetRegion(d, config), d.Get("is_admin").(bool))
	if err != nil {
		return fmt.Errorf("Error creating Kubernikus client: %s", err)
	}

	result, err := klient.ShowCluster(operations.NewShowClusterParams().WithName(d.Id()), klient.authFunc())
	if err != nil {
		switch err.(type) {
		case *operations.ShowClusterDefault:
			result := err.(*operations.ShowClusterDefault)

			if result.Payload.Message == "Not found" {
				d.SetId("")
				return nil
			}

			return fmt.Errorf("Error reading Kubernikus cluster: %s", result.Payload.Message)
		case error:
			return fmt.Errorf("Error reading Kubernikus cluster: %s", err)
		}
		return err
	}

	pretty, _ := json.MarshalIndent(result.Payload, "", "  ")
	log.Printf("[DEBUG] Payload get: %s", string(pretty))

	d.Set("advertise_address", result.Payload.Spec.AdvertiseAddress)
	d.Set("cluster_cidr", result.Payload.Spec.ClusterCIDR)
	d.Set("dns_address", result.Payload.Spec.DNSAddress)
	d.Set("dns_domain", result.Payload.Spec.DNSDomain)
	d.Set("name", result.Payload.Spec.Name)
	d.Set("ssh_public_key", result.Payload.Spec.SSHPublicKey)
	d.Set("service_cidr", result.Payload.Spec.ServiceCIDR)
	d.Set("version", result.Payload.Spec.Version)
	d.Set("phase", result.Payload.Status.Phase)
	d.Set("wormhole", result.Payload.Status.Wormhole)
	d.Set("openstack", kubernikusFlattenOpenstackSpecV1(&result.Payload.Spec.Openstack))
	d.Set("node_pools", kubernikusFlattenNodePoolsV1(result.Payload.Spec.NodePools))

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceCCloudKubernetesV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	log.Printf("[KUBERNETES] Updating Kubernikus Kluster in project %s", config.TenantID)

	klient, err := config.kubernikusV1Client(GetRegion(d, config), d.Get("is_admin").(bool))
	if err != nil {
		return fmt.Errorf("Error creating Kubernikus client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutUpdate)
	cluster := &models.Kluster{
		Spec: models.KlusterSpec{
			NodePools: []models.NodePool{},
			Openstack: models.OpenstackSpec{},
		},
	}

	cluster.Name = d.Id()

	if v, ok := d.GetOk("ssh_public_key"); ok {
		cluster.Spec.SSHPublicKey = v.(string)
	}

	if v, ok := d.GetOk("openstack.0.security_group_name"); ok {
		cluster.Spec.Openstack.SecurityGroupName = v.(string)
	}

	o, n := d.GetChange("node_pools")

	err = kubernikusUpdateNodePoolsV1(klient, cluster, o, n, timeout)
	if err != nil {
		return err
	}

	return resourceCCloudKubernetesV1Read(d, meta)
}

func resourceCCloudKubernetesV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	log.Printf("[KUBERNETES] Deleting Kubernikus Kluster in project %s", config.TenantID)

	klient, err := config.kubernikusV1Client(GetRegion(d, config), d.Get("is_admin").(bool))
	if err != nil {
		return fmt.Errorf("Error creating Kubernikus client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutDelete)

	_, err = klient.TerminateCluster(operations.NewTerminateClusterParams().WithName(d.Id()), klient.authFunc())
	if err != nil {
		return kubernikusHandleErrorV1("Error deleting cluster", err)
	}

	target := "Terminated"
	pending := []string{"Pending", "Creating", "Running", "Terminating"}
	err = kubernikusWaitForClusterV1(klient, d.Id(), target, pending, timeout)
	if err != nil {
		return kubernikusHandleErrorV1("Error waiting for cluster to be deleted", err)
	}

	return nil
}

func kubernikusValidateClusterName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(KlusterNameRegex).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q must be 1 to 20 characters with lowercase and uppercase letters, numbers and hyphens.", k))
	}
	return
}

func kubernikusValidatePoolName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)

	if !regexp.MustCompile(PoolNameRegex).MatchString(value) {
		errors = append(errors, fmt.Errorf(
			"%q must be 1 to 20 characters with lowercase and uppercase letters and numbers.", k))
	}
	return
}

func kubernikusFlattenOpenstackSpecV1(spec *models.OpenstackSpec) []map[string]interface{} {
	var res []map[string]interface{}

	if spec == (&models.OpenstackSpec{}) {
		return res
	}

	return append(res, map[string]interface{}{
		"lb_floating_network_id": spec.LBFloatingNetworkID,
		"lb_subnet_id":           spec.LBSubnetID,
		"network_id":             spec.NetworkID,
		"project_id":             spec.ProjectID,
		"router_id":              spec.RouterID,
		"security_group_name":    spec.SecurityGroupName,
	})
}

func kubernikusFlattenNodePoolsV1(nodePools []models.NodePool) []map[string]interface{} {
	var res []map[string]interface{}
	for _, p := range nodePools {
		res = append(res, map[string]interface{}{
			"availability_zone": p.AvailabilityZone,
			"flavor":            p.Flavor,
			"image":             p.Image,
			"name":              p.Name,
			"size":              p.Size,
		})
	}
	return res
}

func kubernikusExpandOpenstackSpecV1(raw interface{}) *models.OpenstackSpec {
	if raw != nil {
		if v, ok := raw.([]interface{}); ok {
			for _, v := range v {
				if v, ok := v.(map[string]interface{}); ok {
					res := new(models.OpenstackSpec)

					if v, ok := v["lb_floating_network_id"]; ok {
						res.LBFloatingNetworkID = v.(string)
					}
					if v, ok := v["lb_subnet_id"]; ok {
						res.LBSubnetID = v.(string)
					}
					if v, ok := v["network_id"]; ok {
						res.NetworkID = v.(string)
					}
					if v, ok := v["project_id"]; ok {
						res.ProjectID = v.(string)
					}
					if v, ok := v["router_id"]; ok {
						res.RouterID = v.(string)
					}
					if v, ok := v["security_group_name"]; ok {
						res.SecurityGroupName = v.(string)
					}

					return res
				}
			}
		}
	}

	return nil
}

func kubernikusExpandNodePoolsV1(raw interface{}) ([]models.NodePool, error) {
	var names []string

	if raw != nil {
		if v, ok := raw.([]interface{}); ok {
			var res []models.NodePool

			for _, v := range v {
				if v, ok := v.(map[string]interface{}); ok {
					var p models.NodePool

					if v, ok := v["name"]; ok {
						p.Name = v.(string)
						if strSliceContains(names, p.Name) {
							return nil, fmt.Errorf("Duplicate node pool name found: %s", p.Name)
						}
						names = append(names, p.Name)
					}
					if v, ok := v["flavor"]; ok {
						p.Flavor = v.(string)
					}
					if v, ok := v["image"]; ok {
						p.Image = v.(string)
					}
					if v, ok := v["size"]; ok {
						p.Size = int64(v.(int))
					}
					if v, ok := v["availability_zone"]; ok {
						p.AvailabilityZone = v.(string)
					}

					res = append(res, p)
				}
			}

			return res, nil
		}
	}

	return nil, nil
}

func kubernikusWaitForClusterV1(klient *Kubernikus, name string, target string, pending []string, timeout time.Duration) error {
	// Phase: "Pending","Creating","Running","Terminating"
	log.Printf("[DEBUG] Waiting for %s cluster to become %s.", name, target)

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    kubernikusKlusterV1GetPhase(klient, target, name),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForState()

	if err != nil {
		if e, ok := err.(*operations.ShowClusterDefault); ok && target == "Terminated" && e.Payload.Message == "Not found" {
			return nil
		}
	}

	return err
}

func kubernikusKlusterV1GetPhase(klient *Kubernikus, target string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		result, err := klient.ShowCluster(operations.NewShowClusterParams().WithName(name), klient.authFunc())
		if err != nil {
			return nil, "", err
		}

		pretty, _ := json.MarshalIndent(result.Payload, "", "  ")
		log.Printf("[DEBUG] Payload phase response: %s", string(pretty))

		if target != "Terminated" && result.Payload.Status.Phase == models.KlusterPhasePending {
			events, err := klient.GetClusterEvents(operations.NewGetClusterEventsParams().WithName(name), klient.authFunc())
			if err != nil {
				return nil, "", err
			}
			if len(events.Payload) > 0 {
				event := events.Payload[len(events.Payload)-1]
				if event.Reason == "ConfigurationError" {
					return nil, event.Reason, fmt.Errorf(event.Message)
				}
			}
		}
		return result.Payload, string(result.Payload.Status.Phase), nil
	}
}

func kubernikusWaitForNodePoolsV1(klient *Kubernikus, name string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for %s cluster node pools to become active.", name)

	stateConf := &resource.StateChangeConf{
		Target:     []string{"active"},
		Pending:    []string{"pending"},
		Refresh:    kubernikusKlusterV1GetNodePoolState(klient, name),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForState()

	return err
}

func kubernikusKlusterV1GetNodePoolState(klient *Kubernikus, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		result, err := klient.ShowCluster(operations.NewShowClusterParams().WithName(name), klient.authFunc())
		if err != nil {
			return nil, "", err
		}

		pretty, _ := json.MarshalIndent(result.Payload, "", "  ")
		log.Printf("[DEBUG] Payload state response: %s", string(pretty))

		if len(result.Payload.Spec.NodePools) != len(result.Payload.Status.NodePools) {
			return result.Payload, "pending", nil
		}
		for _, a := range result.Payload.Spec.NodePools {
			for _, s := range result.Payload.Status.NodePools {
				if a.Name == s.Name {
					// sometimes status size doesn't reflect the actual size, therefore we use "a.Size"
					if a.Size != s.Healthy {

						// check, whether there are error events
						events, err := klient.GetClusterEvents(operations.NewGetClusterEventsParams().WithName(name), klient.authFunc())
						if err != nil {
							return nil, "", err
						}
						if len(events.Payload) > 0 {
							event := events.Payload[len(events.Payload)-1]
							if strings.Contains(event.Reason, "Error") || strings.Contains(event.Reason, "Failed") {
								return nil, event.Reason, fmt.Errorf("%s node pool: %s", a.Name, event.Message)
							}
						}

						return result.Payload, "pending", nil
					}
				}
			}
		}
		return result.Payload, "active", nil
	}
}

func kubernikusUpdateNodePoolsV1(klient *Kubernikus, cluster *models.Kluster, oldNodePoolsRaw, newNodePoolsRaw interface{}, timeout time.Duration) error {
	var poolsToKeep []models.NodePool
	var poolsToDelete []models.NodePool
	oldNodePools, err := kubernikusExpandNodePoolsV1(oldNodePoolsRaw)
	if err != nil {
		return err
	}
	newNodePools, err := kubernikusExpandNodePoolsV1(newNodePoolsRaw)
	if err != nil {
		return err
	}

	pretty, _ := json.MarshalIndent(oldNodePools, "", "  ")
	log.Printf("[DEBUG] Old: %s", string(pretty))
	pretty, _ = json.MarshalIndent(newNodePools, "", "  ")
	log.Printf("[DEBUG] New: %s", string(pretty))

	// Determine if any node pools removed from the configuration.
	// Then downscale those pools and delete.
	for _, op := range oldNodePools {
		var found bool
		for _, np := range newNodePools {
			if op.Name == np.Name && op.Flavor == np.Flavor && op.Image == np.Image && (np.AvailabilityZone == "" || op.AvailabilityZone == np.AvailabilityZone) {
				tmp := np
				// copy previously "computed" AZ
				if np.AvailabilityZone == "" {
					tmp.AvailabilityZone = op.AvailabilityZone
				}
				poolsToKeep = append(poolsToKeep, np)
				found = true
			}
		}

		if !found {
			tmp := op
			tmp.Size = 0
			poolsToDelete = append(poolsToDelete, tmp)
		}
	}

	pretty, _ = json.MarshalIndent(poolsToKeep, "", "  ")
	log.Printf("[DEBUG] Keep: %s", string(pretty))
	pretty, _ = json.MarshalIndent(poolsToDelete, "", "  ")
	log.Printf("[DEBUG] Downscale: %s", string(pretty))

	if len(poolsToDelete) > 0 {
		// downscale
		cluster.Spec.NodePools = append(poolsToKeep, poolsToDelete...)
		err = kubernikusUpdateAndWait(klient, cluster, timeout)
		if err != nil {
			return err
		}
	}

	// delete old
	cluster.Spec.NodePools = poolsToKeep
	err = kubernikusUpdateAndWait(klient, cluster, timeout)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(poolsToKeep, newNodePools) {
		// create new
		cluster.Spec.NodePools = newNodePools
		err = kubernikusUpdateAndWait(klient, cluster, timeout)
		if err != nil {
			return err
		}
	}

	return nil
}

func kubernikusHandleErrorV1(msg string, err error) error {
	switch err.(type) {
	case *operations.TerminateClusterDefault:
		result := err.(*operations.TerminateClusterDefault)
		return fmt.Errorf("%s: %s", msg, result.Payload.Message)
	case error:
		return fmt.Errorf("%s: %s", msg, err)
	}
	return err
}

func kubernikusUpdateAndWait(klient *Kubernikus, cluster *models.Kluster, timeout time.Duration) error {
	pretty, _ := json.MarshalIndent(cluster, "", "  ")
	log.Printf("[DEBUG] Payload request: %s", string(pretty))

	result, err := klient.UpdateCluster(operations.NewUpdateClusterParams().WithName(cluster.Name).WithBody(cluster), klient.authFunc())
	if err != nil {
		return kubernikusHandleErrorV1("Error updating cluster", err)
	}
	err = kubernikusWaitForNodePoolsV1(klient, cluster.Name, timeout)
	if err != nil {
		return kubernikusHandleErrorV1("Error waiting for cluster node pools active state", err)
	}

	pretty, _ = json.MarshalIndent(result.Payload, "", "  ")
	log.Printf("[DEBUG] Payload response: %s", string(pretty))

	return nil
}

func resourceCCloudKubernetesV1Import(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)

	config := meta.(*Config)
	log.Printf("[KUBERNETES] Reading Kubernikus Kluster in project %s", config.TenantID)

	name := parts[0]
	var isAdmin bool
	var err error
	if len(parts) == 2 {
		isAdmin, err = strconv.ParseBool(parts[1])
		if err != nil {
			return nil, fmt.Errorf("Failed to parse is_admin field: %s", err)
		}
	}

	klient, err := config.kubernikusV1Client(GetRegion(d, config), isAdmin)
	if err != nil {
		return nil, fmt.Errorf("Error creating Kubernikus client: %s", err)
	}

	_, err = klient.ShowCluster(operations.NewShowClusterParams().WithName(name), klient.authFunc())
	if err != nil {
		return nil, kubernikusHandleErrorV1("Error reading Kubernikus cluster", err)
	}

	d.SetId(name)
	d.Set("is_admin", isAdmin)

	return []*schema.ResourceData{d}, nil
}
