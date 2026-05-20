/**
 * @mmup-axis 3 종합 데이터 분석
 * @mmup-stage 1 측정
 * @family C
 * @trinity IP3
 * @sb SB-1
 * @standard IEC 62304 Class C
 */
'use client';

import React, { useEffect, useState, useRef } from 'react';
import { AlertTriangle, CheckCircle, XCircle, Thermometer, Play, Pause, RotateCcw, Zap } from 'lucide-react';
import { PacketGenerator, DataPacket } from '../../lib/packetGenerator';

const PACKET_INTERVAL = 2000;
const TEMPERATURE_CRITICAL = 85;
const TEMPERATURE_WARNING = 75;

export function DataIntegrityConsole() {
  const [packets, setPackets] = useState<DataPacket[]>([]);
  const [isPaused, setIsPaused] = useState(false);
  const [systemTemp, setSystemTemp] = useState(45);
  const [isEmergencyStop, setIsEmergencyStop] = useState(false);
  const [mounted, setMounted] = useState(false);
  const [stats, setStats] = useState({
    total: 0,
    valid: 0,
    invalid: 0,
  });

  const generatorRef = useRef(new PacketGenerator());
  const intervalRef = useRef<number | null>(null);
  const previousPacketRef = useRef<DataPacket | null>(null);

  useEffect(() => {
    setMounted(true);
    const tempInterval = window.setInterval(() => {
      setSystemTemp((prev) => {
        const change = (Math.random() - 0.5) * 3;
        const newTemp = Math.max(40, Math.min(95, prev + change));

        if (newTemp >= TEMPERATURE_CRITICAL && !isEmergencyStop) {
          setIsEmergencyStop(true);
          setIsPaused(true);
        }

        return parseFloat(newTemp.toFixed(1));
      });
    }, 1000);

    return () => clearInterval(tempInterval);
  }, [isEmergencyStop]);

  useEffect(() => {
    if (!isPaused && !isEmergencyStop) {
      intervalRef.current = window.setInterval(() => {
        generateAndAddPacket(false);
      }, PACKET_INTERVAL);
    }

    return () => {
      if (intervalRef.current) {
        clearInterval(intervalRef.current);
      }
    };
  }, [isPaused, isEmergencyStop]);

  const generateAndAddPacket = (injectError: boolean = false) => {
    const newPacket = generatorRef.current.generatePacket(injectError);

    const validation = generatorRef.current.validatePacket(
      newPacket,
      previousPacketRef.current
    );

    newPacket.isValid = validation.isValid;
    if (!validation.isValid) {
      newPacket.validationError = validation.errors.join('; ');
    }

    previousPacketRef.current = newPacket;

    setPackets((prev) => {
      const updated = [newPacket, ...prev].slice(0, 20);
      return updated;
    });

    setStats((prev) => ({
      total: prev.total + 1,
      valid: validation.isValid ? prev.valid + 1 : prev.valid,
      invalid: validation.isValid ? prev.invalid : prev.invalid + 1,
    }));
  };

  const handleInjectError = () => {
    if (isEmergencyStop) return;
    generateAndAddPacket(true);
  };

  const handleReset = () => {
    generatorRef.current.reset();
    previousPacketRef.current = null;
    setPackets([]);
    setStats({ total: 0, valid: 0, invalid: 0 });
    setSystemTemp(45);
    setIsEmergencyStop(false);
    setIsPaused(false);
  };

  const handleEmergencyReset = () => {
    setSystemTemp(45);
    setIsEmergencyStop(false);
    setIsPaused(false);
  };

  const getTempColor = () => {
    if (systemTemp >= TEMPERATURE_CRITICAL) return 'text-red-500';
    if (systemTemp >= TEMPERATURE_WARNING) return 'text-yellow-500';
    return 'text-green-500';
  };

  const getTempBgColor = () => {
    if (systemTemp >= TEMPERATURE_CRITICAL) return 'bg-red-500/20 border-red-500';
    if (systemTemp >= TEMPERATURE_WARNING) return 'bg-yellow-500/20 border-yellow-500';
    return 'bg-green-500/20 border-green-500';
  };

  if (!mounted) {
    return <div className="w-full h-[600px] bg-gray-900 rounded-2xl border border-gray-700 animate-pulse flex items-center justify-center text-gray-500">시스템 부팅 중...</div>;
  }

  return (
    <div className="w-full h-[600px] flex flex-col gap-4 p-6 bg-gradient-to-br from-gray-900 to-gray-800 rounded-2xl border border-cyan-500/30 shadow-2xl">
      {isEmergencyStop && (
        <div className="absolute inset-0 z-50 flex items-center justify-center bg-black/80 backdrop-blur-sm rounded-2xl">
          <div className="bg-red-900/90 border-4 border-red-500 rounded-2xl p-8 max-w-md shadow-2xl animate-pulse">
            <div className="flex items-center gap-4 mb-6">
              <AlertTriangle className="w-16 h-16 text-red-400" />
              <div>
                <h2 className="text-3xl font-bold text-red-400">시스템 과부하</h2>
                <p className="text-red-300 text-lg mt-2">비상 정지</p>
              </div>
            </div>

            <div className="bg-black/40 rounded-lg p-4 mb-6">
              <p className="text-white text-lg mb-2">
                <span className="font-semibold">현재 온도:</span>{' '}
                <span className="text-red-400 text-2xl font-bold">{systemTemp}°C</span>
              </p>
              <p className="text-gray-300 text-sm">
                임계값 {TEMPERATURE_CRITICAL}°C 초과로 모든 작업이 중단되었습니다.
              </p>
            </div>

            <button
              onClick={handleEmergencyReset}
              className="w-full bg-red-600 hover:bg-red-500 text-white font-bold py-3 px-6 rounded-lg transition-all flex items-center justify-center gap-2"
            >
              <RotateCcw className="w-5 h-5" />
              시스템 냉각 및 재시작
            </button>
          </div>
        </div>
      )}

      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-cyan-400 mb-1">
            데이터 무결성 검증 콘솔
          </h2>
          <p className="text-gray-400 text-sm">
            실시간 패킷 검증 및 해시 체인 무결성 확인
          </p>
        </div>

        <div className="flex items-center gap-3">
          <div className={`flex items-center gap-2 px-4 py-2 rounded-lg border ${getTempBgColor()}`}>
            <Thermometer className={`w-5 h-5 ${getTempColor()}`} />
            <span className={`font-mono text-lg font-bold ${getTempColor()}`}>
              {systemTemp}°C
            </span>
          </div>

          <button
            onClick={() => setIsPaused(!isPaused)}
            disabled={isEmergencyStop}
            className={`px-4 py-2 rounded-lg font-semibold transition-all flex items-center gap-2 ${
              isEmergencyStop
                ? 'bg-gray-700 text-gray-500 cursor-not-allowed'
                : isPaused
                ? 'bg-green-600 hover:bg-green-500 text-white'
                : 'bg-yellow-600 hover:bg-yellow-500 text-white'
            }`}
          >
            {isPaused ? <Play className="w-4 h-4" /> : <Pause className="w-4 h-4" />}
            {isPaused ? '재개' : '일시정지'}
          </button>

          <button
            onClick={handleInjectError}
            disabled={isEmergencyStop}
            className={`px-4 py-2 rounded-lg font-semibold transition-all flex items-center gap-2 ${
              isEmergencyStop
                ? 'bg-gray-700 text-gray-500 cursor-not-allowed'
                : 'bg-red-600 hover:bg-red-500 text-white'
            }`}
          >
            <Zap className="w-4 h-4" />
            오류 주입
          </button>

          <button
            onClick={handleReset}
            className="px-4 py-2 rounded-lg font-semibold bg-gray-700 hover:bg-gray-600 text-white transition-all flex items-center gap-2"
          >
            <RotateCcw className="w-4 h-4" />
            초기화
          </button>
        </div>
      </div>

      <div className="grid grid-cols-4 gap-4">
        <div className="bg-gray-800/50 rounded-lg p-4 border border-cyan-500/30">
          <p className="text-gray-400 text-sm mb-1">총 패킷</p>
          <p className="text-2xl font-bold text-cyan-400">{stats.total}</p>
        </div>
        <div className="bg-gray-800/50 rounded-lg p-4 border border-green-500/30">
          <p className="text-gray-400 text-sm mb-1">검증 성공</p>
          <p className="text-2xl font-bold text-green-400">{stats.valid}</p>
        </div>
        <div className="bg-gray-800/50 rounded-lg p-4 border border-red-500/30">
          <p className="text-gray-400 text-sm mb-1">검증 실패</p>
          <p className="text-2xl font-bold text-red-400">{stats.invalid}</p>
        </div>
        <div className="bg-gray-800/50 rounded-lg p-4 border border-yellow-500/30">
          <p className="text-gray-400 text-sm mb-1">성공률</p>
          <p className="text-2xl font-bold text-yellow-400">
            {stats.total > 0 ? ((stats.valid / stats.total) * 100).toFixed(1) : 0}%
          </p>
        </div>
      </div>

      <div className="flex-1 bg-gray-900/50 rounded-lg border border-gray-700 overflow-hidden flex flex-col">
        <div className="bg-gray-800/80 px-4 py-2 border-b border-gray-700">
          <h3 className="text-cyan-400 font-semibold">실시간 패킷 로그</h3>
        </div>

        <div className="flex-1 overflow-y-auto p-4 space-y-3">
          {packets.length === 0 ? (
            <div className="h-full flex items-center justify-center">
              <p className="text-gray-500">패킷 수신 대기 중...</p>
            </div>
          ) : (
            packets.map((packet) => (
              <div
                key={packet.id}
                className={`bg-gray-800/50 rounded-lg border p-4 ${
                  packet.isValid
                    ? 'border-green-500/50 hover:border-green-500'
                    : 'border-red-500/50 hover:border-red-500'
                } transition-all`}
              >
                <div className="flex items-start justify-between mb-3">
                  <div className="flex items-center gap-3">
                    {packet.isValid ? (
                      <CheckCircle className="w-6 h-6 text-green-400 flex-shrink-0" />
                    ) : (
                      <XCircle className="w-6 h-6 text-red-400 flex-shrink-0" />
                    )}
                    <div>
                      <p className="text-white font-semibold">
                        패킷 #{packet.id}
                        {packet.isValid ? (
                          <span className="ml-2 text-green-400 text-sm">✓ 검증 성공</span>
                        ) : (
                          <span className="ml-2 text-red-400 text-sm">⚠ 무결성 검증 실패</span>
                        )}
                      </p>
                      <p className="text-gray-400 text-xs">{packet.timestamp}</p>
                    </div>
                  </div>
                </div>

                {!packet.isValid && packet.validationError && (
                  <div className="bg-red-900/30 border border-red-500/50 rounded p-2 mb-3">
                    <p className="text-red-400 text-sm font-semibold">
                      {packet.validationError}
                    </p>
                  </div>
                )}

                <div className="bg-black/40 rounded-lg p-3 font-mono text-xs overflow-x-auto">
                  <pre className="text-gray-300">
                    {JSON.stringify(
                      {
                        Header: packet.header,
                        Payload: packet.payload,
                        Checksum: packet.checksum,
                        PrevHash: packet.prevHash,
                        CurrentHash: packet.currentHash,
                      },
                      null,
                      2
                    )}
                  </pre>
                </div>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}
