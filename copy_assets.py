import os
import shutil

source_dir = r"C:\Users\kangj\.gemini\antigravity\brain\1a430321-7cea-49c2-85a5-ccc56667317c"
target_dir = r"\\wsl.localhost\Ubuntu\home\kangjh3kang\Manpasik\frontend\flutter-app\assets\images"

files = {
    "header_3d_frame_slim_v5_1771833070093.png": "header_3d_frame_slim.png",
    "bottom_3d_dock_slim_v5_1771833096268.png": "bottom_3d_dock_slim.png",
    "btn_3d_action_slim_v5_1771833141426.png": "btn_3d_action_slim.png"
}

try:
    if not os.path.exists(target_dir):
        os.makedirs(target_dir)
        print(f"Created directory: {target_dir}")

    for src_name, tgt_name in files.items():
        src_path = os.path.join(source_dir, src_name)
        tgt_path = os.path.join(target_dir, tgt_name)
        if os.path.exists(src_path):
            shutil.copy2(src_path, tgt_path)
            print(f"Copied: {src_name} -> {tgt_name}")
        else:
            print(f"Source not found: {src_path}")
    print("All tasks completed successfully (Python).")
except Exception as e:
    print(f"Error occurred: {e}")
