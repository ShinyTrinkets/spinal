#! /usr/bin/env python3

"""
Direct usage:
python version_bump.py `cat VERSION` --patch > VERSION
python version_bump.py `cat VERSION` --minor > VERSION

With make:
make version BUMP=minor
make version BUMP=major
"""

import sys
import click

MIN_DIGITS = 2
MAX_DIGITS = 3

@click.command()
@click.argument('version')
@click.option('--major', 'bump_idx', flag_value=0, help='Increment major number.')
@click.option('--minor', 'bump_idx', flag_value=1, help='Increment minor number.')
@click.option('--patch', 'bump_idx', flag_value=2, default=True, help='Increment patch number.')
def cli(version, bump_idx):
    """
    Bumps a MAJOR.MINOR.PATCH version string at the specified index location or 'patch' digit.
    An optional 'v' prefix is allowed and will be included in the output if found.
    """
    prefix = version[0] if version[0].isalpha() else ''
    digits = version.lower().lstrip('v').split('.')

    if len(digits) > MAX_DIGITS:
        click.secho('ERROR: Too many digits', fg='red', err=True)
        sys.exit(1)

    digits = (digits + ['0'] * MAX_DIGITS)[:MAX_DIGITS]  # Extend total digits to max
    digits[bump_idx] = str(int(digits[bump_idx]) + 1)  # Increment the desired digit

    # Zero rightmost digits after bump position
    for i in range(bump_idx + 1, MAX_DIGITS):
        digits[i] = '0'

    digits = digits[:max(MIN_DIGITS, bump_idx + 1)]  # Trim rightmost digits
    click.echo(prefix + '.'.join(digits), nl=False)


if __name__ == '__main__':
    cli()
