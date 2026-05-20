import struct, json, sys, math

def analyze(path):
    with open(path, 'rb') as f:
        data = f.read()
    
    magic = struct.unpack_from('<I', data, 0)[0]
    version = struct.unpack_from('<I', data, 4)[0]
    json_len = struct.unpack_from('<I', data, 12)[0]
    gltf = json.loads(data[20:20+json_len])
    bin_start = 20 + json_len + 8  # skip json+bin chunk headers

    accessors   = gltf.get('accessors', [])
    buffer_views = gltf.get('bufferViews', [])
    meshes      = gltf.get('meshes', [])

    print(f"File: {path}")
    print(f"Magic: {hex(magic)}, Version: {version}")
    print(f"Meshes: {len(meshes)}, Accessors: {len(accessors)}")
    print()

    def read_floats(acc_idx):
        acc = accessors[acc_idx]
        bv  = buffer_views[acc['bufferView']]
        off = bv.get('byteOffset', 0) + acc.get('byteOffset', 0)
        cnt = acc['count']
        ch  = {'VEC3':3,'VEC2':2,'SCALAR':1}.get(acc['type'], 1)
        return [struct.unpack_from('<f', data, bin_start+off+i*4)[0]
                for i in range(cnt*ch)]

    def read_indices(acc_idx):
        acc = accessors[acc_idx]
        bv  = buffer_views[acc['bufferView']]
        off = bv.get('byteOffset', 0) + acc.get('byteOffset', 0)
        cnt = acc['count']
        ct  = acc['componentType']
        fmt = '<I' if ct==5125 else '<H' if ct==5123 else '<B'
        size = 4 if ct==5125 else 2 if ct==5123 else 1
        return [struct.unpack_from(fmt, data, bin_start+off+i*size)[0]
                for i in range(cnt)]

    for i, mesh in enumerate(meshes):
        for j, prim in enumerate(mesh.get('primitives', [])):
            attrs = prim.get('attributes', {})
            pos_idx = attrs.get('POSITION')
            if pos_idx is None:
                continue
            acc = accessors[pos_idx]
            cnt = acc['count']
            print(f"Mesh[{i}] Prim[{j}]:")
            print(f"  vertices={cnt}, type={acc['type']}, componentType={acc['componentType']}")
            print(f"  hasNormals={'NORMAL' in attrs}, hasIndices={'indices' in prim}")
            
            pos = read_floats(pos_idx)
            xs = [pos[k*3]   for k in range(cnt)]
            ys = [pos[k*3+1] for k in range(cnt)]
            zs = [pos[k*3+2] for k in range(cnt)]
            print(f"  X: [{min(xs):.4f}, {max(xs):.4f}]  width={max(xs)-min(xs):.4f}")
            print(f"  Y: [{min(ys):.4f}, {max(ys):.4f}]  height={max(ys)-min(ys):.4f}")
            print(f"  Z: [{min(zs):.4f}, {max(zs):.4f}]  depth={max(zs)-min(zs):.4f}")
            
            if 'NORMAL' in attrs:
                nrm = read_floats(attrs['NORMAL'])
                nzs = [nrm[k*3+2] for k in range(len(nrm)//3)]
                print(f"  Normal.Z range: [{min(nzs):.4f}, {max(nzs):.4f}]")
            
            if 'indices' in prim:
                idx = read_indices(prim['indices'])
                print(f"  indices={len(idx)} ({len(idx)//3} triangles), componentType={accessors[prim['indices']]['componentType']}")
            print()

if __name__ == '__main__':
    path = sys.argv[1] if len(sys.argv) > 1 else 'web/models/Holo_M1.glb'
    analyze(path)
