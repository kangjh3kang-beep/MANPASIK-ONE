import re

with open('/tmp/gl_raw.js', 'r') as f:
    src = f.read()

# Remove ALL import statements (handles multi-line too)
src = re.sub(
    r'import\s+[\s\S]*?from\s+[\'"][^\'"]+[\'"]\s*;',
    '',
    src,
    flags=re.MULTILINE
)
src = re.sub(r"import\s+['\"][^'\"]+['\"]\s*;", '', src, flags=re.MULTILINE)

# Remove export {...} blocks
src = re.sub(r'export\s*\{[^}]*\}\s*;?', '', src, flags=re.MULTILINE)
# Remove 'export default X' -> just keep X
src = re.sub(r'^export\s+default\s+', '', src, flags=re.MULTILINE)
# Remove bare 'export' keyword before class/function/const declarations
src = re.sub(r'^export\s+', '', src, flags=re.MULTILINE)

# Verify zero imports remain
remaining = [l for l in src.split('\n') if re.search(r'^\s*import\s', l)]
print('Remaining import lines:', len(remaining))

# Instead of manual alias list, use new Function approach:
# At runtime, we auto-alias ALL THREE properties into local scope,
# then run GLTFLoader code, then capture the GLTFLoader class.

# Escape backticks in src for template literal
src_escaped = src  # will be embedded as-is inside a JS string via new Function

# Build the wrapper:
# We use new Function('THREE', code)(THREE) so all vars are local.
# Inside we build aliases for every key of THREE, then run GLTFLoader src.
# GLTFLoader class will be defined as a local var we return.

wrapper = r"""(function(THREE){
  // Auto-alias ALL THREE exports into local scope at runtime
  var _aliasLines = Object.keys(THREE).map(function(k){
    return 'var ' + k + ' = _T["' + k + '"];';
  }).join('\n');

  // GLTFLoader source (import/export stripped)
  var _gltfSrc = _aliasLines + '\n' + __GLTF_SRC__ + '\n;return typeof GLTFLoader !== "undefined" ? GLTFLoader : null;';

  try {
    var _GLTFLoader = new Function('_T', _gltfSrc)(THREE);
    if (_GLTFLoader) {
      THREE.GLTFLoader = _GLTFLoader;
      console.log('[GLTFLoader UMD] Loaded OK:', typeof _GLTFLoader);
    } else {
      console.error('[GLTFLoader UMD] GLTFLoader class not found after evaluation');
    }
  } catch(e) {
    console.error('[GLTFLoader UMD] Error:', e.message, e.stack && e.stack.split('\\n')[1]);
  }
})(window.THREE || (window.THREE = {}));
"""

# JSON-encode the source to safely embed it as a JS string literal
import json
src_json = json.dumps(src)

wrapper = wrapper.replace('__GLTF_SRC__', src_json)

out_path = '/home/kangjh3kang/Manpasik/frontend/flutter-app/web/js/GLTFLoader.umd.js'
with open(out_path, 'w') as f:
    f.write(wrapper)
print('Written: ' + str(len(wrapper)) + ' chars -> ' + out_path)
print('First 200 chars:')
print(wrapper[:200])
