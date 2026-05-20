// Package service: IoT 프로토콜 핸들러 (MQTT/BLE/CoAP)
//
// 인메모리 프로토콜 핸들러를 제공하여 실제 브로커/스택과 통합되기 전까지
// 비즈니스 로직과 명령 라우팅을 검증할 수 있도록 합니다.
package service

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ============================================================================
// 프로토콜 핸들러 인터페이스
// ============================================================================

// ProtocolHandler는 IoT 프로토콜 통신 인터페이스입니다.
type ProtocolHandler interface {
	Protocol() string
	Send(ctx context.Context, deviceID string, payload []byte) error
	Subscribe(deviceID string, handler MessageHandler) error
	Unsubscribe(deviceID string) error
	Stats() *ProtocolStats
}

// MessageHandler는 디바이스에서 수신한 메시지를 처리합니다.
type MessageHandler func(deviceID string, payload []byte)

// ProtocolStats는 프로토콜 통계입니다.
type ProtocolStats struct {
	Protocol        string
	MessagesSent    int64
	MessagesReceived int64
	Errors          int64
	ActiveSessions  int
	LastActivityAt  time.Time
}

// ============================================================================
// MQTT 핸들러 (인메모리)
// ============================================================================

// MQTTHandler는 인메모리 MQTT 브로커 시뮬레이터입니다.
type MQTTHandler struct {
	mu          sync.RWMutex
	subscribers map[string]MessageHandler
	stats       ProtocolStats
}

// NewMQTTHandler는 새 MQTT 핸들러를 생성합니다.
func NewMQTTHandler() *MQTTHandler {
	return &MQTTHandler{
		subscribers: make(map[string]MessageHandler),
		stats:       ProtocolStats{Protocol: "mqtt"},
	}
}

// Protocol은 프로토콜 이름을 반환합니다.
func (h *MQTTHandler) Protocol() string { return "mqtt" }

// Send는 디바이스에 메시지를 전송합니다.
func (h *MQTTHandler) Send(ctx context.Context, deviceID string, payload []byte) error {
	if deviceID == "" {
		h.incrErrors()
		return errors.New("device_id required")
	}
	if len(payload) > 256*1024 {
		h.incrErrors()
		return errors.New("payload too large (>256KB)")
	}

	h.mu.Lock()
	h.stats.MessagesSent++
	h.stats.LastActivityAt = time.Now().UTC()
	h.mu.Unlock()
	return nil
}

// Subscribe는 디바이스 메시지 수신을 구독합니다.
func (h *MQTTHandler) Subscribe(deviceID string, handler MessageHandler) error {
	if deviceID == "" || handler == nil {
		return errors.New("device_id and handler required")
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.subscribers[deviceID] = handler
	h.stats.ActiveSessions = len(h.subscribers)
	return nil
}

// Unsubscribe는 구독을 해제합니다.
func (h *MQTTHandler) Unsubscribe(deviceID string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.subscribers, deviceID)
	h.stats.ActiveSessions = len(h.subscribers)
	return nil
}

// Stats는 프로토콜 통계를 반환합니다.
func (h *MQTTHandler) Stats() *ProtocolStats {
	h.mu.RLock()
	defer h.mu.RUnlock()
	stats := h.stats
	return &stats
}

// SimulateReceive는 시뮬레이션 메시지 수신입니다 (테스트 용도).
func (h *MQTTHandler) SimulateReceive(deviceID string, payload []byte) {
	h.mu.RLock()
	handler, ok := h.subscribers[deviceID]
	h.mu.RUnlock()

	h.mu.Lock()
	h.stats.MessagesReceived++
	h.stats.LastActivityAt = time.Now().UTC()
	h.mu.Unlock()

	if ok {
		handler(deviceID, payload)
	}
}

func (h *MQTTHandler) incrErrors() {
	h.mu.Lock()
	h.stats.Errors++
	h.mu.Unlock()
}

// ============================================================================
// BLE 핸들러 (재시도 로직 포함)
// ============================================================================

// BLEHandler는 BLE 디바이스 통신 핸들러입니다.
type BLEHandler struct {
	mu          sync.RWMutex
	connected   map[string]bool
	subscribers map[string]MessageHandler
	stats       ProtocolStats
	maxRetries  int
}

// NewBLEHandler는 새 BLE 핸들러를 생성합니다.
func NewBLEHandler() *BLEHandler {
	return &BLEHandler{
		connected:   make(map[string]bool),
		subscribers: make(map[string]MessageHandler),
		stats:       ProtocolStats{Protocol: "ble"},
		maxRetries:  3,
	}
}

// Protocol은 프로토콜 이름을 반환합니다.
func (h *BLEHandler) Protocol() string { return "ble" }

// Connect는 디바이스에 BLE 연결을 수립합니다.
func (h *BLEHandler) Connect(deviceID string) error {
	if deviceID == "" {
		return errors.New("device_id required")
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.connected[deviceID] = true
	return nil
}

// Disconnect는 BLE 연결을 종료합니다.
func (h *BLEHandler) Disconnect(deviceID string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.connected, deviceID)
	delete(h.subscribers, deviceID)
	return nil
}

// IsConnected는 디바이스의 연결 상태를 반환합니다.
func (h *BLEHandler) IsConnected(deviceID string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.connected[deviceID]
}

// Send는 BLE 디바이스에 메시지를 전송합니다.
//
// 페이로드 크기는 ATT MTU 한계(512 bytes)를 초과할 수 없습니다.
// 실제 운영에서는 외부 BLE 스택에서 전송 실패 시 transport 어댑터가
// `SendOnce` 훅을 통해 재시도하도록 설계되어 있어, 본 인메모리 핸들러는
// 검증 후 곧바로 카운터를 증가시킵니다.
func (h *BLEHandler) Send(ctx context.Context, deviceID string, payload []byte) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	h.mu.RLock()
	connected := h.connected[deviceID]
	h.mu.RUnlock()

	if !connected {
		h.incrErrors()
		return errors.New("device not connected")
	}

	// 페이로드 크기 검증 (BLE ATT MTU)
	if len(payload) > 512 {
		h.incrErrors()
		return errors.New("BLE payload exceeds 512 bytes")
	}

	h.mu.Lock()
	h.stats.MessagesSent++
	h.stats.LastActivityAt = time.Now().UTC()
	h.mu.Unlock()
	return nil
}

// incrErrors는 에러 카운터를 1 증가시킵니다.
func (h *BLEHandler) incrErrors() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.stats.Errors++
}

// SendWithRetry는 외부 transport 어댑터를 통해 exponential backoff 재시도를 수행합니다.
//
// transport: 실제 BLE 스택을 추상화한 함수 (실패 시 nil이 아닌 error 반환).
// maxRetries: 최대 재시도 횟수 (0이면 핸들러 기본값 사용).
func (h *BLEHandler) SendWithRetry(
	ctx context.Context,
	deviceID string,
	payload []byte,
	transport func(deviceID string, payload []byte) error,
	maxRetries int,
) error {
	if !h.IsConnected(deviceID) {
		return errors.New("device not connected")
	}
	if len(payload) > 512 {
		return errors.New("BLE payload exceeds 512 bytes")
	}
	if transport == nil {
		return errors.New("transport function required")
	}
	if maxRetries <= 0 {
		maxRetries = h.maxRetries
	}

	var lastErr error
	backoff := 10 * time.Millisecond
	for attempt := 0; attempt < maxRetries; attempt++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err := transport(deviceID, payload); err == nil {
			h.mu.Lock()
			h.stats.MessagesSent++
			h.stats.LastActivityAt = time.Now().UTC()
			h.mu.Unlock()
			return nil
		} else {
			lastErr = err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}
		backoff *= 2
	}

	h.incrErrors()
	return fmt.Errorf("BLE send failed after %d retries: %w", maxRetries, lastErr)
}

// Subscribe는 BLE 알림을 구독합니다.
func (h *BLEHandler) Subscribe(deviceID string, handler MessageHandler) error {
	if !h.IsConnected(deviceID) {
		return errors.New("device not connected")
	}
	h.mu.Lock()
	defer h.mu.Unlock()
	h.subscribers[deviceID] = handler
	h.stats.ActiveSessions = len(h.subscribers)
	return nil
}

// Unsubscribe는 BLE 알림 구독을 해제합니다.
func (h *BLEHandler) Unsubscribe(deviceID string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.subscribers, deviceID)
	h.stats.ActiveSessions = len(h.subscribers)
	return nil
}

// Stats는 통계를 반환합니다.
func (h *BLEHandler) Stats() *ProtocolStats {
	h.mu.RLock()
	defer h.mu.RUnlock()
	stats := h.stats
	return &stats
}

// ============================================================================
// CoAP 핸들러
// ============================================================================

// CoAPHandler는 CoAP 프로토콜 핸들러 (REST 매핑)입니다.
type CoAPHandler struct {
	mu       sync.RWMutex
	resources map[string][]byte // path -> payload (cache)
	stats    ProtocolStats
}

// NewCoAPHandler는 새 CoAP 핸들러를 생성합니다.
func NewCoAPHandler() *CoAPHandler {
	return &CoAPHandler{
		resources: make(map[string][]byte),
		stats:     ProtocolStats{Protocol: "coap"},
	}
}

// Protocol은 프로토콜 이름을 반환합니다.
func (h *CoAPHandler) Protocol() string { return "coap" }

// GET은 CoAP GET 요청을 처리합니다.
func (h *CoAPHandler) GET(deviceID, path string) ([]byte, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	key := deviceID + ":" + path
	data, ok := h.resources[key]
	if !ok {
		return nil, errors.New("4.04 Not Found")
	}
	return data, nil
}

// PUT은 CoAP PUT 요청을 처리합니다.
func (h *CoAPHandler) PUT(deviceID, path string, payload []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	key := deviceID + ":" + path
	h.resources[key] = payload
	h.stats.MessagesSent++
	h.stats.LastActivityAt = time.Now().UTC()
	return nil
}

// POST는 CoAP POST 요청을 처리합니다.
func (h *CoAPHandler) POST(deviceID, path string, payload []byte) error {
	return h.PUT(deviceID, path, payload)
}

// DELETE는 CoAP DELETE 요청을 처리합니다.
func (h *CoAPHandler) DELETE(deviceID, path string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	key := deviceID + ":" + path
	delete(h.resources, key)
	return nil
}

// Send는 CoAP를 통해 명령을 전송합니다 (POST 매핑).
func (h *CoAPHandler) Send(ctx context.Context, deviceID string, payload []byte) error {
	return h.POST(deviceID, "/cmd", payload)
}

// Subscribe는 CoAP Observe 구독입니다 (인메모리 시뮬레이션).
func (h *CoAPHandler) Subscribe(deviceID string, handler MessageHandler) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.stats.ActiveSessions++
	return nil
}

// Unsubscribe는 Observe 구독 해제입니다.
func (h *CoAPHandler) Unsubscribe(deviceID string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.stats.ActiveSessions > 0 {
		h.stats.ActiveSessions--
	}
	return nil
}

// Stats는 통계를 반환합니다.
func (h *CoAPHandler) Stats() *ProtocolStats {
	h.mu.RLock()
	defer h.mu.RUnlock()
	stats := h.stats
	return &stats
}

// ============================================================================
// 디바이스 명령 큐 (우선순위 + TTL)
// ============================================================================

// QueuedCommand는 큐에 들어간 디바이스 명령입니다.
type QueuedCommand struct {
	ID         string
	DeviceID   string
	Payload    []byte
	Priority   int       // 1(낮음) ~ 4(긴급)
	EnqueuedAt time.Time
	TTL        time.Duration
	Attempts   int
}

// IsExpired는 TTL 만료 여부를 반환합니다.
func (c *QueuedCommand) IsExpired() bool {
	if c.TTL == 0 {
		return false
	}
	return time.Since(c.EnqueuedAt) > c.TTL
}

// CommandQueue는 우선순위 + TTL 기반 명령 큐입니다.
type CommandQueue struct {
	mu       sync.Mutex
	commands []*QueuedCommand
}

// NewCommandQueue는 새 명령 큐를 생성합니다.
func NewCommandQueue() *CommandQueue {
	return &CommandQueue{commands: make([]*QueuedCommand, 0)}
}

// Enqueue는 명령을 큐에 추가합니다.
func (q *CommandQueue) Enqueue(cmd *QueuedCommand) {
	if cmd == nil {
		return
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	cmd.EnqueuedAt = time.Now().UTC()
	q.commands = append(q.commands, cmd)
}

// DequeueByPriority는 우선순위가 가장 높은 만료되지 않은 명령을 가져옵니다.
func (q *CommandQueue) DequeueByPriority() *QueuedCommand {
	q.mu.Lock()
	defer q.mu.Unlock()

	// 만료된 것 제거
	q.purgeExpiredLocked()

	if len(q.commands) == 0 {
		return nil
	}

	// 우선순위 + 입력 시각 정렬
	sort.SliceStable(q.commands, func(i, j int) bool {
		if q.commands[i].Priority != q.commands[j].Priority {
			return q.commands[i].Priority > q.commands[j].Priority
		}
		return q.commands[i].EnqueuedAt.Before(q.commands[j].EnqueuedAt)
	})

	cmd := q.commands[0]
	q.commands = q.commands[1:]
	return cmd
}

// PurgeExpired는 만료된 명령들을 제거합니다.
func (q *CommandQueue) PurgeExpired() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.purgeExpiredLocked()
}

func (q *CommandQueue) purgeExpiredLocked() int {
	purged := 0
	keep := q.commands[:0]
	for _, c := range q.commands {
		if c.IsExpired() {
			purged++
			continue
		}
		keep = append(keep, c)
	}
	q.commands = keep
	return purged
}

// Size는 큐의 크기를 반환합니다.
func (q *CommandQueue) Size() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.commands)
}

// SizeByDevice는 특정 디바이스의 대기 명령 수를 반환합니다.
func (q *CommandQueue) SizeByDevice(deviceID string) int {
	q.mu.Lock()
	defer q.mu.Unlock()
	count := 0
	for _, c := range q.commands {
		if c.DeviceID == deviceID && !c.IsExpired() {
			count++
		}
	}
	return count
}

// ============================================================================
// 연결 모니터 (heartbeat + 자동 재연결)
// ============================================================================

// ConnectionMonitor는 디바이스 연결 상태를 모니터링합니다.
type ConnectionMonitor struct {
	mu             sync.RWMutex
	lastSeen       map[string]time.Time
	heartbeatTimeout time.Duration
}

// NewConnectionMonitor는 새 연결 모니터를 생성합니다.
func NewConnectionMonitor(timeout time.Duration) *ConnectionMonitor {
	return &ConnectionMonitor{
		lastSeen:         make(map[string]time.Time),
		heartbeatTimeout: timeout,
	}
}

// RecordHeartbeat는 디바이스의 heartbeat를 기록합니다.
func (m *ConnectionMonitor) RecordHeartbeat(deviceID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lastSeen[deviceID] = time.Now().UTC()
}

// IsAlive는 디바이스가 살아있는지 확인합니다.
func (m *ConnectionMonitor) IsAlive(deviceID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	last, ok := m.lastSeen[deviceID]
	if !ok {
		return false
	}
	return time.Since(last) < m.heartbeatTimeout
}

// LastSeen은 디바이스의 마지막 heartbeat 시각을 반환합니다.
func (m *ConnectionMonitor) LastSeen(deviceID string) (time.Time, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	t, ok := m.lastSeen[deviceID]
	return t, ok
}

// FindDeadDevices는 timeout을 초과한 디바이스 목록을 반환합니다.
func (m *ConnectionMonitor) FindDeadDevices() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var dead []string
	now := time.Now().UTC()
	for id, last := range m.lastSeen {
		if now.Sub(last) >= m.heartbeatTimeout {
			dead = append(dead, id)
		}
	}
	sort.Strings(dead)
	return dead
}

// Cleanup은 dead 디바이스들을 제거합니다.
func (m *ConnectionMonitor) Cleanup() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	cleaned := 0
	now := time.Now().UTC()
	for id, last := range m.lastSeen {
		if now.Sub(last) >= m.heartbeatTimeout*2 {
			delete(m.lastSeen, id)
			cleaned++
		}
	}
	return cleaned
}

// ============================================================================
// 프로토콜 라우터 (멀티 프로토콜 통합)
// ============================================================================

// ProtocolRouter는 디바이스별 프로토콜에 따라 핸들러를 라우팅합니다.
type ProtocolRouter struct {
	mu       sync.RWMutex
	handlers map[string]ProtocolHandler // protocol -> handler
	devices  map[string]string           // deviceID -> protocol
}

// NewProtocolRouter는 새 프로토콜 라우터를 생성합니다.
func NewProtocolRouter() *ProtocolRouter {
	return &ProtocolRouter{
		handlers: make(map[string]ProtocolHandler),
		devices:  make(map[string]string),
	}
}

// Register는 프로토콜 핸들러를 등록합니다.
func (r *ProtocolRouter) Register(handler ProtocolHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.handlers[handler.Protocol()] = handler
}

// AssignDevice는 디바이스를 특정 프로토콜에 할당합니다.
func (r *ProtocolRouter) AssignDevice(deviceID, protocol string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.handlers[protocol]; !ok {
		return fmt.Errorf("protocol %q not registered", protocol)
	}
	r.devices[deviceID] = protocol
	return nil
}

// Send는 디바이스 프로토콜에 따라 적절한 핸들러로 라우팅합니다.
func (r *ProtocolRouter) Send(ctx context.Context, deviceID string, payload []byte) error {
	r.mu.RLock()
	protocol, ok := r.devices[deviceID]
	if !ok {
		r.mu.RUnlock()
		return fmt.Errorf("device %q not assigned to any protocol", deviceID)
	}
	handler := r.handlers[protocol]
	r.mu.RUnlock()
	return handler.Send(ctx, deviceID, payload)
}

// AggregateStats는 모든 핸들러의 통계를 합칩니다.
func (r *ProtocolRouter) AggregateStats() map[string]*ProtocolStats {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]*ProtocolStats)
	for proto, h := range r.handlers {
		out[proto] = h.Stats()
	}
	return out
}

// EnsureUUIDForCommand는 명령에 UUID가 없으면 생성합니다.
func EnsureUUIDForCommand(cmd *QueuedCommand) {
	if cmd != nil && cmd.ID == "" {
		cmd.ID = uuid.New().String()
	}
}
