module github.com/onosproject/onos-test

go 1.13

require (
	github.com/atomix/go-client v0.0.0-20200218200323-6fd69e684d05
	github.com/atomix/kubernetes-controller v0.0.0-20200202101151-b31765af9a0f
	github.com/dustinkirkland/golang-petname v0.0.0-20190613200456-11339a705ed2
	github.com/fatih/color v1.7.0
	github.com/gogo/protobuf v1.3.1
	github.com/iancoleman/strcase v0.0.0-20191112232945-16388991a334
	github.com/joncalhoun/pipe v0.0.0-20170510025636-72505674a733
	github.com/onosproject/onos-config v0.0.0-20200204191831-5c2803ee469d
	github.com/onosproject/onos-ric v0.0.0-20200225182040-dcf370614b8e
	github.com/onosproject/onos-topo v0.0.0-20200218171206-55029b503689
	github.com/openconfig/gnmi v0.0.0-20190823184014-89b2bf29312c
	github.com/spf13/cobra v0.0.6
	github.com/stretchr/testify v1.5.1
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a // indirect
	golang.org/x/tools v0.0.0-20200313205530-4303120df7d8 // indirect
	google.golang.org/grpc v1.27.1
	gopkg.in/yaml.v2 v2.2.8
	helm.sh/helm/v3 v3.1.1
	k8s.io/api v0.17.3
	k8s.io/apiextensions-apiserver v0.17.2
	k8s.io/apimachinery v0.17.3
	k8s.io/cli-runtime v0.17.3
	k8s.io/client-go v0.17.3
	sigs.k8s.io/controller-runtime v0.1.12
)

replace github.com/docker/docker => github.com/docker/engine v1.4.2-0.20200229013735-71373c6105e3
replace golang.org/x/tools => github.com/golangci/tools v0.0.0-20190915081525-6aa350649b1c