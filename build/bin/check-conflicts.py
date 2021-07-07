#!/usr/bin/env python3

# Copyright 2021-present Open Networking Foundation.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# CLI checker utility functions

# using exact matches b/c would otherwise need different of input (pci's need to be exact, hard to generalize)
# values based off of model.yaml in sdran-helm-charts/ran-simulator/files/model that is automatically loaded

import sys
import time

expected_output = ["NCGI                 #UEs Max UEs    TxDB       Lat       Lng Azimuth     Arc   A3Offset     TTT  A3Hyst CellOffset FreqOffset      PCI    Color Idle Conn Neighbors",
                   "13842601454c002         0   99999   40.00    52.486    13.412     120     120          0       0       0          0          0      218    green    0,    0, 13842601454c001,13842601454c003",
                   "13842601454c001         0   99999   40.00    52.486    13.412       0     120          0       0       0          0          0       80    green    2,    0, 13842601454c003,13842601454c002",
                   "138426014550002         0   99999   40.00    52.504    13.453     120     120          0       0       0          0          0      115    green    4,    4, 138426014550001,138426014550003",
                   "13842601454c003         0   99999   40.00    52.486    13.412     240     120          0       0       0          0          0      148    green    1,    0, 13842601454c001,13842601454c002",
                   "138426014550001         0   99999   40.00    52.504    13.453       0     120          0       0       0          0          0       69    green    0,    2, 138426014550003,138426014550002",
                   "138426014550003         0   99999   40.00    52.504    13.453     240     120          0       0       0          0          0      218    green    1,    0, 138426014550002,138426014550001",
                   "Cell is now ecgi:87893173159133185 location:<lat:52.50431527434924 lng:13.453261970488306 > sector:<arc:120 centroid:<lat:52.50431527434924 lng:13.453261970488306 > height:46 tilt:-15 > color:\"green\" max_ues:99999 neighbors:87893173159133187 neighbors:87893173159133186 tx_power_db:40 measurement_params:<event_a3_params:<> > pci:115 rrc_connected_count:2 Cell 138426014550001 updated",
                   "ID                Node ID   Dlearfcn   Cell Type   PCI   PCI Pool",
                   "13842601454c003   5153      0          FEMTO       148   [1:512]",
                   "138426014550001   5154      0          FEMTO       145   [1:512]",
                   "138426014550002   5154      0          FEMTO       115   [1:512]",
                   "138426014550003   5154      0          FEMTO       218   [1:512]",
                   "13842601454c001   5153      0          FEMTO       80    [1:512]",
                   "13842601454c002   5153      0          FEMTO       218   [1:512]",
                   "ID   Node ID   Dlearfcn   Cell Type   PCI   PCI Pool",
                   "ID                Total Resolved Conflicts   Most Recent Resolution",
                   "13842601454c001   1                          148=>80",
                   "138426014550001   2                          115=>145"]

if __name__ == '__main__':
    check = 'PCI app get resolve conflicts'
    input_lines = sys.stdin.read().splitlines()
    ok = True
    if len(expected_output) != len(input_lines):
        ok = False
    else: 
        # check each line against expected output
        for i in range(len(expected_output)):
            if expected_output[i] != input_lines[i]:
                if i == 0:
                    print("Failed at ransim get cells header")
                elif i < 7:
                    print("Failed at ransim get cells body")
                elif i < 8:
                    print("Failed at ransim set cell 38426014550001 --pci 115")
                elif i < 9:
                    print("Failed at pci get cells header")
                elif i < 10:
                    print("Failed at pci get cells body")
                elif i < 16:
                    print("Failed at pci get conflicts header")
                elif i < 17:
                    print("Failed at pci get resolved header")
                elif i < 19:
                    print("Failed at pci get resolved body")
                ok = False
                break

    if not ok:
        print("Check " + check + " failed")
        exit(1)

    print("Check " + check + " passed")
