#!/usr/bin/env python3
""""
 Copyright 2021-present Open Networking Foundation.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.

 CLI checker utility functions
"""

import sys
import re
import time


def do_check(name, cli_lines, expected_headers, expected_regexes):
    connections_found = 0
    headers_found = 0

    if len(cli_lines) == 0:
        print('No CLI output')
        exit(1)
    cli_headers = cli_lines[0].split()

    if cli_headers != expected_headers:
        print('CLI headers are incorrect: ' + cli_lines[0])
        return False
    headers_found += 1

    for connection in cli_lines[1:]:
        connection_fields = connection.split()
        for index in range(0, len(expected_regexes)-1):
            expected_field = expected_regexes[index]

            if len(connection_fields) != len(expected_regexes):
                print('Wrong number of fields in: ' + connection)
                return False

            if not re.fullmatch(expected_field, connection_fields[index]):
                print('Connection incorrect: ' + connection_fields[index] + ' does not match ' + expected_field)
                return False

        connections_found += 1

    if headers_found == 0:
        print('CLI Headers not found')
        return False

    if connections_found == 0:
        print('No connections found')
        return False

    print(name + 'CLI output for is correct!')
    return True


def check_cli_output(name, expected_headers, expected_regexes):
    cli_lines = str.splitlines(sys.stdin.read())

    return do_check(name, cli_lines, expected_headers, expected_regexes)
