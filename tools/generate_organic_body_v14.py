import numpy as np
import pyvista as pv
import json
import math
import os
import time

# Configuration
OUTPUT_DIR = 'frontend/flutter-app/assets/data'
OUTPUT_FILE = 'holo_body_v16.json' # V16
GRID_RES = 135

def create_high_fidelity_body():
    """
    Generates a high-fidelity organic human body mesh using SDF.
    V16: Human Proportions (Not Alien)
    """
    print(f"Initializing V16 Body Generator (Resolution: {GRID_RES})...")
    start_time = time.time()

    # Grid Setup (Head Up = Negative Y)
    x = np.linspace(-0.6, 0.6, GRID_RES)
    y = np.linspace(-1.0, 1.1, int(GRID_RES * 2.1)) 
    z = np.linspace(-0.35, 0.35, int(GRID_RES * 0.7))
    
    grid = pv.RectilinearGrid(x, y, z)
    X, Y, Z = np.meshgrid(x, y, z, indexing='ij')
    
    # -------------------------------------------------------------------------
    # SDF Primitives
    # -------------------------------------------------------------------------
    def length(vx, vy, vz):
        return np.sqrt(vx*vx + vy*vy + vz*vz)

    def mix(a, b, t):
        return a * (1.0 - t) + b * t

    def sdSphere(px, py, pz, r):
        return length(px, py, pz) - r

    def sdEllipsoid(px, py, pz, r):
        k0 = length(px/r[0], py/r[1], pz/r[2])
        k1 = length(px/(r[0]*r[0]), py/(r[1]*r[1]), pz/(r[2]*r[2]))
        return k0 * (k0 - 1.0) / k1

    def sdCapsule(px, py, pz, ax, ay, az, bx, by, bz, r):
        pax, pay, paz = px-ax, py-ay, pz-az
        bax, bay, baz = bx-ax, by-ay, bz-az
        h = np.clip((pax*bax + pay*bay + paz*baz) / (bax*bax + bay*bay + baz*baz + 1e-6), 0.0, 1.0)
        return length(pax - bax * h, pay - bay * h, paz - baz * h) - r

    def sdTaperedCapsule(px, py, pz, ax, ay, az, bx, by, bz, ra, rb):
        pax, pay, paz = px-ax, py-ay, pz-az
        bax, bay, baz = bx-ax, by-ay, bz-az
        h = np.clip((pax*bax + pay*bay + paz*baz) / (bax*bax + bay*bay + baz*baz + 1e-6), 0.0, 1.0)
        return length(pax - bax * h, pay - bay * h, paz - baz * h) - (ra * (1.0 - h) + rb * h)

    def smin(a, b, k=0.08):
        h = np.clip(0.5 + 0.5 * (b - a) / k, 0.0, 1.0)
        return mix(b, a, h) - k * h * (1.0 - h)

    # -------------------------------------------------------------------------
    # Anatomy Modeling
    # -------------------------------------------------------------------------
    d = np.full_like(X, 10.0)

    # --- V19: PERFECT HUMAN ANATOMY (7.5 Heads) ---
    # Reference: Head Height ~ 0.22 units. Total Height ~ 1.75 units.
    # Center (0,0,0) is roughly the navel/waist.
    
    # --- HEAD & NECK ---
    # Head (Cranium + Jaw) - Proportionally smaller than V14
    d = smin(d, sdEllipsoid(X, Y-(-0.92), Z-0.02, [0.07, 0.088, 0.078]), 0.04)
    # Jawline (Defined)
    d = smin(d, sdTaperedCapsule(X,Y,Z, 0, -0.90, 0.03, 0, -0.80, 0.06, 0.055, 0.038), 0.03)
    # Neck (Muscular Traps transition)
    d = smin(d, sdCapsule(X,Y,Z, 0, -0.78, -0.02, 0, -0.68, -0.02, 0.048), 0.05)
    
    # --- TORSO (V-Taper) ---
    # Traps 
    d = smin(d, sdTaperedCapsule(X,Y,Z, -0.12, -0.72, -0.03, 0.12, -0.72, -0.03, 0.02, 0.02), 0.08)
    
    # Chest (Pectorals - Wide & Flat)
    pec_y = -0.60
    pec_z = 0.04
    for sign in [-1, 1]:
        d = smin(d, sdEllipsoid(X-(sign*0.09), Y-pec_y, Z-pec_z, [0.085, 0.065, 0.03]), 0.06)

    # Ribcage / Lats (V-Shape Base)
    d = smin(d, sdTaperedCapsule(X,Y,Z, 0, -0.65, -0.04, 0, -0.35, -0.02, 0.16, 0.11), 0.10)
    
    # Abdomen (Six Pack Definition)
    d = smin(d, sdCapsule(X,Y,Z, 0, -0.45, 0.07, 0, -0.22, 0.06, 0.065), 0.04)
    for y_div in [-0.40, -0.32, -0.24]: # Abs separations
         d = smin(d, sdCapsule(X,Y,Z, -0.05, y_div, 0.07, 0.05, y_div, 0.07, 0.006), 0.03)

    # --- ARMS (Muscular Definition) ---
    for sign in [-1, 1]:
        # Clavicle (Collarbone)
        d = smin(d, sdCapsule(X,Y,Z, 0, -0.70, 0.02, sign*0.20, -0.71, 0.0, 0.022), 0.03)

        # Deltoids (Shoulders - Round & Broad)
        p_shoulder = (sign*0.23, -0.69, -0.01)
        d = smin(d, sdEllipsoid(X-p_shoulder[0], Y-p_shoulder[1], Z-p_shoulder[2], [0.085, 0.085, 0.08]), 0.05)
        
        # Biceps/Triceps (Upper Arm)
        p_elbow = (sign*0.28, -0.36, 0.03)
        d = smin(d, sdTaperedCapsule(X,Y,Z, *p_shoulder, *p_elbow, 0.065, 0.055), 0.04)
        
        # Forearm (Tapered)
        p_wrist = (sign*0.31, -0.08, 0.10)
        d = smin(d, sdTaperedCapsule(X,Y,Z, *p_elbow, *p_wrist, 0.052, 0.038), 0.04)
        
        # Hands (Generic shape for scale)
        d = smin(d, sdEllipsoid(X-(sign*0.315), Y-(-0.03), Z-0.11, [0.035, 0.05, 0.02]), 0.03)

    # --- LEGS (Athletic) ---
    for sign in [-1, 1]:
        p_hip = (sign*0.12, -0.12, 0.0)
        p_knee = (sign*0.14, 0.42, 0.03)
        p_ankle = (sign*0.15, 0.90, -0.03)
        
        # Glutes (Buttocks)
        d = smin(d, sdSphere(X-(sign*0.12), Y-(-0.12), Z-(-0.06), 0.11), 0.08)
        
        # Thighs (Quads/Hamstrings - Thick at top)
        d = smin(d, sdTaperedCapsule(X,Y,Z, *p_hip, *p_knee, 0.12, 0.075), 0.06)
        
        # Knees
        d = smin(d, sdSphere(X-p_knee[0], Y-p_knee[1], Z-p_knee[2], 0.055), 0.03)
        
        # Calves (Defined Gastrocnemius)
        d = smin(d, sdTaperedCapsule(X,Y,Z, *p_knee, *p_ankle, 0.065, 0.045), 0.05)
        calf_muscle = (p_knee[0]*0.7 + p_ankle[0]*0.3, p_knee[1]*0.7 + p_ankle[1]*0.3, p_knee[2] - 0.04)
        d = smin(d, sdEllipsoid(X-calf_muscle[0], Y-calf_muscle[1], Z-calf_muscle[2], [0.048, 0.09, 0.035]), 0.06)

        # Feet
        d = smin(d, sdEllipsoid(X-(p_ankle[0]+sign*0.02), Y-(p_ankle[1]+0.12), Z-(p_ankle[2]+0.06), [0.04, 0.02, 0.07]), 0.04)

    print(f"SDF Calculation done in {time.time() - start_time:.2f}s")
    
    # -------------------------------------------------------------------------
    # Generate Low-Poly Mesh (For Performance)
    # -------------------------------------------------------------------------
    print("Extracting surface (Marching Cubes)...")
    grid.point_data["values"] = d.flatten(order='F')
    mesh = grid.contour(isosurfaces=[0.0])
    
    # Decimate to reasonable count (e.g. 5000-8000 polys for mobile/web)
    target_tris = 12000 
    if mesh.n_cells > target_tris:
        print(f"Decimating {mesh.n_cells} -> {target_tris}...")
        mesh = mesh.decimate_pro(1.0 - (target_tris / mesh.n_cells))
    
    print("Smoothing & Computing Normals...")
    mesh = mesh.smooth(n_iter=20, relaxation_factor=0.08) # More smooth for organic look
    mesh = mesh.compute_normals(cell_normals=False, point_normals=True, auto_orient_normals=True)
    
    # Export V19 Data
    vertices = mesh.points.flatten().tolist()
    
    # Safe Face Extraction
    try:
        raw_faces = mesh.faces
        indices = []
        i = 0
        while i < len(raw_faces):
            n = raw_faces[i]
            if n == 3:
                indices.extend(raw_faces[i+1:i+1+n])
            i += n + 1
    except:
        indices = mesh.faces.reshape(-1, 4)[:, 1:4].flatten().tolist()
        
    normals = mesh.point_data["Normals"].flatten().tolist()
    
    # Output File: holo_body_v19.json
    OUTPUT_FILE_V19 = 'holo_body_v19.json'
    
    data = {
        "version": "v19_perfect_human",
        "vertexCount": len(vertices) // 3,
        "faceCount": len(indices) // 3,
        "vertices": [round(float(v), 5) for v in vertices],
        "indices": [int(i) for i in indices],
        "normals": [round(float(n), 4) for n in normals]
    }
    
    os.makedirs(OUTPUT_DIR, exist_ok=True)
    out_path = os.path.join(OUTPUT_DIR, OUTPUT_FILE_V19)
    
    with open(out_path, 'w') as f:
        json.dump(data, f)
        
    print(f"Success! V19 Perfect Human Body saved to {out_path}")

if __name__ == "__main__":
    try:
        create_high_fidelity_body()
    except Exception:
        import traceback
        traceback.print_exc()
