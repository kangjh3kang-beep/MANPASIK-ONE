import re
import os

file_path = os.path.expanduser("~/Manpasik/frontend/flutter-app/lib/generated/manpasik.pb.dart")

if not os.path.exists(file_path):
    print(f"Error: File not found at {file_path}")
    exit(1)

with open(file_path, 'r', encoding='utf-8') as f:
    lines = f.readlines()

new_lines = []
in_class = False
class_name = ""
brace_depth = 0
has_clone = False
has_create = False
has_info = False

class_start_pattern = re.compile(r'class\s+(\w+)\s+extends\s+\$pb\.GeneratedMessage\s+\{')

for i, line in enumerate(lines):
    # Check for class start
    match = class_start_pattern.search(line)
    if match:
        in_class = True
        class_name = match.group(1)
        brace_depth = 0
        has_clone = False
        has_create = False
        has_info = False
        # Count braces in this line
        brace_depth += line.count('{') - line.count('}')
        new_lines.append(line)
        continue

    if in_class:
        # Check for existing methods
        if f'{class_name} clone()' in line:
            has_clone = True
        if f'{class_name} createEmptyInstance()' in line:
            has_create = True
        if 'BuilderInfo get info_' in line:
            has_info = True

        # Track brace depth
        open_braces = line.count('{')
        close_braces = line.count('}')
        brace_depth += open_braces - close_braces

        # Check for end of class
        if brace_depth == 0 and close_braces > 0:
            # We are at the closing brace of the class
            # Insert missing methods before the closing brace
            
            # Use specific indentation (usually 2 spaces)
            indent = "  "
            
            methods_to_add = []
            
            if not has_clone:
                methods_to_add.append(f"\n{indent}@$core.override")
                methods_to_add.append(f"{indent}{class_name} clone() => {class_name}()..mergeFromMessage(this);")
            
            if not has_create:
                methods_to_add.append(f"\n{indent}@$core.override")
                methods_to_add.append(f"{indent}{class_name} createEmptyInstance() => CreateMessage();")

            if not has_info:
                methods_to_add.append(f"\n{indent}@$core.override")
                methods_to_add.append(f"{indent}$pb.BuilderInfo get info_ => _i;")

            if methods_to_add:
                # Add methods before the current line (which contains the closing brace)
                # But wait, if the closing brace is on a line with other stuff, we might split it?
                # Usually generated code puts '}' on its own line or at end.
                # Let's assume '}' is the last significant char.
                
                # Check if line allows insertion before '}'
                # Simple strategy: Just append methods to new_lines before appending the current line
                for m_line in methods_to_add:
                    new_lines.append(m_line + "\n")
            
            in_class = False
            class_name = ""
            
        new_lines.append(line)
    else:
        new_lines.append(line)

with open(file_path, 'w', encoding='utf-8') as f:
    f.writelines(new_lines)

print(f"Successfully patched {file_path}")
