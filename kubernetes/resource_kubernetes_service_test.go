package kubernetes

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	api "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestAccKubernetesService_basic(t *testing.T) {
	var conf api.Service
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		IDRefreshName:     "kubernetes_service.test",
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckKubernetesServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesServiceConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceExists("kubernetes_service.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.name", name),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.#", "1"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "spec.0.cluster_ip"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.name", ""),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.node_port", "0"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.port", "8080"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.target_port", "80"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.session_affinity", "None"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.type", "ClusterIP"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.publish_not_ready_addresses", "false"),
					testAccCheckServicePorts(&conf, []api.ServicePort{
						{
							Port:       int32(8080),
							Protocol:   api.ProtocolTCP,
							TargetPort: intstr.FromInt(80),
						},
					}),
				),
			},
			{
				ResourceName:            "kubernetes_service.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.0.resource_version", "wait_for_load_balancer"},
			},
			{
				Config: testAccKubernetesServiceConfig_modified(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceExists("kubernetes_service.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.name", name),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.#", "1"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "spec.0.cluster_ip"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.name", ""),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.node_port", "0"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.port", "8081"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.target_port", "80"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.session_affinity", "None"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.type", "ClusterIP"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.publish_not_ready_addresses", "true"),
					testAccCheckServicePorts(&conf, []api.ServicePort{
						{
							Port:       int32(8081),
							Protocol:   api.ProtocolTCP,
							TargetPort: intstr.FromInt(80),
						},
					}),
				),
			},
			{
				Config: testAccKubernetesServiceConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceExists("kubernetes_service.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.name", name),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.#", "1"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "spec.0.cluster_ip"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.name", ""),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.node_port", "0"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.port", "8080"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.target_port", "80"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.session_affinity", "None"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.type", "ClusterIP"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.publish_not_ready_addresses", "false"),
					testAccCheckServicePorts(&conf, []api.ServicePort{
						{
							Port:       int32(8080),
							Protocol:   api.ProtocolTCP,
							TargetPort: intstr.FromInt(80),
						},
					}),
				),
			},
		},
	})
}

func TestAccKubernetesService_loadBalancer(t *testing.T) {
	var conf api.Service
	name := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); skipIfNoLoadBalancersAvailable(t) },
		IDRefreshName:     "kubernetes_service.test",
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckKubernetesServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesServiceConfig_loadBalancer(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceExists("kubernetes_service.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.#", "1"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "spec.0.port.0.node_port"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.port", "8888"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.target_port", "80"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "spec.0.cluster_ip"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_ips.#", "2"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_ips.1", "10.0.0.4"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_ips.0", "10.0.0.3"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_name", "ext-name-"+name),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_traffic_policy", "Cluster"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.load_balancer_source_ranges.#", "2"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.load_balancer_source_ranges.0", "10.0.0.5/32"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.load_balancer_source_ranges.1", "10.0.0.6/32"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.selector.%", "1"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.selector.App", "MyApp"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.type", "LoadBalancer"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "status.0.load_balancer.0.ingress.0.ip"),
					testAccCheckServicePorts(&conf, []api.ServicePort{
						{
							Port:       int32(8888),
							Protocol:   api.ProtocolTCP,
							TargetPort: intstr.FromInt(80),
						},
					}),
				),
			},
			{
				Config: testAccKubernetesServiceConfig_loadBalancer_modified(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceExists("kubernetes_service.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.#", "1"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "spec.0.cluster_ip"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_ips.#", "2"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_ips.0", "10.0.0.4"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_ips.1", "10.0.0.5"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_name", "ext-name-modified-"+name),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_traffic_policy", "Local"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.load_balancer_source_ranges.#", "2"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.load_balancer_source_ranges.0", "10.0.0.1/32"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.load_balancer_source_ranges.1", "10.0.0.2/32"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.#", "1"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "spec.0.port.0.node_port"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.port", "9999"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.target_port", "81"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.selector.%", "2"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.selector.App", "MyModifiedApp"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.selector.NewSelector", "NewValue"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.type", "LoadBalancer"),
					testAccCheckServicePorts(&conf, []api.ServicePort{
						{
							Port:       int32(9999),
							Protocol:   api.ProtocolTCP,
							TargetPort: intstr.FromInt(81),
						},
					}),
				),
			},
		},
	})
}

func TestAccKubernetesService_loadBalancer_healthcheck(t *testing.T) {
	var conf api.Service
	name := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); skipIfNoLoadBalancersAvailable(t) },
		IDRefreshName:     "kubernetes_service.test",
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckKubernetesServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesServiceConfig_loadBalancer_healthcheck(name, 31111),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceExists("kubernetes_service.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_traffic_policy", "Local"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.type", "LoadBalancer"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.health_check_node_port", "31111"),
				),
			},
			{
				Config: testAccKubernetesServiceConfig_loadBalancer_healthcheck(name, 31112),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceExists("kubernetes_service.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_traffic_policy", "Local"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.type", "LoadBalancer"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.health_check_node_port", "31112"),
				),
			},
		},
	})
}

func TestAccKubernetesService_loadBalancer_annotations_aws(t *testing.T) {
	var conf api.Service
	name := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); skipIfNoLoadBalancersAvailable(t) },
		IDRefreshName:     "kubernetes_service.test",
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckKubernetesServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesServiceConfig_loadBalancer_annotations_aws(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceExists("kubernetes_service.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.annotations.%", "3"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.annotations.service.beta.kubernetes.io/aws-load-balancer-backend-protocol", "http"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.annotations.service.beta.kubernetes.io/aws-load-balancer-connection-idle-timeout", "300"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.annotations.service.beta.kubernetes.io/aws-load-balancer-ssl-ports", "*"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.#", "1"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "spec.0.port.0.node_port"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.port", "8888"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.target_port", "80"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "spec.0.cluster_ip"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_ips.#", "2"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_ips.1", "10.0.0.4"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_ips.0", "10.0.0.3"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_name", "ext-name-"+name),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.load_balancer_source_ranges.#", "2"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.load_balancer_source_ranges.0", "10.0.0.5/32"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.load_balancer_source_ranges.1", "10.0.0.6/32"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.selector.%", "1"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.selector.App", "MyApp"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.type", "LoadBalancer"),
					testAccCheckServicePorts(&conf, []api.ServicePort{
						{
							Port:       int32(8888),
							Protocol:   api.ProtocolTCP,
							TargetPort: intstr.FromInt(80),
						},
					}),
				),
			},
			{
				Config: testAccKubernetesServiceConfig_loadBalancer_annotations_aws_modified(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceExists("kubernetes_service.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.annotations.%", "4"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.annotations.service.beta.kubernetes.io/aws-load-balancer-backend-protocol", "http"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.annotations.service.beta.kubernetes.io/aws-load-balancer-connection-idle-timeout", "60"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.annotations.service.beta.kubernetes.io/aws-load-balancer-ssl-ports", "*"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.annotations.service.beta.kubernetes.io/aws-load-balancer-cross-zone-load-balancing-enabled", "true"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.#", "1"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "spec.0.cluster_ip"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_ips.#", "2"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_ips.0", "10.0.0.4"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_ips.1", "10.0.0.5"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_name", "ext-name-modified-"+name),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.load_balancer_source_ranges.#", "2"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.load_balancer_source_ranges.0", "10.0.0.1/32"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.load_balancer_source_ranges.1", "10.0.0.2/32"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.#", "1"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "spec.0.port.0.node_port"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.port", "9999"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.target_port", "81"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.selector.%", "2"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.selector.App", "MyModifiedApp"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.selector.NewSelector", "NewValue"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.type", "LoadBalancer"),
					testAccCheckServicePorts(&conf, []api.ServicePort{
						{
							Port:       int32(9999),
							Protocol:   api.ProtocolTCP,
							TargetPort: intstr.FromInt(81),
						},
					}),
				),
			},
		},
	})
}

func TestAccKubernetesService_nodePort(t *testing.T) {
	var conf api.Service
	name := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		IDRefreshName:     "kubernetes_service.test",
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckKubernetesServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesServiceConfig_nodePort(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceExists("kubernetes_service.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.#", "1"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "spec.0.cluster_ip"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_ips.#", "2"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_ips.0", "10.0.0.4"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_ips.1", "10.0.0.5"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_name", "ext-name-"+name),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.load_balancer_ip", "12.0.0.125"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.#", "2"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.name", "first"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "spec.0.port.0.node_port"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.port", "10222"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.target_port", "22"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.1.name", "second"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "spec.0.port.1.node_port"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.1.port", "10333"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.1.protocol", "TCP"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.1.target_port", "33"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.selector.%", "1"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.selector.App", "MyApp"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.session_affinity", "ClientIP"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.type", "NodePort"),
					testAccCheckServicePorts(&conf, []api.ServicePort{
						{
							Name:       "first",
							Port:       int32(10222),
							Protocol:   api.ProtocolTCP,
							TargetPort: intstr.FromInt(22),
						},
						{
							Name:       "second",
							Port:       int32(10333),
							Protocol:   api.ProtocolTCP,
							TargetPort: intstr.FromInt(33),
						},
					}),
				),
			},
		},
	})
}

func TestAccKubernetesService_noTargetPort(t *testing.T) {
	var conf api.Service
	name := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); skipIfNoLoadBalancersAvailable(t) },
		IDRefreshName:     "kubernetes_service.test",
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckKubernetesServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesServiceConfig_noTargetPort(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceExists("kubernetes_service.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.#", "1"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "spec.0.cluster_ip"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_ips.#", "0"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.#", "2"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.name", "http"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "spec.0.port.0.node_port"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.port", "80"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.target_port", "80"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.1.name", "https"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "spec.0.port.1.node_port"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.1.port", "443"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.1.protocol", "TCP"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.1.target_port", "443"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.selector.%", "1"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.selector.App", "MyOtherApp"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.session_affinity", "None"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.type", "LoadBalancer"),
					testAccCheckServicePorts(&conf, []api.ServicePort{
						{
							Name:       "http",
							Port:       int32(80),
							Protocol:   api.ProtocolTCP,
							TargetPort: intstr.FromInt(80),
						},
						{
							Name:       "https",
							Port:       int32(443),
							Protocol:   api.ProtocolTCP,
							TargetPort: intstr.FromInt(443),
						},
					}),
				),
			},
		},
	})
}

func TestAccKubernetesService_stringTargetPort(t *testing.T) {
	var conf api.Service
	name := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); skipIfNoLoadBalancersAvailable(t) },
		IDRefreshName:     "kubernetes_service.test",
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckKubernetesServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesServiceConfig_stringTargetPort(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceExists("kubernetes_service.test", &conf),
					testAccCheckServicePorts(&conf, []api.ServicePort{
						{
							Port:       int32(8080),
							Protocol:   api.ProtocolTCP,
							TargetPort: intstr.FromString("http-server"),
						},
					}),
				),
			},
		},
	})
}

func TestAccKubernetesService_externalName(t *testing.T) {
	var conf api.Service
	name := acctest.RandomWithPrefix("tf-acc-test")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		IDRefreshName:     "kubernetes_service.test",
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckKubernetesServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesServiceConfig_externalName(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceExists("kubernetes_service.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.cluster_ip", ""),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_ips.#", "0"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.external_name", "terraform.io"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.load_balancer_ip", ""),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.load_balancer_source_ranges.#", "0"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.#", "0"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.selector.%", "0"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.session_affinity", "None"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.type", "ExternalName"),
				),
			},
		},
	})
}

func TestAccKubernetesService_generatedName(t *testing.T) {
	var conf api.Service
	prefix := "tf-acc-test-gen-"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		IDRefreshName:     "kubernetes_service.test",
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckKubernetesServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesServiceConfig_generatedName(prefix),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceExists("kubernetes_service.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.annotations.%", "0"),
					//testAccCheckMetaAnnotations(&conf.ObjectMeta, map[string]string{}),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.labels.%", "0"),
					//testAccCheckMetaLabels(&conf.ObjectMeta, map[string]string{}),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.generate_name", prefix),
					resource.TestMatchResourceAttr("kubernetes_service.test", "metadata.0.name", regexp.MustCompile("^"+prefix)),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "metadata.0.uid"),
				),
			},
			{
				ResourceName:            "kubernetes_service.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.0.resource_version"},
			},
		},
	})
}

func TestAccKubernetesService_regression(t *testing.T) {
	var conf1, conf2 api.Service
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		IDRefreshName:     "kubernetes_service.test",
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      testAccCheckKubernetesServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: requiredProviders() + testAccKubernetesServiceConfig_regression("kubernetes-released", name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceExists("kubernetes_service.test", &conf1),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.name", name),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.#", "1"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "spec.0.cluster_ip"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.name", ""),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.node_port", "0"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.port", "8080"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.target_port", "80"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.session_affinity", "None"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.type", "ClusterIP"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.publish_not_ready_addresses", "false"),
					testAccCheckServicePorts(&conf1, []api.ServicePort{
						{
							Port:       int32(8080),
							Protocol:   api.ProtocolTCP,
							TargetPort: intstr.FromInt(80),
						},
					}),
				),
			},
			{
				Config: requiredProviders() + testAccKubernetesServiceConfig_regression("kubernetes-local", name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceExists("kubernetes_service.test", &conf2),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.name", name),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.#", "1"),
					resource.TestCheckResourceAttrSet("kubernetes_service.test", "spec.0.cluster_ip"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.name", ""),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.node_port", "0"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.port", "8080"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.protocol", "TCP"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.port.0.target_port", "80"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.session_affinity", "None"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.type", "ClusterIP"),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.publish_not_ready_addresses", "false"),
					testAccCheckKubernetesServiceForceNew(&conf1, &conf2, false),
					testAccCheckServicePorts(&conf2, []api.ServicePort{
						{
							Port:       int32(8080),
							Protocol:   api.ProtocolTCP,
							TargetPort: intstr.FromInt(80),
						},
					}),
				),
			},
		},
	})
}

func TestAccKubernetesService_stateUpgradeV0_loadBalancerIngress(t *testing.T) {
	var conf1, conf2 api.Service
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); skipIfNotRunningInEks(t) },
		IDRefreshName:     "kubernetes_service.test",
		ExternalProviders: testAccExternalProviders,
		CheckDestroy:      testAccCheckKubernetesServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: requiredProviders() + testAccKubernetesServiceConfig_stateUpgradev0("kubernetes-released", name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceExists("kubernetes_service.test", &conf1),
					resource.TestCheckResourceAttr("kubernetes_service.test", "metadata.0.name", name),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.type", "LoadBalancer"),
				),
			},
			{
				Config: requiredProviders() + testAccKubernetesServiceConfig_stateUpgradev0("kubernetes-local", name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesServiceExists("kubernetes_service.test", &conf2),
					resource.TestCheckResourceAttr("kubernetes_service.test", "spec.0.type", "LoadBalancer"),
					testAccCheckKubernetesServiceForceNew(&conf1, &conf2, false),
				),
			},
		},
	})
}

func testAccCheckKubernetesServiceForceNew(old, new *api.Service, wantNew bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if wantNew {
			if old.ObjectMeta.UID == new.ObjectMeta.UID {
				return fmt.Errorf("Expecting new resource for Service %s", old.ObjectMeta.UID)
			}
		} else {
			if old.ObjectMeta.UID != new.ObjectMeta.UID {
				return fmt.Errorf("Expecting Service UIDs to be the same: expected %s got %s", old.ObjectMeta.UID, new.ObjectMeta.UID)
			}
		}
		return nil
	}
}

func testAccCheckServicePorts(svc *api.Service, expected []api.ServicePort) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(expected) == 0 && len(svc.Spec.Ports) == 0 {
			return nil
		}

		ports := svc.Spec.Ports

		// Ignore NodePorts as these are assigned randomly
		for k := range ports {
			ports[k].NodePort = 0
		}

		if !reflect.DeepEqual(ports, expected) {
			return fmt.Errorf("Service ports don't match.\nExpected: %#v\nGiven: %#v",
				expected, svc.Spec.Ports)
		}

		return nil
	}
}

func testAccCheckKubernetesServiceDestroy(s *terraform.State) error {
	conn, err := testAccProvider.Meta().(KubeClientsets).MainClientset()

	if err != nil {
		return err
	}
	ctx := context.TODO()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kubernetes_service" {
			continue
		}

		namespace, name, err := idParts(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
		if err == nil {
			if resp.Name == rs.Primary.ID {
				return fmt.Errorf("Service still exists: %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckKubernetesServiceExists(n string, obj *api.Service) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn, err := testAccProvider.Meta().(KubeClientsets).MainClientset()
		if err != nil {
			return err
		}
		ctx := context.TODO()

		namespace, name, err := idParts(rs.Primary.ID)
		if err != nil {
			return err
		}

		out, err := conn.CoreV1().Services(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return err
		}

		*obj = *out
		return nil
	}
}

func testAccKubernetesServiceConfig_basic(name string) string {
	return fmt.Sprintf(`resource "kubernetes_service" "test" {
  metadata {
    annotations = {
      TestAnnotationOne = "one"
      TestAnnotationTwo = "two"
    }

    labels = {
      TestLabelOne   = "one"
      TestLabelTwo   = "two"
      TestLabelThree = "three"
    }

    name = "%s"
  }

  spec {
    port {
      port        = 8080
      target_port = 80
    }
  }
}
`, name)
}

func testAccKubernetesServiceConfig_regression(provider, name string) string {
	return fmt.Sprintf(`resource "kubernetes_service" "test" {
  provider = %s
  metadata {
    annotations = {
      TestAnnotationOne = "one"
      TestAnnotationTwo = "two"
    }

    labels = {
      TestLabelOne   = "one"
      TestLabelTwo   = "two"
      TestLabelThree = "three"
    }

    name = "%s"
  }

  spec {
    port {
      port        = 8080
      target_port = 80
    }
  }
}
`, provider, name)
}

func testAccKubernetesServiceConfig_stateUpgradev0(provider, name string) string {
	return fmt.Sprintf(`resource "kubernetes_service" "test" {
  provider = "%s"
  metadata {
    name = "%s"
  }

  spec {
    type = "LoadBalancer"
    port {
      port        = 8080
      target_port = 80
    }
  }
}
`, provider, name)
}

func testAccKubernetesServiceConfig_modified(name string) string {
	return fmt.Sprintf(`resource "kubernetes_service" "test" {
  metadata {
    annotations = {
      TestAnnotationOne = "one"
      Different         = "1234"
    }

    labels = {
      TestLabelOne   = "one"
      TestLabelThree = "three"
    }

    name = "%s"
  }

  spec {
    port {
      port        = 8081
      target_port = 80
    }

    publish_not_ready_addresses = "true"
  }
}
`, name)
}

func testAccKubernetesServiceConfig_loadBalancer(name string) string {
	return fmt.Sprintf(`resource "kubernetes_service" "test" {
  metadata {
    name = "%[1]s"
  }

  spec {
    external_name               = "ext-name-%[1]s"
    external_ips                = ["10.0.0.3", "10.0.0.4"]
    load_balancer_source_ranges = ["10.0.0.5/32", "10.0.0.6/32"]

    selector = {
      App = "MyApp"
    }

    port {
      port        = 8888
      target_port = 80
    }

    type = "LoadBalancer"
  }
}
`, name)
}

func testAccKubernetesServiceConfig_loadBalancer_modified(name string) string {
	return fmt.Sprintf(`resource "kubernetes_service" "test" {
  metadata {
    name = "%[1]s"
  }

  spec {
    external_name               = "ext-name-modified-%[1]s"
    external_ips                = ["10.0.0.4", "10.0.0.5"]
    load_balancer_source_ranges = ["10.0.0.1/32", "10.0.0.2/32"]
    external_traffic_policy     = "Local"

    selector = {
      App         = "MyModifiedApp"
      NewSelector = "NewValue"
    }

    port {
      port        = 9999
      target_port = 81
    }

    type = "LoadBalancer"
  }
}
`, name)
}

func testAccKubernetesServiceConfig_loadBalancer_annotations_aws(name string) string {
	return fmt.Sprintf(`resource "kubernetes_service" "test" {
  metadata {
    name = "%[1]s"
    annotations = {
      "service.beta.kubernetes.io/aws-load-balancer-backend-protocol"        = "http"
      "service.beta.kubernetes.io/aws-load-balancer-connection-idle-timeout" = "300"
      "service.beta.kubernetes.io/aws-load-balancer-ssl-ports"               = "*"
    }
  }

  spec {
    external_name               = "ext-name-%[1]s"
    external_ips                = ["10.0.0.3", "10.0.0.4"]
    load_balancer_source_ranges = ["10.0.0.5/32", "10.0.0.6/32"]

    selector = {
      App = "MyApp"
    }

    port {
      port        = 8888
      target_port = 80
    }

    type = "LoadBalancer"
  }
}
`, name)
}

func testAccKubernetesServiceConfig_loadBalancer_annotations_aws_modified(name string) string {
	return fmt.Sprintf(`resource "kubernetes_service" "test" {
  metadata {
    name = "%[1]s"
    annotations = {
      "service.beta.kubernetes.io/aws-load-balancer-backend-protocol"                  = "http"
      "service.beta.kubernetes.io/aws-load-balancer-connection-idle-timeout"           = "60"
    "service.beta.kubernetes.io/aws-load-balancer-ssl-ports"                         = "*"
    "service.beta.kubernetes.io/aws-load-balancer-cross-zone-load-balancing-enabled" = "true"
    }
  }

  spec {
    external_name               = "ext-name-modified-%[1]s"
    external_ips                = ["10.0.0.4", "10.0.0.5"]
    load_balancer_source_ranges = ["10.0.0.1/32", "10.0.0.2/32"]

    selector = {
      App         = "MyModifiedApp"
      NewSelector = "NewValue"
    }

    port {
      port        = 9999
      target_port = 81
    }

    type = "LoadBalancer"
  }
}
`, name)
}

func testAccKubernetesServiceConfig_loadBalancer_healthcheck(name string, nodePort int) string {
	return fmt.Sprintf(`resource "kubernetes_service" "test" {
  metadata {
    name = "%[1]s"
  }

  spec {
    external_name               = "ext-name-%[1]s"
    external_ips                = ["10.0.0.3", "10.0.0.4"]
    load_balancer_source_ranges = ["10.0.0.5/32", "10.0.0.6/32"]
    external_traffic_policy     = "Local"
    health_check_node_port      = %[2]d

    selector = {
      App = "MyApp"
    }

    port {
      port        = 8888
      target_port = 80
    }

    type = "LoadBalancer"
  }
}
`, name, nodePort)
}

func testAccKubernetesServiceConfig_nodePort(name string) string {
	return fmt.Sprintf(`resource "kubernetes_service" "test" {
  metadata {
    name = "%[1]s"
  }

  spec {
    external_name    = "ext-name-%[1]s"
    external_ips     = ["10.0.0.4", "10.0.0.5"]
    load_balancer_ip = "12.0.0.125"

    selector = {
      App = "MyApp"
    }

    session_affinity = "ClientIP"

    port {
      name        = "first"
      port        = 10222
      target_port = 22
    }

    port {
      name        = "second"
      port        = 10333
      target_port = 33
    }

    type = "NodePort"
  }
}
`, name)
}

func testAccKubernetesServiceConfig_stringTargetPort(name string) string {
	return fmt.Sprintf(`resource "kubernetes_service" "test" {
  metadata {
    name = "%s"

    labels = {
      app  = "helloweb"
      tier = "frontend"
    }
  }

  spec {
    type = "LoadBalancer"

    selector = {
      app  = "helloweb"
      tier = "frontend"
    }

    port {
      port        = 8080
      target_port = "http-server"
    }
  }
}
`, name)
}

func testAccKubernetesServiceConfig_noTargetPort(name string) string {
	return fmt.Sprintf(`resource "kubernetes_service" "test" {
  metadata {
    name = "%s"
  }

  spec {
    selector = {
      App = "MyOtherApp"
    }

    port {
      name = "http"
      port = 80
    }

    port {
      name = "https"
      port = 443
    }

    type = "LoadBalancer"
  }
}
`, name)
}

func testAccKubernetesServiceConfig_externalName(name string) string {
	return fmt.Sprintf(`resource "kubernetes_service" "test" {
  metadata {
    name = "%s"
  }

  spec {
    type          = "ExternalName"
    external_name = "terraform.io"
  }
}
`, name)
}

func testAccKubernetesServiceConfig_generatedName(prefix string) string {
	return fmt.Sprintf(`resource "kubernetes_service" "test" {
  metadata {
    generate_name = "%s"
  }

  spec {
    port {
      port        = 8080
      target_port = 80
    }
  }
}
`, prefix)
}
