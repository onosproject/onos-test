# SPDX-FileCopyrightText: 2022-present Intel Corporation
# SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0


*** Settings ***
Documentation     Shared keywords for micro onos testing
Library           String
Library           Process

*** Variables ***
${NAMESPACE}              onos-test
${REGISTRY_SETTINGS}      ${EMPTY}
${ONF_PROXY}              mirror.registry.opennetworking.org
${ONOS_CHART_REPO}        ${EMPTY}
${SDRAN_CHART_REPO}       ${EMPTY}
${USE_LATEST}             true
${TAG_OPTIONS}            --set global.image.tag\=latest
${COMMON_GNMI_ARGS}       -address onos-config:5150  "-timeout" "5s" "-en" "JSON"
                          ...     "-alsologtostderr" "-insecure" "-client_crt" "/etc/ssl/certs/client1.crt"
                          ...     "-client_key" "/etc/ssl/certs/client1.key" "-ca_crt" "/etc/ssl/certs/onfca.crt"

*** Keywords ***
Setup Cluster
    [Documentation]
    ...  Create a kind cluster and install atomix and onos operator on it
    ...  The following environment variables can be used to alter the default behavior:
    ...  NAMESPACE - k8s namespace to use for the test.
    ...     Defaults to "micro-onos", override to use a different k8s namespace
    ...  START_KIND_CLUSTER - determines if the KinD cluster should be created.
    ...     Defaults to "true", set to false to use an existing cluster.
    ...  SDRAN_CHART_REPO - repo name for helm for the sd-ran charts.
    ...     Defaults to "sdran", override to use local charts
    ...  ONOS_CHART_REPO - repo name for helm for the ONOS charts.
    ...     Defaults to "onos", override to use local charts
    ...  USE_LATEST - specifies if the tag "latest" should be used for images.
    ...     Defaults to "true", set to "false" to run the tags defined in the charts

    # Configure variables
    Set Suite Variable    ${NAMESPACE}     %{NAMESPACE=onos-test}

    # set up artifacts directory
    ${result}=    Run Process    ${CURDIR}/../../build/bin/setup-artifacts      shell=yes
    Should Be Equal As Integers    ${result.rc}    0


    # Make sure the helm repos are up to date
    ${repos_update}=    Catenate    helm repo remove cord atomix onos sdran aether ;
                                    ...    helm repo add cord https://charts.opencord.org &&
                                    ...    helm repo add atomix https://charts.atomix.io &&
                                    ...    helm repo add onos https://charts.onosproject.org &&
                                    ...    helm repo add sdran https://sdrancharts.onosproject.org &&
                                    ...    helm repo add aether https://charts.aetherproject.org &&
                                    ...    helm repo update
    ${result}=          Run Process    ${repos_update}  shell=yes
    Should Be Equal As Integers    ${result.rc}    0

    # Create a kind cluster
    Set Local Variable    ${start_kind_cluster}    %{START_KIND_CLUSTER=true}
    IF    "${start_kind_cluster}" == "true"
        Run Process    kind delete cluster -q && kind create cluster -q    shell=yes
    END
    ${result}=   Run Process    kubectl delete ns ${NAMESPACE}; kubectl create ns ${NAMESPACE}   shell=yes
    Should Be Equal As Integers    ${result.rc}    0

    # Set up the proxy variables
    Set Local Variable    ${use_proxy}    %{USE_PROXY=true}
    IF    "${use_proxy}" == "true"
        ${proxy}=    Catenate    mirror.registry.opennetworking.org
        Set Suite Variable     ${REGISTRY_SETTINGS}    --set global.image.registry\=${proxy} --set global.store.consensus.image.registry\=${proxy} --set global.storage.consensus.image\=${proxy}/atomix/raft-storage-node:v0.5.3
    ELSE
        Set Suite Variable    ${REGISTRY_SETTINGS}     ${EMPTY}
    END

    # Install Atomix charts
    IF    "${start_kind_cluster}" == "true"
        Run Process    helm install --set image.registry\=${ONF_PROXY} --set init.image.registry\=${ONF_PROXY} --set broker.image.registry\=${ONF_PROXY} atomix-controller atomix/atomix-controller -n kube-system --wait    shell=yes
        Run Process    helm install --set image.registry\=${ONF_PROXY} --set driver.image.registry\=${ONF_PROXY} --set node.image.registry\=${ONF_PROXY} atomix-raft-storage atomix/atomix-raft-storage -n kube-system --wait    shell=yes

        # Install onos-operator chart
        Run Process    helm install --set global.image.registry\=${ONF_PROXY} -n kube-system  onos-operator onos/onos-operator --wait    shell=yes
    END

    # Set up repo overrides
    Set Suite Variable    ${ONOS_CHART_REPO}    %{ONOS_CHART_REPO=onos}
    IF    "${ONOS_CHART_REPO}" == ""
        Set Suite Variable    ${ONOS_CHART_REPO}    onos
    END
    Set Suite Variable    ${SDRAN_CHART_REPO}    %{SDRAN_CHART_REPO=sdran}
    IF    "${SDRAN_CHART_REPO}" == ""
        Set Suite Variable    ${SDRAN_CHART_REPO}    sdran
    END

    # Set up tags
    Set Suite Variable    ${USE_LATEST}    %{USE_LATEST=true}
    IF    "${USE_LATEST}" == "true"
        Set Suite Variable    ${USE_LATEST}    true
        Set Suite Variable    ${TAG_OPTIONS}   --set global.image.tag\=latest
    ELSE
        Set Suite Variable    ${USE_LATEST}    false
        Set Suite Variable    ${TAG_OPTIONS}   ${EMPTY}
    END

Teardown Suite
    ${result}=     Run Process    ${CURDIR}/../../build/bin/archive-artifacts    shell=yes
    Should Be Equal As Integers    ${result.rc}    0


ONOS CLI No Headers
    [Arguments]             ${command}
    ${cli_command}=         Catenate    kubectl -n ${NAMESPACE} exec -t deploy/onos-cli -- ${command} --no-headers
    ${cli_output}=          Run Process    ${cli_command}     shell=yes
    [Return]                ${cli_output}

ONOS CLI
    [Arguments]             ${command}
    ${cli_command}=         Catenate    kubectl -n ${NAMESPACE} exec -t deploy/onos-cli -- ${command}
    ${cli_output}=          Run Process    ${cli_command}     shell=yes
    [Return]                ${cli_output}

GNMI CLI Get
    [Arguments]             ${proto}
    ${gnmi_get_cmd}=        Catenate     gnmi_cli -get ${COMMON_GNMI_ARGS} ${proto}
    ${gnmi_get_result}=     ONOS CLI    ${gnmi_get_cmd}
    [Return]                ${gnmi_get_result}

GNMI CLI Capabilities
    ${gnmi_cap_cmd}=        Catenate     gnmi_cli -capabilities  ${COMMON_GNMI_ARGS}
    ${gnmi_cap_result}=     ONOS CLI    ${gnmi_cap_cmd}
    [Return]                ${gnmi_cap_result}

Check Pods
    [Documentation]     Checks that the k8s pods in the test namespace are all ready with no errors
    ${kube_output}=     Run Process     kubectl get pods -n ${NAMESPACE}  shell=yes
    Should Be Equal As Integers    ${kube_output.rc}    0
    ${kube_lines}=      Split To Lines    ${kube_output.stdout}
    FOR     ${kube_entry}   IN  @{kube_lines}
        @{kube_words}=    Split String    ${kube_entry}
        
        IF  "${kube_words}[0]" != "NAME"
            Should Match Regexp    ${kube_words}[1]     (\\d+)\/\\1
            Should Match   ${kube_words}[2]    Running
            Should Be Equal As Integers   ${kube_words}[3]    0
        END
    END

Install onos-config helm chart
    ${onos_helm_cmd}=     Catenate    helm install -n ${NAMESPACE} onos ${ONOS_CHART_REPO}/onos-umbrella ${TAG_OPTIONS} ${REGISTRY_SETTINGS} --wait
    ${result}=  Run Process      ${onos_helm_cmd}   shell=yes
    Should Be Equal As Integers    ${result.rc}    0

Install device simulator helm chart
    ${sim_helm_cmd}=    Catenate    helm install -n ${NAMESPACE} device-1 ${ONOS_CHART_REPO}/device-simulator ${TAG_OPTIONS} ${REGISTRY_SETTINGS} --wait
    ${result}=  Run Process       ${sim_helm_cmd}   shell=yes
    Should Be Equal As Integers    ${result.rc}    0


