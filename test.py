import os
import subprocess
import sys
import time


def main():
    working = os.path.abspath(os.getcwd())
    repetitions = 0
    state = 0
    while True:
        new_state = sum(_checksums(working))
        if state != new_state:
            repetitions += 1
            _display_repetitions_banner(repetitions)
            _run_tests(working)
            state = new_state
        time.sleep(.75)


def _checksums(working):
    for root, dirs, files in os.walk(working):
        for f in files:
            if f.endswith('.go'):
                stats = os.stat(os.path.join(root, f))
                yield stats.st_mtime + stats.st_size


def _display_repetitions_banner(repetitions):
    number = ' {} '.format(repetitions)
    half_delimiter = (EVEN if not repetitions % 2 else ODD) * \
                     ((80 - len(number)) / 2)
    print '\n{0}{1}{0}\n'.format(half_delimiter, number)


def _run_tests(working):
    os.chdir(working)
    _run_test()
    for root, dirs, files in os.walk(working):
        for d in dirs:
            dir_path = os.path.join(root, d)
            if '.git' in dir_path:
                continue

            os.chdir(os.path.join(root, d))
            _run_test()


def _run_test():
    os.system('go test -i')
    os.system('go test {0}'.format(sys.argv[-1] if sys.argv[-1] == '-v' else ''))
    print


EVEN = '='
ODD = '-'


if __name__ == '__main__':
    main()