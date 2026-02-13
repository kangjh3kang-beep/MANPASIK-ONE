#!/usr/bin/env python3
"""Generate cd.yml - run then delete this script."""
import os, sys

# Read the template from the accompanying .txt file
here = os.path.dirname(os.path.abspath(__file__))
src = os.path.join(here, '_cd_template.txt')
dst = os.path.join(here, 'cd.yml')

with open(src, 'r', encoding='utf-8') as f:
    content = f.read()

with open(dst, 'w', encoding='utf-8') as f:
    f.write(content)

print(f'Written {len(content)} bytes to cd.yml')
