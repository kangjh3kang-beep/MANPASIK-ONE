// harness/cartridge_manifest.rs

#[derive(Debug, Default, Clone)]
pub struct CartridgeManifest {
    pub csi_version: u8,
    pub stage: u8,
    // Add other relevant fields...
}

impl CartridgeManifest {
    pub fn is_v1(&self) -> bool {
        // v1.0 = 16핀 CSI 로직
        self.csi_version == 1
    }

    pub fn is_v2(&self) -> bool {
        self.csi_version == 2
    }
}
