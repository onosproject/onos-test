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

const envoyConfig = `
## refs:
      ## - https://www.envoyproxy.io/docs/envoy/latest/start/start#quick-start-to-run-simple-example
      ## - https://raw.githubusercontent.com/envoyproxy/envoy/master/configs/google_com_proxy.v2.yaml
      admin:
        access_log_path: /tmp/admin_access.log
        address:
          socket_address: { address: 0.0.0.0, port_value: 9901 }
      static_resources:
        listeners:
        {{- range $key, $value := .Values.onosservices }}
          - name: listener_{{$key}}
            address:
              socket_address: { address: 0.0.0.0, port_value: {{ $value.proxy }} }
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
                                  cluster: onos_{{ $key }}_service
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
        {{- end }}
        clusters:
        {{- range $key, $value := .Values.onosservices }}
          - name: onos_{{ $key }}_service
            connect_timeout: 0.25s
            type: logical_dns
            http2_protocol_options: {}
            lb_policy: round_robin
            # win/mac hosts: Use address: host.docker.internal instead of address: localhost in the line below
            hosts: [{ socket_address: { address: onos-{{$key}}, port_value: {{ $value.grpc }} }}]
            tls_context:
              common_tls_context:
                tls_certificates:
                  certificate_chain: { "filename": "/etc/envoy-proxy/certs/tls.crt" }
                  private_key: { "filename": "/etc/envoy-proxy/certs/tls.key" }
                validation_context:
                  trusted_ca: { filename: "/etc/envoy-proxy/certs/tls.cacrt" }
        {{- end }}`

var envoyConfigMaps = map[string]string{
	"/config/envoy.yaml": envoyConfig,
}

const (
	guiType    = "gui"
	guiImage   = "onosproject/onos-gui:latest"
	envoyImage = "envoyproxy/envoy:v1.11.1"
	guiService = "onos-gui"
	guiPort    = 80
)

// Enabled indicates whether the Gui is enabled
func (c *Gui) Enabled() bool {
	return GetArg(c.name, "enabled").Bool(c.enabled)
}

// SetEnabled sets whether the Gui is enabled
func (c *Gui) SetEnabled(enabled bool) {
	c.enabled = enabled
}

func newGui(cluster *Cluster) *Gui {
	service := newService(cluster)
	service.SetPorts([]Port{{Name: "grpc", Port: guiPort}})
	service.SetName(guiService)
	service.SetLabels(getLabels(guiType))
	service.SetConfigMaps(envoyConfigMaps)

	guiContainer := newContainer(cluster)
	var containers []*Container
	guiContainer.SetName(guiService)
	guiContainer.SetImage(guiImage)
	containers = append(containers, guiContainer)

	envoyContainer := newContainer(cluster)
	envoyContainer.SetName("onos-envoy")
	envoyContainer.SetImage(envoyImage)
	envoyContainer.SetCommand("/usr/local/bin/envoy")
	envoyContainer.SetArgs("-c", "/config/envoy.yaml")

	containers = append(containers, envoyContainer)

	service.SetContainers(containers)

	return &Gui{
		Service: service,
	}

}

// Gui provides methods for managing the onos-gui service
type Gui struct {
	*Service
	enabled bool
}
