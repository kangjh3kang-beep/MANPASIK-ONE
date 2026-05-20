// BLE 및 UART를 통해 하드웨어(STM32)에서 들어오는 원시 바이트(Raw Bytes) 스트림을 파싱하여
// 16핀(CSI v1.0) 차동 신호 프레임으로 복원하는 모듈입니다. (Phase 1 구현)

#[derive(Debug, PartialEq)]
pub enum ParserError {
    InvalidHeader,
    ChecksumMismatch,
    InsufficientData,
}

const HEADER_1: u8 = 0xAA;
const HEADER_2: u8 = 0x55;
const FRAME_PAYLOAD_SIZE: usize = 64; // 16채널 * 4바이트(f32)
const FRAME_TOTAL_SIZE: usize = 2 + 1 + FRAME_PAYLOAD_SIZE + 1; // Header(2) + Status(1) + Payload(64) + Checksum(1) = 68 bytes

pub struct BleStreamParser {
    buffer: Vec<u8>,
}

impl BleStreamParser {
    pub fn new() -> Self {
        Self {
            buffer: Vec::with_capacity(256),
        }
    }

    /// Flutter(FFI)나 하드웨어 UART로부터 들어오는 원시 바이트들을 버퍼에 누적합니다.
    pub fn feed_bytes(&mut self, data: &[u8]) {
        self.buffer.extend_from_slice(data);
    }

    /// 완성된 프레임이 존재하는지 검사하고 16채널 f32 배열을 추출합니다.
    /// 에러 발생 시 Result<T, E> 원칙(CLAUDE.md)에 따라 에러를 격리 및 반환합니다.
    pub fn try_parse_frame(&mut self) -> Result<Option<[f32; 16]>, ParserError> {
        while self.buffer.len() >= FRAME_TOTAL_SIZE {
            // 헤더 검사
            if self.buffer[0] != HEADER_1 || self.buffer[1] != HEADER_2 {
                self.buffer.remove(0); // 1바이트씩 쉬프트하며 헤더 동기화 서치
                continue;
            }

            // 페이로드 추출 및 체크섬 검증
            let payload = &self.buffer[3..3 + FRAME_PAYLOAD_SIZE];
            let received_checksum = self.buffer[3 + FRAME_PAYLOAD_SIZE];

            let calculated_checksum = payload.iter().fold(0, |acc, &x| acc ^ x);

            if received_checksum != calculated_checksum {
                // 체크섬 실패 시, 해당 프레임 전체를 드롭하고 다음 프레임을 찾습니다.
                self.buffer.drain(0..FRAME_TOTAL_SIZE);
                return Err(ParserError::ChecksumMismatch);
            }

            // f32 어레이 복원
            let mut channels = [0.0f32; 16];
            for i in 0..16 {
                let start = i * 4;
                let bytes = [
                    payload[start],
                    payload[start + 1],
                    payload[start + 2],
                    payload[start + 3],
                ];
                channels[i] = f32::from_le_bytes(bytes);
            }

            // 프레임 소비 완료
            self.buffer.drain(0..FRAME_TOTAL_SIZE);
            return Ok(Some(channels));
        }

        // 버퍼에 아직 68바이트가 안 모임
        Ok(None)
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_valid_frame_parsing() {
        let mut parser = BleStreamParser::new();
        let mut mock_frame = vec![0; FRAME_TOTAL_SIZE];
        mock_frame[0] = HEADER_1;
        mock_frame[1] = HEADER_2;
        mock_frame[2] = 0x01; // status

        // 채널 0을 1.0f32 로 설정 (0x00 0x00 0x80 0x3F)
        let f_bytes = 1.0f32.to_le_bytes();
        mock_frame[3] = f_bytes[0];
        mock_frame[4] = f_bytes[1];
        mock_frame[5] = f_bytes[2];
        mock_frame[6] = f_bytes[3];

        let payload = &mock_frame[3..3 + FRAME_PAYLOAD_SIZE];
        mock_frame[3 + FRAME_PAYLOAD_SIZE] = payload.iter().fold(0, |acc, &x| acc ^ x); // Calculate CRC

        parser.feed_bytes(&mock_frame);

        let result = parser.try_parse_frame().unwrap();
        assert!(result.is_some());

        let channels = result.unwrap();
        assert_eq!(channels[0], 1.0f32);
        assert_eq!(channels[1], 0.0f32);
    }
}
