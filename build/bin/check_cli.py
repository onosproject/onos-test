#!/usr/bin/env python3
"""
 SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>

 SPDX-License-Identifier: Apache-2.0
"""

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
        for index in range(int(0), int(len(expected_regexes))):
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
