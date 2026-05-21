package hl7mllp

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	hl7 "github.com/manpasik/backend/shared/hl7-parser"
)

// MessageHandler 는 수신된 HL7 메시지를 처리하는 콜백.
//
// raw 는 MLLP framing 이 제거된 순수 HL7 메시지. msg 는 파싱된 결과 (nil 이면
// 파싱 실패 — handler 가 raw 만 보고 처리하거나 ACK_REJECT 반환).
//
// 반환값:
//   - ackCode: "AA" (accept) / "AE" (error) / "AR" (reject)
//   - textMsg: ACK 본문에 포함할 추가 정보 (선택)
//   - err: 호스트 측 에러 (서버가 로그/메트릭 기록용 — ACK 자체는 ackCode 사용)
type MessageHandler func(ctx context.Context, raw []byte, msg *hl7.Message) (ackCode, textMsg string, err error)

// ServerConfig 는 MLLPServer 설정.
type ServerConfig struct {
	// Addr 는 listen 주소 (예: "0.0.0.0:2575" — HL7 표준 포트).
	Addr string

	// ReadTimeout 은 단일 메시지 읽기 타임아웃 (기본 30초).
	ReadTimeout time.Duration

	// WriteTimeout 은 ACK 쓰기 타임아웃 (기본 10초).
	WriteTimeout time.Duration

	// IdleTimeout 은 연결 유휴 한계 (기본 5분). 초과 시 연결 종료.
	IdleTimeout time.Duration

	// MaxConnections 는 동시 접속 상한 (0 = 무제한).
	MaxConnections int

	// Handler 는 메시지 처리 콜백 (필수).
	Handler MessageHandler

	// OnError 는 처리 중 에러 콜백 (선택). 운영 메트릭/로그용.
	OnError func(err error)
}

// MLLPServer 는 TCP 기반 HL7 MLLP 수신 서버.
type MLLPServer struct {
	cfg ServerConfig

	mu        sync.Mutex
	listener  net.Listener
	stopCh    chan struct{}
	doneCh    chan struct{}
	activeConn int32 // 현재 active 연결 (atomic)

	// 통계 — atomic 접근
	messagesAccepted int64
	messagesError    int64
	messagesRejected int64
}

// NewMLLPServer 생성. cfg.Handler 가 nil 이면 에러.
func NewMLLPServer(cfg ServerConfig) (*MLLPServer, error) {
	if cfg.Handler == nil {
		return nil, errors.New("Handler 필수")
	}
	if cfg.Addr == "" {
		cfg.Addr = "0.0.0.0:2575"
	}
	if cfg.ReadTimeout <= 0 {
		cfg.ReadTimeout = 30 * time.Second
	}
	if cfg.WriteTimeout <= 0 {
		cfg.WriteTimeout = 10 * time.Second
	}
	if cfg.IdleTimeout <= 0 {
		cfg.IdleTimeout = 5 * time.Minute
	}
	return &MLLPServer{cfg: cfg}, nil
}

// Start 는 서버를 시작. 비동기 — Stop() 으로 종료.
//
// 같은 인스턴스를 두 번 Start 하면 에러.
func (s *MLLPServer) Start() error {
	s.mu.Lock()
	if s.stopCh != nil {
		s.mu.Unlock()
		return errors.New("이미 시작됨")
	}
	listener, err := net.Listen("tcp", s.cfg.Addr)
	if err != nil {
		s.mu.Unlock()
		return fmt.Errorf("listen %s: %w", s.cfg.Addr, err)
	}
	s.listener = listener
	stopCh := make(chan struct{})
	doneCh := make(chan struct{})
	s.stopCh = stopCh
	s.doneCh = doneCh
	s.mu.Unlock()

	go s.acceptLoop(stopCh, doneCh)
	return nil
}

// Addr 는 실제 바인딩된 주소 (테스트에서 :0 사용 시 OS 가 할당한 포트 확인용).
func (s *MLLPServer) Addr() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.listener == nil {
		return ""
	}
	return s.listener.Addr().String()
}

// Stop 은 listener 종료 + accept 루프 종료 대기.
//
// 진행 중인 메시지 처리는 호출 측 ctx 로 cancellation 가능.
func (s *MLLPServer) Stop() error {
	s.mu.Lock()
	if s.stopCh == nil {
		s.mu.Unlock()
		return nil
	}
	close(s.stopCh)
	listener := s.listener
	doneCh := s.doneCh
	s.stopCh = nil
	s.listener = nil
	s.mu.Unlock()

	if listener != nil {
		_ = listener.Close()
	}
	if doneCh != nil {
		<-doneCh
	}
	return nil
}

// Stats 는 운영 통계 스냅샷.
type Stats struct {
	MessagesAccepted int64
	MessagesError    int64
	MessagesRejected int64
	ActiveConn       int32
}

// Stats 는 현재 통계 반환.
func (s *MLLPServer) Stats() Stats {
	return Stats{
		MessagesAccepted: atomic.LoadInt64(&s.messagesAccepted),
		MessagesError:    atomic.LoadInt64(&s.messagesError),
		MessagesRejected: atomic.LoadInt64(&s.messagesRejected),
		ActiveConn:       atomic.LoadInt32(&s.activeConn),
	}
}

func (s *MLLPServer) acceptLoop(stopCh chan struct{}, doneCh chan struct{}) {
	defer close(doneCh)
	for {
		select {
		case <-stopCh:
			return
		default:
		}
		conn, err := s.listener.Accept()
		if err != nil {
			// Stop 시 Accept 가 에러 반환 — 의도된 종료
			select {
			case <-stopCh:
				return
			default:
			}
			s.reportErr(fmt.Errorf("accept: %w", err))
			continue
		}
		// max connections 검사
		if s.cfg.MaxConnections > 0 &&
			atomic.LoadInt32(&s.activeConn) >= int32(s.cfg.MaxConnections) {
			_ = conn.Close()
			s.reportErr(errors.New("max connections 초과 — 연결 거부"))
			continue
		}
		atomic.AddInt32(&s.activeConn, 1)
		go s.handleConn(conn, stopCh)
	}
}

func (s *MLLPServer) handleConn(conn net.Conn, stopCh chan struct{}) {
	defer atomic.AddInt32(&s.activeConn, -1)
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		select {
		case <-stopCh:
			return
		default:
		}

		_ = conn.SetReadDeadline(time.Now().Add(s.cfg.IdleTimeout))
		raw, err := ReadFrame(reader)
		if err != nil {
			// 유휴 시간 초과 또는 클라이언트 종료 — 정상 종료
			return
		}

		// 처리 시간 측정용 ctx
		ctx, cancel := context.WithTimeout(context.Background(), s.cfg.ReadTimeout)
		msg, parseErr := hl7.Parse(string(raw))
		var ackCode, ackText string
		if parseErr != nil {
			ackCode = "AR" // reject (파싱 불가)
			ackText = "HL7 파싱 실패: " + parseErr.Error()
			atomic.AddInt64(&s.messagesRejected, 1)
		} else {
			code, text, herr := s.cfg.Handler(ctx, raw, msg)
			if code == "" {
				code = "AA"
			}
			ackCode = code
			ackText = text
			if herr != nil {
				s.reportErr(herr)
			}
			switch ackCode {
			case "AA":
				atomic.AddInt64(&s.messagesAccepted, 1)
			case "AE":
				atomic.AddInt64(&s.messagesError, 1)
			case "AR":
				atomic.AddInt64(&s.messagesRejected, 1)
			}
		}
		cancel()

		// ACK 생성 + 송신
		ack := BuildACK(msg, ackCode, ackText)
		_ = conn.SetWriteDeadline(time.Now().Add(s.cfg.WriteTimeout))
		if werr := WriteFrame(writer, []byte(ack)); werr != nil {
			s.reportErr(fmt.Errorf("ACK 송신: %w", werr))
			return
		}
	}
}

func (s *MLLPServer) reportErr(err error) {
	if s.cfg.OnError != nil {
		s.cfg.OnError(err)
	}
}
