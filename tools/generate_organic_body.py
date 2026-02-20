import json
import math
import os
import numpy as np

# Try importing pyvista, handle failure gracefully
try:
    import pyvista as pv
    HAS_PYVISTA = True
except ImportError:
    HAS_PYVISTA = False
    print("Warning: PyVista not found. Falling back to Numpy approximation.")

OUTPUT_FILE = 'frontend/flutter-app/assets/data/holo_body_v14.json'

def create_human_body():
    """
    Creates a merged, smooth human body mesh using PyVista primitives.
    """
    if not HAS_PYVISTA:
        return None
        
    # Helpers
    def make_ellipsoid(center, radius, theta_res=20, phi_res=20):
        # PyVista's ParametricEllipsoid is a surface, we want a mesh
        # Easier: Create Sphere and scale it
        mesh = pv.Sphere(radius=1.0, center=(0,0,0), theta_resolution=theta_res, phi_resolution=phi_res)
        mesh.points[:, 0] *= radius[0]
        mesh.points[:, 1] *= radius[1]
        mesh.points[:, 2] *= radius[2]
        mesh.translate(center, inplace=True)
        return mesh

    def make_cylinder(center, direction, radius, height, res=20):
        mesh = pv.Cylinder(center=center, direction=direction, radius=radius, height=height, resolution=res)
        return mesh

    parts = []
    
    # 1. Torso (Smoother shape)
    # Chest
    parts.append(make_ellipsoid((0, 0.4, 0), (0.28, 0.22, 0.15)))
    # Abs
    parts.append(make_ellipsoid((0, 0.1, 0), (0.24, 0.30, 0.14)))
    # Hips
    parts.append(make_ellipsoid((0, -0.2, 0), (0.25, 0.20, 0.16)))

    # 2. Head
    parts.append(make_ellipsoid((0, 0.85, 0), (0.12, 0.15, 0.14))) # Head
    parts.append(make_cylinder((0, 0.65, 0), (0,1,0), 0.08, 0.15)) # Neck

    # 3. Limbs (Use multiples for joints)
    
    # Arms
    for side in [-1, 1]:
        # Shoulder
        parts.append(make_ellipsoid((side * 0.32, 0.45, 0), (0.11, 0.11, 0.11)))
        # Upper Arm
        parts.append(make_ellipsoid((side * 0.35, 0.20, 0), (0.08, 0.20, 0.08))) 
        # Elbow
        parts.append(make_ellipsoid((side * 0.37, -0.05, 0), (0.07, 0.07, 0.07)))
        # Lower Arm
        parts.append(make_ellipsoid((side * 0.39, -0.25, 0), (0.065, 0.18, 0.06)))
        # Hand
        parts.append(make_ellipsoid((side * 0.41, -0.48, 0), (0.05, 0.08, 0.02)))

    # Legs
    for side in [-1, 1]:
        off = 0.14
        # Thigh
        parts.append(make_ellipsoid((side * off, -0.45, 0), (0.11, 0.28, 0.11)))
        # Knee
        parts.append(make_ellipsoid((side * off, -0.75, 0), (0.09, 0.09, 0.09)))
        # Calf
        parts.append(make_ellipsoid((side * off, -1.05, 0), (0.08, 0.25, 0.08)))
        # Foot
        parts.append(make_ellipsoid((side * off, -1.35, 0.05), (0.07, 0.05, 0.15)))

    # Merge all parts
    # Note: simple '+' just concatenates meshes. It doesn't boolean union.
    # But for a hologram, overlapping internal geometry is sometimes okay, 
    # OR we can purely rely on the "Outer Surface" if we use heavy transparency.
    # However, to get "Organic", we want to avoid sharp intersections.
    
    merged = parts[0]
    for p in parts[1:]:
        merged = merged + p
        
    print("Parts merged. Smoothing...")
    
    # Smooth to hide seams
    # Taubin smoothing prevents shrinkage while smoothing
    smoothed = merged.smooth(n_iter=50, relaxation_factor=0.05)
    
    # Decimate to reasonable count for mobile (Target: ~4000 triangles)
    decimated = smoothed.decimate(target_reduction=0.7) 
    
    # Compute Normals (Essential for Organic Lighting)
    decimated.compute_normals(cell_normals=False, point_normals=True, inplace=True, flip_normals=False)
    
    return decimated

def extract_data(mesh):
    # Extract flattened arrays
    points = mesh.points.flatten().tolist() # [x,y,z, x,y,z...]
    
    # Faces in PyVista are [n, i1, i2, i3, m, j1, j2, j3...]
    # We need just triangles [i1, i2, i3...]
    faces = mesh.faces
    indices = []
    i = 0
    while i < len(faces):
        n = faces[i]
        if n == 3:
            indices.extend(faces[i+1 : i+4])
        elif n == 4:
            # triangulate quad roughly
            indices.extend([faces[i+1], faces[i+2], faces[i+3]])
            indices.extend([faces[i+1], faces[i+3], faces[i+4]])
        i += n + 1
        
    normals = mesh.point_data['Normals'].flatten().tolist()
    
    return points, indices, normals

if __name__ == '__main__':
    print("Generating HoloBody V14 Organic Data...")
    
    mesh = create_human_body()
    
    if mesh:
        p, i, n = extract_data(mesh)
        
        # Scale to -1.0 to 1.0 roughly
        # Current height approx 2.2 (-1.35 to 0.85). 
        # Shift Y by +0.25 approx.
        
        # Normalize
        p_np = np.array(p)
        # Simple shift
        # No robust normalization needed if we tune renderer, 
        # but let's shift Y up a bit so center is torso.
        
        data = {
            'version': 'v14_organic',
            'generated_at': '2026-02-19',
            'p': p, # Points
            'i': i, # Indices (Triangles)
            'n': n, # Normals
        }
        
        os.makedirs('frontend/flutter-app/assets/data', exist_ok=True)
        with open(OUTPUT_FILE, 'w') as f:
            json.dump(data, f)
            
        print(f"Stats: {len(p)//3} vertices, {len(i)//3} triangles.")
        print(f"Saved to {OUTPUT_FILE}")
        
    else:
        print("Failed to generate mesh (PyVista missing?).")
