import json
import math
import os
import numpy as np

# Configuration
OUTPUT_FILE = 'frontend/flutter-app/assets/data/holo_body_v13.json'

# Data Storage
# Vertices: [x, y, z, x, y, z...]
# Indices: [i1, i2, i3, i4...] (pairs of indices for GL_LINES)
vertices = []
indices = []
offset = 0

def add_line_loop(points, close=True):
    global offset
    # points: [[x,y,z], [x,y,z]...]
    
    start_wrapper_idx = len(vertices) // 3
    
    for p in points:
        vertices.extend([round(float(p[0]), 4), round(float(p[1]), 4), round(float(p[2]), 4)])
        
    count = len(points)
    for i in range(count - 1):
        indices.extend([start_wrapper_idx + i, start_wrapper_idx + i + 1])
    
    if close:
        indices.extend([start_wrapper_idx + count - 1, start_wrapper_idx])

def add_grid_mesh(rows_of_points):
    # rows_of_points: List of lists of points (e.g. latitude rings)
    # Connects row i to row i+1
    
    # First, separate start indices for each row
    row_start_indices = []
    current_idx = len(vertices) // 3
    
    for row in rows_of_points:
        row_start_indices.append(current_idx)
        for p in row:
            vertices.extend([round(float(p[0]), 4), round(float(p[1]), 4), round(float(p[2]), 4)])
        current_idx += len(row)
        
    # Draw horizontal lines (rings)
    for r_idx, row in enumerate(rows_of_points):
        base = row_start_indices[r_idx]
        count = len(row)
        for i in range(count):
            next_i = (i + 1) % count
            indices.extend([base + i, base + next_i])
            
    # Draw vertical lines (longitude)
    if len(rows_of_points) > 1:
        for r_idx in range(len(rows_of_points) - 1):
            base_curr = row_start_indices[r_idx]
            base_next = row_start_indices[r_idx+1]
            
            count_curr = len(rows_of_points[r_idx])
            count_next = len(rows_of_points[r_idx+1])
            
            # Assume equal count for best grid
            limit = min(count_curr, count_next)
            for i in range(limit):
                indices.extend([base_curr + i, base_next + i])


# -----------------------------------------------------------------------------
# Parametric Generators (NumPy / Math)
# -----------------------------------------------------------------------------

def generate_ellipsoid_grid(center, axes, v_res=15, h_res=24, partial_v=1.0):
    # center: (x,y,z), axes: (rx,ry,rz)
    # v_res: vertical rings (latitude)
    # h_res: horizontal segments (longitude)
    
    rows = []
    
    # Generate rings from bottom to top
    for i in range(v_res + 1):
        v = i / v_res # 0 to 1
        phi = math.acos(2 * v - 1) # Spherical mapping 0 to pi
        
        # Adjust for partial (e.g. hemisphere) if needed, but here full sphere usually
        
        ring = []
        # Ring radius at this height
        # x = rx * sin(phi) * cos(theta)
        # y = ry * sin(phi) * sin(theta)
        # z = rz * cos(phi)
        
        for j in range(h_res):
            theta = 2 * math.pi * (j / h_res)
            
            x = center[0] + axes[0] * math.sin(phi) * math.cos(theta)
            y = center[1] + axes[1] * math.sin(phi) * math.sin(theta)
            z = center[2] + axes[2] * math.cos(phi)
            ring.append([x,y,z])
        rows.append(ring)
        
    add_grid_mesh(rows)

def generate_cylinder_grid(start, end, r_start, r_end, rings=10, segments=12):
    p_start = np.array(start)
    p_end = np.array(end)
    vec = p_end - p_start
    length = np.linalg.norm(vec)
    
    if length < 0.001: return

    axis = vec / length
    if abs(axis[0]) < 0.9: temp = np.array([1,0,0])
    else: temp = np.array([0,1,0])
    u = np.cross(axis, temp)
    u = u / np.linalg.norm(u)
    v = np.cross(axis, u)
    
    rows = []
    for i in range(rings + 1):
        t = i / rings
        center = p_start + vec * t
        r = r_start + (r_end - r_start) * t
        
        ring = []
        for j in range(segments):
            theta = 2 * math.pi * (j / segments)
            pos = center + r * (math.cos(theta) * u + math.sin(theta) * v)
            ring.append([pos[0], pos[1], pos[2]])
        rows.append(ring)
        
    add_grid_mesh(rows)

def generate_skeleton_hand_mesh(wrist, is_left):
     dir_x = -1.0 if is_left else 1.0
     palm_center = np.array([wrist[0] + dir_x * 0.02, wrist[1] + 0.04, 0])
     
     # Palm Grid (Diamond shape)
     palm_pts = [
         [wrist[0], wrist[1], 0],
         [wrist[0] + dir_x * 0.04, wrist[1] + 0.02, 0],
         [palm_center[0], palm_center[1] + 0.02, 0],
         [wrist[0] - dir_x * 0.01, wrist[1] + 0.03, 0]
     ]
     add_line_loop(palm_pts, True)
     
     # Fingers (Line Segments with joints)
     for i in range(5):
        angle = -0.3 + i * 0.15 
        if is_left: angle *= -1
        fl = 0.06 + 0.04 * math.sin(i * math.pi / 4)
        
        tip = np.array([
            palm_center[0] + math.sin(angle) * fl * 3.0 * dir_x,
            palm_center[1] + math.cos(angle) * fl,
            0
        ])
        
        # Segmented finger for "Tech" look
        joints = 3
        prev = palm_center
        for k in range(1, joints + 1):
            curr = palm_center + (tip - palm_center) * (k/joints)
            
            # Add line
            # Instead of separate lines, we want connectivity. 
            # But line loop helper is good enough.
            add_line_loop([prev, curr], False)
            
            # Joint marker (small cross or box)
            # Simple small cross
            s = 0.005
            add_line_loop([[curr[0]-s, curr[1], curr[2]], [curr[0]+s, curr[1], curr[2]]], False)
            
            prev = curr

# -----------------------------------------------------------------------------
# Main Builder
# -----------------------------------------------------------------------------

def build_body_wireframe():
    # HEAD (Refined Grid)
    # Scaled to -1.0 ~ 1.0 coords
    generate_ellipsoid_grid((0, -0.85, 0), (0.10, 0.13, 0.11), v_res=16, h_res=20)
    
    # NECK
    generate_cylinder_grid((0, -0.73, 0), (0, -0.65, 0), 0.06, 0.065, rings=3, segments=12)
    
    # TORSO
    # Chest
    generate_ellipsoid_grid((0, -0.52, 0), (0.24, 0.18, 0.13), v_res=10, h_res=20)
    # Abs/Waist
    generate_cylinder_grid((0, -0.50, 0), (0, -0.25, 0), 0.18, 0.15, rings=5, segments=16)
    # Hips
    generate_ellipsoid_grid((0, -0.15, 0), (0.21, 0.15, 0.13), v_res=8, h_res=20)
    
    # ARMS
    s_w = 0.28
    # L
    generate_ellipsoid_grid((-s_w, -0.60, 0), (0.09, 0.09, 0.09), v_res=6, h_res=12) # Shoulder
    generate_cylinder_grid((-s_w, -0.60, 0), (-s_w - 0.05, -0.35, 0), 0.07, 0.06, rings=4, segments=8) # Upper
    generate_cylinder_grid((-s_w - 0.05, -0.35, 0), (-s_w - 0.08, -0.10, 0), 0.055, 0.04, rings=4, segments=8) # Lower
    
    # R
    generate_ellipsoid_grid((s_w, -0.60, 0), (0.09, 0.09, 0.09), v_res=6, h_res=12)
    generate_cylinder_grid((s_w, -0.60, 0), (s_w + 0.05, -0.35, 0), 0.07, 0.06, rings=4, segments=8)
    generate_cylinder_grid((s_w + 0.05, -0.35, 0), (s_w + 0.08, -0.10, 0), 0.055, 0.04, rings=4, segments=8)
    
    # HANDS
    generate_skeleton_hand_mesh((-s_w - 0.08, -0.08), True)
    generate_skeleton_hand_mesh((s_w + 0.08, -0.08), False)
    
    # LEGS
    l_off = 0.12
    # L
    generate_cylinder_grid((-l_off, -0.15, 0), (-l_off - 0.02, 0.30, 0), 0.09, 0.07, rings=6, segments=12) # Thigh
    generate_ellipsoid_grid((-l_off - 0.02, 0.30, 0), (0.07, 0.07, 0.07), v_res=4, h_res=12) # Knee
    generate_cylinder_grid((-l_off - 0.02, 0.30, 0), (-l_off - 0.02, 0.70, 0), 0.065, 0.05, rings=6, segments=12) # Calf
    
    # R
    generate_cylinder_grid((l_off, -0.15, 0), (l_off + 0.02, 0.30, 0), 0.09, 0.07, rings=6, segments=12)
    generate_ellipsoid_grid((l_off + 0.02, 0.30, 0), (0.07, 0.07, 0.07), v_res=4, h_res=12)
    generate_cylinder_grid((l_off + 0.02, 0.30, 0), (l_off + 0.02, 0.70, 0), 0.065, 0.05, rings=6, segments=12)
    
    # FEET
    generate_ellipsoid_grid((-l_off - 0.02, 0.73, 0.05), (0.06, 0.04, 0.12), v_res=5, h_res=12)
    generate_ellipsoid_grid((l_off + 0.02, 0.73, 0.05), (0.06, 0.04, 0.12), v_res=5, h_res=12)

if __name__ == '__main__':
    print("Generating HoloBody V13 Wireframe Mesh...")
    
    # Try using PyVista if available (for future extensions), 
    # but currently mostly relying on numpy parametric generation for Grid control.
    # Why? PyVista decimate often creates chaotic triangles. 
    # We want structured LAT/LONG grids which look much more "Sci-Fi".
    # So the Grid Generator above is actually SUPERIOR for this specific aesthetic than raw STL meshing.
    
    build_body_wireframe()
    
    data = {
        'version': 'v13_wireframe',
        'generated_at': '2026-02-19',
        'p': vertices, # Flattened [x,y,z, x,y,z...]
        'i': indices   # Flattened [0,1, 1,2...]
    }
    
    os.makedirs('frontend/flutter-app/assets/data', exist_ok=True)
    with open(OUTPUT_FILE, 'w') as f:
        json.dump(data, f)
        
    print(f"Stats: {len(vertices)//3} vertices, {len(indices)//2} lines.")
    print(f"Saved to {OUTPUT_FILE}")
