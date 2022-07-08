#! /usr/bin/env python

import sys, re

# This script reads a copy of lscpu-arm.c and outputs the various tables to be
# included in cpu_linux_arm.go.  Of course, this may need adjustment as other
# changes are made to lscpu-arm.c.

def main():
    lines = iter(open(sys.argv[1], "r"))

    # part_lists are the lists of part numbers, keyed by name
    part_lists = {}

    # hw_implementer is the hw_implementer array
    hw_implementer = None

    try:
        while True:
            line = next(lines)
            mo = re.match(r"static const struct ([a-z0-9_]*) ([a-z0-9_]*)\[\] = {", line)
            if mo:
                if mo.group(1) == "id_part":
                    part_lists[mo.group(2)] = parse_part_list(lines)
                if mo.group(1) == "hw_impl":
                    hw_implementer = parse_hw_impl_list(lines)
    except StopIteration:
        pass # all finished!

    print("var hwVariant = map[uint64]hwImpl{")
    for hw_impl, ident, name in hw_implementer:
        print(f"\t{hw_impl}: hwImpl{{")
        print(f"\t\tname: {name},")
        print(f"\t\tparts: map[uint64]string{{")
        for k, v in part_lists[ident]:
            print(f"\t\t\t{k}: {v},")
        print(f"\t\t}},")
        print(f"\t}},")
    print(f"}}")


def parse_part_list(lines):
    """Parse lines of the form `{ k, v },`, where k is a number and v is a string"""
    rv = []
    while True:
        line = next(lines)
        mo = re.match(r'\s*{\s*([0-9a-fx]*),\s*(".*")\s*},.*', line)
        if mo:
            k, v = mo.group(1, 2)
            if k == "-1":
                continue
            rv.append((k,v))
        else:
            return rv


def parse_hw_impl_list(lines):
    """Parse lines of the form `{ k, ident, name },`, where k is a number, ident is
    an identifier, and name is a string"""
    rv = []
    while True:
        line = next(lines)
        mo = re.match(r'\s*{\s*([0-9a-fx]*),\s*([a-z0-9_]*),\s*(".*")\s*},.*', line)
        if mo:
            k, ident, name = mo.group(1, 2, 3)
            if k == "-1":
                continue
            rv.append((k, ident, name))
        else:
            return rv


if __name__ == "__main__":
    main()
