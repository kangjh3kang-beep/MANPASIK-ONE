package hl7mllp

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

// Client 는 외부 LIS/EHR 로 HL7 메시지를 송신하는 MLLP 클라이언트.
//
// 만파식이 내부 측정 결과를 외부 EHR 에 푸시할 때 사용 (Bundle → HL7v2 역변환
// 후 송신). 단일 TCP 연결을 재사용하여 여러 메시지를 보낼 수 있습니다.
type Client struct {
	addr            string
	dialTimeout     time.Duration
	readTimeout     time.Duration
	writeTimeout    time.Duration

	mu     sync.Mutex
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

// ClientConfig 는 Client 설정.
type ClientConfig struct {
	Addr         string        // 예: "lab.hospital.org:2575"
	DialTimeout  time.Duration // 기본 10초
	ReadTimeout  time.Duration // ACK 수신 타임아웃 (기본 30초)
	WriteTimeout time.Duration // 송신 타임아웃 (기본 10초)
}

// NewClient 생성 — Addr 필수. 실제 dial 은 첫 Send 또는 Connect() 시 수행.
func NewClient(cfg ClientConfig) (*Client, error) {
	if cfg.Addr == "" {
		return nil, errors.New("Addr 필수")
	}
	if cfg.DialTimeout <= 0 {
		cfg.DialTimeout = 10 * time.Second
	}
	if cfg.ReadTimeout <= 0 {
		cfg.ReadTimeout = 30 * time.Second
	}
	if cfg.WriteTimeout <= 0 {
		cfg.WriteTimeout = 10 * time.Second
	}
	return &Client{
		addr:         cfg.Addr,
		dialTimeout:  cfg.DialTimeout,
		readTimeout:  cfg.ReadTimeout,
		writeTimeout: cfg.WriteTimeout,
	}, nil
}

// Connect 는 명시적으로 TCP 연결. 이미 연결되어 있으면 noop.
func (c *Client) Connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn != nil {
		return nil
	}
	conn, err := net.DialTimeout("tcp", c.addr, c.dialTimeout)
	if err != nil {
		return fmt.Errorf("dial %s: %w", c.addr, err)
	}
	c.conn = conn
	c.reader = bufio.NewReader(conn)
	c.writer = bufio.NewWriter(conn)
	return nil
}

// Close 는 연결 종료.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.conn == nil {
		return nil
	}
	err := c.conn.Close()
	c.conn = nil
	c.reader = nil
	c.writer = nil
	return err
}

// Send 는 raw HL7 메시지 송신 + ACK 응답 수신.
//
// 반환:
//   - ackCode: "AA"/"AE"/"AR"/"" (빈 문자열 = ACK 파싱 실패)
//   - ackText: ACK 본문 (MSA-3)
//   - err: 송신/수신 에러
func (c *Client) Send(raw []byte) (ackCode, ackText string, err error) {
	if err := c.Connect(); err != nil {
		return "", "", err
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	_ = c.conn.SetWriteDeadline(time.Now().Add(c.writeTimeout))
	if werr := WriteFrame(c.writer, raw); werr != nil {
		return "", "", fmt.Errorf("write: %w", werr)
	}

	_ = c.conn.SetReadDeadline(time.Now().Add(c.readTimeout))
	ackRaw, rerr := ReadFrame(c.reader)
	if rerr != nil {
		return "", "", fmt.Errorf("ACK 수신: %w", rerr)
	}
	ackCode, ackText = ParseACK(string(ackRaw))
	if ackCode == "" {
		return "", "", errors.New("ACK 파싱 실패")
	}
	return ackCode, ackText, nil
}

// SendString 은 Send 의 string 편의 래퍼.
func (c *Client) SendString(raw string) (ackCode, ackText string, err error) {
	return c.Send([]byte(raw))
}
