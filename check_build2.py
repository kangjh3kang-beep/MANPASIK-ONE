#!/usr/bin/env python3

f = '/home/kangjh3kang/Manpasik/frontend/flutter-app/build/web/main.dart.js'
text = open(f, encoding='utf-8').read()

# Dart2JS encodes Korean as \uXXXX sequences
# 리 = \ub9ac, 더 = \ub354
# 뇌 = \ub1cc, 활 = \ud65c, 동 = \ub3d9
# 심 = \uc2ec, 박 = \ubc15, 수 = \uc218

import re

tests = {
    '리더': r'\ub9ac\ub354',
    '뇌활동': r'\ub1cc\ud65c\ub3d9',
    '심박수': r'\uc2ec\ubc15\uc218',
    '혈당': r'\ud608\ub2f9',
    '바디스캔': r'\ubc14\ub514\uc2a4\uce94',
    '산소포화도': r'\uc0b0\uc18c\ud3ec\ud654\ub3c4',
    '연결됨': r'\uc5f0\uacb0\ub428',
    '배터리': r'\ubc30\ud130\ub9ac',
    '차동측정': r'\ucc28\ub3d9\uce21\uc815',
    '활성': r'\ud65c\uc131',
}

for name, pattern in tests.items():
    idx = text.find(pattern)
    print(f"{'OK' if idx >= 0 else 'MISS'}: {name} (pos={idx})")

# Also check old English
eng_tests = {
    'Reader Device': 'Reader Device',
    'BRAIN ACTIV': 'BRAIN ACTIV',
    'HEART RATE': 'HEART RATE',
    'BODY SCAN': 'BODY SCAN',
    'GLUCOSE': 'GLUCOSE',
}
print("---")
for name, pattern in eng_tests.items():
    idx = text.find(pattern)
    if idx >= 0:
        print(f"ENG FOUND: {name} (pos={idx})")
    else:
        print(f"ENG GONE: {name}")
