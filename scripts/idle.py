#!/usr/bin/env python

"""
This script scans the current working directory for changes to .go files and runs 
`go test` in each folder where *_test.go files are found. It does this indefinitely 
or until a KeyboardInterrupt is raised (<Ctrl+c>). This script passes the verbosity 
command line argument (-v) to `go test`.
"""

import os
import subprocess
import sys
import time


def main(verbose, output):
    working = os.path.abspath(os.path.join(os.getcwd()))    
    output = OutputWriter(output)
    scanner = WorkspaceScanner(working)
    runner = TestRunner(working, output, verbose)

    while True:
        if scanner.scan():
            output.start()
            runner.run()
            output.finish()


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
                    stats = os.stat(os.path.join(root, f))
                    yield stats.st_mtime + stats.st_size


class TestRunner(object):
    def __init__(self, top, out, verbosity):
        self.repetitions = 0
        self.top = top
        self.out = out
        self.working = self.top
        self.verbosity = verbosity

    def run(self):
        self.repetitions += 1
        self._display_repetitions_banner()
        self._run_tests()

    def _display_repetitions_banner(self):
        number = ' {} '.format(self.repetitions)
        half_delimiter = (EVEN if not self.repetitions % 2 else ODD) * \
                         ((80 - len(number)) / 2)
        self.out.write('\n{0}{1}{0}\n'.format(half_delimiter, number))

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
            output = subprocess.check_output('go test ' + self.verbosity, shell=True)
            self.write_output(output)
        except subprocess.CalledProcessError as error:
            self.write_output(error.output)

        self.out.write('\n')

    def write_output(self, output):
        self.out.write(output)

    def _chdir(self, new):
        os.chdir(new)
        self.working = new


class OutputWriter(object):
    def __init__(self, output):
        self.output = output
        self.logfile = None
        self.working = None
        self.finished = None
        
    def start(self):
        if self.output:
            self.finished = os.path.join(output, 'latest.txt')
            self.working = os.path.join(output, 'running.txt')
            self.logfile = open(self.working, 'w')

    def write(self, value):
        sys.stdout.write(value)
        sys.stdout.flush()
        if self.logfile is not None:
            self.write_to_log(value)
        
    def write_to_log(self, value):
        output = value\
            .replace(RED_COLOR, '')\
            .replace(GREEN_COLOR, '')\
            .replace(RESET_COLOR, '')
        self.logfile.write(output)
        self.logfile.flush()

    def finish(self):
        if self.logfile is not None:
            self.logfile.close()
            os.rename(self.working, self.finished)


EVEN = '='
ODD = '-'
RESET_COLOR = '\033[0m'
RED_COLOR = '\033[31m'
GREEN_COLOR = '\033[32m'


def parse_bool_arg(name):
    for arg in sys.argv:
        if arg == name:
            return True
    return False


def parse_string_arg(name):
    for arg in sys.argv:
        if arg.startswith(name + '='):
            return arg.split('=')[-1]
    return None


if __name__ == '__main__':
    verbose = '-v' if parse_bool_arg('-v') else ''
    output = parse_string_arg('--output')
    main(verbose, output)
