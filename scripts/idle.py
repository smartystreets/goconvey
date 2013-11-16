#!/usr/bin/env python

"""
This script scans the current working directory for changes to .go files and 
runs `go test` in each folder where *_test.go files are found. It does this 
indefinitely or until a KeyboardInterrupt is raised (<Ctrl+c>). This script 
passes the verbosity command line argument (-v) to `go test`.
"""


import os
import subprocess
import sys
import time


def main(verbose):
    working = os.path.abspath(os.path.join(os.getcwd()))    
    scanner = WorkspaceScanner(working)
    runner = TestRunner(working, verbose)

    while True:
        if scanner.scan():
            runner.run()


class WorkspaceScanner(object):
    def __init__(self, top):
        self.state = 0
        self.top = top

    def scan(self):
        time.sleep(.75)
        new_state = sum(self._checksums())
        if self.state != new_state:
            self.state = new_state
            return True
        return False

    def _checksums(self):
        for root, dirs, files in os.walk(self.top):
            for f in files:
                if f.endswith('.go'):
                    try:
                        stats = os.stat(os.path.join(root, f))
                        yield stats.st_mtime + stats.st_size
                    except OSError:
                        pass


class TestRunner(object):
    def __init__(self, top, verbosity):
        self.repetitions = 0
        self.top = top
        self.working = self.top
        self.verbosity = verbosity

    def run(self):
        self.repetitions += 1
        self._display_repetitions_banner()
        self._run_tests()

    def _display_repetitions_banner(self):
        number = ' {} '.format(self.repetitions if self.repetitions % 50 else
            'Wow, are you going for a top score? Keep it up!')
        half_delimiter = (EVEN if not self.repetitions % 2 else ODD) * \
                         ((80 - len(number)) / 2)
        write('\n{0}{1}{0}\n'.format(half_delimiter, number))

    def _run_tests(self):
        self._chdir(self.top)
        if self.tests_found():
            self._run_test()
        
        for root, dirs, files in os.walk(self.top):
            self.search_for_tests(root, dirs, files)

    def search_for_tests(self, root, dirs, files):
        for d in dirs:
            if '.git' in d or '.git' in root:
                continue

            self._chdir(os.path.join(root, d))
            if self.tests_found():
                self._run_test()

    def tests_found(self):
        for f in os.listdir(self.working):
            if f.endswith('_test.go'):
                return True

        return False

    def _run_test(self):
        subprocess.call('go test -i', shell=True)
        try:
            output = subprocess.check_output(
                'go test ' + self.verbosity, shell=True)
            self.write_output(output)
        except subprocess.CalledProcessError as error:
            self.write_output(error.output)

        write('\n')

    def write_output(self, output):
        write(output)

    def _chdir(self, new):
        os.chdir(new)
        self.working = new


def write(value):
    sys.stdout.write(value)
    sys.stdout.flush()


EVEN = '='
ODD  = '-'
RESET_COLOR  = '\033[0m'
RED_COLOR    = '\033[31m'
YELLOW_COLOR = '\033[33m'
GREEN_COLOR  = '\033[32m'


def parse_bool_arg(name):
    for arg in sys.argv:
        if arg == name:
            return True
    return False


if __name__ == '__main__':
    verbose = '-v' if parse_bool_arg('-v') else ''
    main(verbose)
