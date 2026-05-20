#!/usr/bin/env python3
import sys

f = '/home/kangjh3kang/Manpasik/frontend/flutter-app/build/web/main.dart.js'
data = open(f, 'rb').read()

# Check raw UTF-8 bytes
tests = {
    '리더 디바이스': '리더 디바이스'.encode('utf-8'),
    '뇌활동': '뇌활동'.encode('utf-8'),
    '심박수': '심박수'.encode('utf-8'),
    '혈당': '혈당'.encode('utf-8'),
    '바디스캔': '바디스캔'.encode('utf-8'),
    '산소포화도': '산소포화도'.encode('utf-8'),
    '연결됨': '연결됨'.encode('utf-8'),
    '배터리': '배터리'.encode('utf-8'),
    '차동측정시스템': '차동측정시스템'.encode('utf-8'),
    '활성': '활성'.encode('utf-8'),
}

for name, b in tests.items():
    idx = data.find(b)
    print(f"{'OK' if idx >= 0 else 'MISS'}: {name} (pos={idx})")

print("---")
eng_tests = {
    'Reader Device Overview': b'Reader Device Overview',
    'BRAIN ACTIV': b'BRAIN ACTIV',
    'HEART RATE': b'HEART RATE',
    'BODY SCAN ACTIVE': b'BODY SCAN ACTIVE',
}
for name, b in eng_tests.items():
    idx = data.find(b)
    if idx >= 0:
        print(f"ENG FOUND: {name} (pos={idx})")
    else:
        print(f"ENG GONE: {name}")
