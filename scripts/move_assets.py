import shutil
import os

# Autopilot: Moving generated assets to Flutter project
src_dir = r"C:\Users\kangj\Documents\brain\1a430321-7cea-49c2-85a5-ccc56667317c" # Actually current brain dir
# Wait, I don't know the exact absolute path of the brain dir in the python env easily without providing it.
# I'll use the paths I know from the tool.

src_frame = r"C:\Users\kangj\.gemini\antigravity\brain\1a430321-7cea-49c2-85a5-ccc56667317c\sanggam_frame_v3_1771812355845.png"
src_button = r"C:\Users\kangj\.gemini\antigravity\brain\1a430321-7cea-49c2-85a5-ccc56667317c\jagae_button_v3_1771812373781.png"

dest_dir = r"\\wsl.localhost\Ubuntu\home\kangjh3kang\Manpasik\frontend\flutter-app\assets\images"

def move_assets():
    if not os.path.exists(dest_dir):
        os.makedirs(dest_dir)
    
    print(f"[Autopilot] Moving {src_frame} to {dest_dir}...")
    shutil.copy2(src_frame, os.path.join(dest_dir, "sanggam_frame.png"))
    
    print(f"[Autopilot] Moving {src_button} to {dest_dir}...")
    shutil.copy2(src_button, os.path.join(dest_dir, "jagae_button.png"))
    
    print("[Autopilot] SUCCESS: Assets moved successfully.")

if __name__ == "__main__":
    move_assets()
