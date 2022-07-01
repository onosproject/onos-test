# SPDX-FileCopyrightText: 2022-present Intel Corporation
# SPDX-FileCopyrightText: 2022-present Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0


*** Settings ***
Documentation             Smoke test for onos-config
Suite Setup               Setup Suite
Suite Teardown            Teardown Suite
Resource                  ${CURDIR}/../libraries/onos.robot

*** Variables ***
${NAMESPACE}              onos-test
${REGISTRY_SETTINGS}      ${EMPTY}
${ONF_PROXY}              mirror.registry.opennetworking.org
${ONOS_CHART_REPO}        ${EMPTY}
${SDRAN_CHART_REPO}       ${EMPTY}
${USE_LATEST}             true
${TAG_OPTIONS}            --set global.image.tag\=latest

*** Keywords ***
Setup Suite
    [Documentation]       Start the onos-config cluster
    Setup Cluster
    Install onos-config helm chart
    Install device simulator helm chart

Create topo entries
    [Documentation]     Create topo entries for the simulated device
    ${cli_output}=    ONOS CLI    onos topo create kind devicesim devicesim
    Should Be Equal As Integers    ${cli_output.rc}    0

    ${cli_command_create}=    CATENATE     onos topo create entity devicesim-1
                                           ...    -a onos.topo.TLSOptions='{"insecure":true,"plain":true}'
                                           ...    -a onos.topo.Configurable='{"address":"devicesim1-device-simulator:11161","version":"1.0.0","type":"devicesim"}'
                                           ...    -k devicesim
    ${cli_output}=    ONOS CLI    ${cli_command_create}
    Should Be Equal As Integers    ${cli_output.rc}    0

*** Test Cases ***

Create topo entities
    [Tags]    create_topo
    Create topo entries

Check topo objects
    [Tags]      check_topo
    ${topo_objects}=   ONOS CLI No Headers    onos topo get objects
    @{lines}=   Split To Lines    ${topo_objects.stdout}
    FOR     ${topo_entry}   IN  @{lines}
        @{words}=    Split String    ${topo_entry}

        IF    "${words}[0]" == "KIND"
            Should Match Regexp    ${words[1]}    ^devicesim$
            Should Match Regexp    ${words[2]}    ^devicesim$
            Should Match Regexp    ${words[3]}    ^<None>$
            Should Match Regexp    ${words[4]}    ^<None>$
            Should Match Regexp    ${words[5]}    ^<None>$
            Should Match Regexp    ${words[6]}    ^<None>$

        ELSE IF    "${words}[0]" == "ENTITY"
            Should Match Regexp    ${words[1]}    ^(gnmi:onos-config.*)|(devicesim.*)$
            Should Match Regexp    ${words[2]}    ^(onos-config|devicesim.*)$
            Should Match Regexp    ${words[3]}    ^<None>$
            Should Match Regexp    ${words[4]}    ^<None>$
            Should Match Regexp    ${words[5]}    ^<None>$
            Should Match Regexp    ${words[6]}    ^onos\.topo\..*$
        ELSE
            FAIL    Unknown Object Type ${words[0]}
        END
    END

Check plugins
    [Tags]      check_plugins
    ${plugin_output}=   ONOS CLI No Headers      onos config get plugins
    ${plugin_count}=    Run Process              echo "${plugin_output.stdout}" | grep -c Loaded    shell=yes
    Should Be Equal    ${plugin_count.stdout}    3

Check pods
    [Tags]      check_pods
    Check Pods

Check GNMI All Targets
    [Tags]      check_gnmi_all_targets
    ${gnmi_all_target}=    GNMI CLI Get    -proto "path: <target: '*'>"

    Should Contain    ${gnmi_all_target.stdout}      name: "all-targets"
    Should Contain    ${gnmi_all_target.stdout}      target: "*"
    Should Contain    ${gnmi_all_target.stdout}      targets\\": [\\"devicesim-1\\"]

Check GNMI Capabilities
    [Tags]      check_gnmi_capabilities
    ${gnmi_caps}=    GNMI CLI Capabilities
    Should Contain X Times    ${gnmi_caps.stdout}    supported_models    8
    Should Contain    ${gnmi_caps.stdout}            supported_encodings: JSON
    Should Contain    ${gnmi_caps.stdout}            supported_encodings: JSON_IETF
    Should Contain    ${gnmi_caps.stdout}            supported_encodings: PROTO
    Should Contain    ${gnmi_caps.stdout}            gNMI_version: "0.7.0"