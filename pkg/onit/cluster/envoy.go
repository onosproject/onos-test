// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cluster

const (
	envoyType    = "envoy"
	envoyService = "onos-envoy"
	envoyImage   = "envoyproxy/envoy-alpine:latest"
	envoyPort    = 8080
)

var envoyCommand = []string{
	"/usr/local/bin/envoy",
	"-c",
	"/etc/envoy-proxy/config/envoy-config.yaml",
}

const envoyConfig = `
admin:
  access_log_path: /tmp/admin_access.log
  address:
    socket_address: { address: 0.0.0.0, port_value: 9901 }
static_resources:
  listeners:
    - name: listener_0
      address:
        socket_address: { address: 0.0.0.0, port_value: 8080 }
      filter_chains:
        - filters:
            - name: envoy.http_connection_manager
              config:
                codec_type: auto
                stat_prefix: ingress_http
                route_config:
                  name: local_route
                  virtual_hosts:
                    - name: local_service
                      domains: ["*"]
                      routes:
                        - match: { prefix: "/" }
                          route:
                            cluster: onos_config_service
                            max_grpc_timeout: 0s
                      cors:
                        allow_origin:
                          - "*"
                        allow_methods: GET, PUT, DELETE, POST, OPTIONS
                        allow_headers: keep-alive,user-agent,cache-control,content-type,content-transfer-encoding,custom-header-1,x-accept-content-transfer-encoding,x-accept-response-streaming,x-user-agent,x-grpc-web,grpc-timeout
                        max_age: "1728000"
                        expose_headers: custom-header-1,grpc-status,grpc-message
                http_filters:
                  - name: envoy.grpc_web
                  - name: envoy.cors
                  - name: envoy.router
  clusters:
    - name: onos_config_service
      connect_timeout: 0.25s
      type: logical_dns
      http2_protocol_options: {}
      lb_policy: round_robin
      # win/mac hosts: Use address: host.docker.internal instead of address: localhost in the line below
      hosts: [{ socket_address: { address: onos-topo, port_value: 5150 }}]
      tls_context:
        common_tls_context:
          tls_certificates:
            certificate_chain: { "filename": "/certs/client.crt" }
            private_key: { "filename": "/certs/client.key" }
          validation_context:
            trusted_ca: { filename: "/certs/onf.cacrt" }`

var envoySecrets = map[string]string{
	"/certs/onf.cacrt":  caCert,
	"/certs/client.crt": clientCert,
	"/certs/client.key": clientKey,
}

// Enabled indicates whether the Gui is enabled
func (c *Envoy) Enabled() bool {
	return GetArg(c.name, "enabled").Bool(c.enabled)
}

// SetEnabled sets whether the Envoy is enabled
func (c *Envoy) SetEnabled(enabled bool) {
	c.enabled = enabled
}

func newEnvoy(cluster *Cluster) *Envoy {
	var envoyConfigMaps = map[string]string{
		"/etc/envoy-proxy/config/envoy-config.yaml": envoyConfig,
	}
	return &Envoy{
		Service: newService(cluster, envoyService, []Port{{Name: "envoy", Port: envoyPort}}, getLabels(envoyType), envoyImage, envoySecrets, nil, envoyCommand, envoyConfigMaps, nil),
	}
}

// Envoy provides methods for managing the envoy service
type Envoy struct {
	*Service
	enabled bool
}
