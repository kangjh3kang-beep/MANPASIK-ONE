import json
import math
import random
import numpy as np
import os

# Configuration
OUTPUT_FILE = 'manpasik_body_v11.json'
TOTAL_POINTS_TARGET = 30000

# Data Storage
particles = []
outlines = []

def add_particle(x, y, z, size=1.0, alpha=0.5):
    particles.append({
        'x': round(float(x), 4),
        'y': round(float(y), 4),
        'z': round(float(z), 4),
        's': round(float(size), 2),
        'a': round(float(alpha), 2)
    })

def add_outline(points):
    # points is list of [x, y, z]
    line = []
    for p in points:
        line.append({
            'x': round(float(p[0]), 4),
            'y': round(float(p[1]), 4),
            'z': round(float(p[2]), 4)
        })
    outlines.append(line)

# -----------------------------------------------------------------------------
# Generators (SDF / Parametric)
# -----------------------------------------------------------------------------

def generate_ellipsoid(center, axes, count, shell_ratio=0.8):
    # center: (x, y, z), axes: (rx, ry, rz)
    # shell_ratio: 0.0 to 1.0 (1.0 means all points on surface)
    print(f"Generating Ellipsoid at {center} with {count} points...")
    
    # Generate random points in unit sphere
    # u, v, w uniform in [0,1]
    # To get uniform in sphere:
    # theta = 2*pi*u
    # phi = acos(2*v - 1)
    # r = w^(1/3)
    
    for _ in range(count):
        theta = 2 * math.pi * np.random.rand()
        phi = math.acos(2 * np.random.rand() - 1)
        
        r_raw = np.random.rand()
        # Shell distribution logic: curve r towards 1.0
        # e.g. r = r^(1/N)
        r = math.pow(r_raw, 1.0/ (1.0 + 10.0 * shell_ratio)) 
        
        # Apply scaling
        x = center[0] + axes[0] * r * math.sin(phi) * math.cos(theta)
        y = center[1] + axes[1] * r * math.sin(phi) * math.sin(theta)
        z = center[2] + axes[2] * r * math.cos(phi)
        
        add_particle(x, y, z)

def generate_limb_segment(start, end, r_start, r_end, count):
    # Cylinder/Cone segment
    print(f"Generating Limb from {start} to {end}...")
    
    p_start = np.array(start)
    p_end = np.array(end)
    vec = p_end - p_start
    length = np.linalg.norm(vec)
    axis = vec / length if length > 0 else np.array([0,1,0])
    
    # Find perpendicular vectors for circle
    # temp vector not parallel to axis
    if abs(axis[0]) < 0.9: temp = np.array([1,0,0])
    else: temp = np.array([0,1,0])
    
    u = np.cross(axis, temp)
    u = u / np.linalg.norm(u)
    v = np.cross(axis, u)
    
    for _ in range(count):
        t = np.random.rand() # 0 to 1 along axis
        
        theta = 2 * math.pi * np.random.rand()
        rad_t = r_start + (r_end - r_start) * t
        
        # Surface bias
        r = rad_t * math.pow(np.random.rand(), 1.0/6.0)
        
        # Circle pos
        circle_pos = r * (math.cos(theta) * u + math.sin(theta) * v)
        # Axis pos
        axis_pos = p_start + vec * t
        
        final_pos = axis_pos + circle_pos
        add_particle(final_pos[0], final_pos[1], final_pos[2])

def generate_skeleton_hand(wrist, is_left):
    # Skeleton lines
    dir_x = -1.0 if is_left else 1.0
    
    palm_len = 0.08
    finger_len = [0.06, 0.09, 0.10, 0.085, 0.05] # Thumb to Pinky? No, Pinky, Ring, Mid, Index, Thumb usually
    # Let's simple fan
    
    # Palm Center
    palm_center = np.array([wrist[0] + dir_x * 0.02, wrist[1] + 0.04, 0])
    
    # Draw Fingers
    for i in range(5):
        # i=0(Pinky) to 4(Thumb)
        angle = -0.3 + i * 0.15 # Radians spread
        if is_left: angle *= -1
        
        fl = 0.06 + 0.04 * math.sin(i * math.pi / 4) # Curve length
        
        tip = np.array([
            palm_center[0] + math.sin(angle) * fl * 3.0 * dir_x, # Spread X
            palm_center[1] + math.cos(angle) * fl, # Y
            0
        ])
        
        # Add Line
        add_outline([
            [wrist[0], wrist[1], 0],
            [palm_center[0], palm_center[1], 0],
            [tip[0], tip[1], 0]
        ])
        
        # Add particles along bone
        segments = 10
        for k in range(segments):
            t = k/segments
            p = palm_center + (tip - palm_center) * t
            add_particle(p[0], p[1], p[2], size=0.4, alpha=0.9)

# -----------------------------------------------------------------------------
# Main Generation Logic (Metric Scaled to -1.0 to 1.0)
# -----------------------------------------------------------------------------
def build_body():
    # HEAD
    generate_ellipsoid((0, -0.85, 0), (0.10, 0.13, 0.11), 3000, shell_ratio=0.9)
    # Jaw Line
    jaw_pts = []
    for i in range(21):
        t = i/20 * math.pi # 0 to pi
        x = 0.09 * math.cos(t) # Width
        z = 0.10 * math.sin(t) # Depth curve
        y = -0.78 + 0.05 * abs(x) # Chin dip
        jaw_pts.append([x, y, z])
    add_outline(jaw_pts)
    
    # NECK
    generate_limb_segment((0, -0.75, 0), (0, -0.65, 0), 0.06, 0.065, 800)
    
    # CHEST
    generate_ellipsoid((0, -0.55, 0), (0.24, 0.15, 0.12), 5000, shell_ratio=0.7)
    
    # ABDOMEN
    generate_ellipsoid((0, -0.35, 0), (0.19, 0.14, 0.11), 3000, shell_ratio=0.6)
    
    # HIPS
    generate_ellipsoid((0, -0.15, 0), (0.21, 0.15, 0.13), 3000, shell_ratio=0.7)
    
    # ARMS
    # L Shoulder
    generate_ellipsoid((-0.26, -0.60, 0), (0.09, 0.09, 0.09), 800)
    # L Upper
    generate_limb_segment((-0.26, -0.60, 0), (-0.32, -0.35, 0), 0.07, 0.06, 1200)
    # L Elbow
    generate_ellipsoid((-0.32, -0.35, 0), (0.06, 0.06, 0.06), 400)
    # L Lower
    generate_limb_segment((-0.32, -0.35, 0), (-0.35, -0.10, 0), 0.055, 0.04, 1200)
    
    # R Shoulder
    generate_ellipsoid((0.26, -0.60, 0), (0.09, 0.09, 0.09), 800)
    # R Upper
    generate_limb_segment((0.26, -0.60, 0), (0.32, -0.35, 0), 0.07, 0.06, 1200)
    # R Elbow
    generate_ellipsoid((0.32, -0.35, 0), (0.06, 0.06, 0.06), 400)
    # R Lower
    generate_limb_segment((0.32, -0.35, 0), (0.35, -0.10, 0), 0.055, 0.04, 1200)
    
    # HANDS
    generate_skeleton_hand((-0.35, -0.08), True)
    generate_skeleton_hand((0.35, -0.08), False)
    
    # LEGS
    # L Thigh
    generate_limb_segment((-0.12, -0.15, 0), (-0.14, 0.30, 0), 0.09, 0.07, 2000)
    # L Knee
    generate_ellipsoid((-0.14, 0.30, 0), (0.07, 0.07, 0.07), 500)
    # L Calf
    generate_limb_segment((-0.14, 0.30, 0), (-0.14, 0.70, 0), 0.065, 0.05, 1800)
    
    # R Thigh
    generate_limb_segment((0.12, -0.15, 0), (0.14, 0.30, 0), 0.09, 0.07, 2000)
    # R Knee
    generate_ellipsoid((0.14, 0.30, 0), (0.07, 0.07, 0.07), 500)
    # R Calf
    generate_limb_segment((0.14, 0.30, 0), (0.14, 0.70, 0), 0.065, 0.05, 1800)
    
    # FEET
    # L Foot
    generate_ellipsoid((-0.14, 0.73, 0.05), (0.06, 0.04, 0.12), 800)
    # R Foot
    generate_ellipsoid((0.14, 0.73, 0.05), (0.06, 0.04, 0.12), 800)


if __name__ == '__main__':
    print("Initializing HoloBody V11 Generator...")
    build_body()
    
    print(f"Generated {len(particles)} particles and {len(outlines)} outline paths.")
    
    data = {
        'version': 'v11',
        'generated_at': '2026-02-19',
        'particles': particles,
        'outlines': outlines
    }
    
    # Ensure directory exists
    os.makedirs('frontend/flutter-app/assets/data', exist_ok=True)
    
    with open('frontend/flutter-app/assets/data/holo_body_v11.json', 'w') as f:
        json.dump(data, f)
        
    print(f"Data saved to frontend/flutter-app/assets/data/holo_body_v11.json")
