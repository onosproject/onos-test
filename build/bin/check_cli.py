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


import sys
import re


def print_with_identifier(name, message):
    identifier = '**** (' + name + ') '
    print(identifier + message)


def do_check(name, cli_lines, expected_headers, expected_regexes):
    items_found = int(0)
    headers_found = int(0)

    if int(len(cli_lines)) == 0:
        print_with_identifier(name, 'No CLI output')
        exit(int(1))
    cli_headers = cli_lines[0].split()

    if cli_headers != expected_headers:
        print_with_identifier(name, 'CLI headers are incorrect: ' + cli_lines[0])
        return False
    headers_found += 1

    for item in cli_lines[1:]:
        fields = item.split()
        for index in range(int(0), int(len(expected_regexes)-1)):
            expected_field = expected_regexes[index]

            if int(len(fields)) != len(expected_regexes):
                print_with_identifier(name, 'Wrong number of fields in: ' + item)
                print(fields)
                return False

            if not re.fullmatch(expected_field, fields[index]):
                print_with_identifier(name, 'Item incorrect: ' + fields[index] + ' does not match ' + expected_field)
                print(fields)
                return False

        items_found += 1

    if headers_found == 0:
        print_with_identifier(name, 'CLI Headers not found')
        return False

    if items_found == 0:
        print_with_identifier(name, 'No items found')
        return False

    print_with_identifier(name, 'CLI output is correct!')
    return True


def check_cli_output(name, expected_headers, expected_regexes):
    cli_lines = str.splitlines(sys.stdin.read())

    return do_check(name, cli_lines, expected_headers, expected_regexes), cli_lines
