'use client';

import React, { Suspense, useRef } from 'react';
import { Canvas, useFrame } from '@react-three/fiber';
import { Float, MeshDistortMaterial, Sphere, OrbitControls } from '@react-three/drei';
import type { Mesh } from 'three';

function MolecularCore() {
  const ref = useRef<Mesh>(null);
  useFrame((_, delta) => {
    if (ref.current) {
      ref.current.rotation.y += delta * 0.15;
      ref.current.rotation.x += delta * 0.05;
    }
  });
  return (
    <Float speed={1.2} rotationIntensity={0.4} floatIntensity={0.6}>
      <Sphere ref={ref} args={[1.4, 96, 96]}>
        <MeshDistortMaterial
          color="#0284c7"
          emissive="#38bdf8"
          emissiveIntensity={0.18}
          roughness={0.18}
          metalness={0.35}
          distort={0.32}
          speed={1.4}
        />
      </Sphere>
    </Float>
  );
}

function Orbit({ radius, speed, color, size }: { radius: number; speed: number; color: string; size: number }) {
  const ref = useRef<Mesh>(null);
  useFrame((state) => {
    if (ref.current) {
      const t = state.clock.getElapsedTime() * speed;
      ref.current.position.x = Math.cos(t) * radius;
      ref.current.position.z = Math.sin(t) * radius;
      ref.current.position.y = Math.sin(t * 0.7) * 0.4;
    }
  });
  return (
    <mesh ref={ref}>
      <sphereGeometry args={[size, 32, 32]} />
      <meshStandardMaterial color={color} emissive={color} emissiveIntensity={0.6} />
    </mesh>
  );
}

export function MedicalHero3D() {
  return (
    // @ts-expect-error React Three Fiber RC type mismatch with React 19
    <Canvas camera={{ position: [0, 0, 4.4], fov: 45 }} dpr={[1, 2]}>
      <ambientLight intensity={0.7} />
      <directionalLight position={[3, 4, 5]} intensity={1.2} />
      <pointLight position={[-3, -2, -2]} intensity={0.6} color="#7dd3fc" />
      <Suspense fallback={null}>
        <MolecularCore />
        <Orbit radius={2.2} speed={0.6} color="#0ea5e9" size={0.09} />
        <Orbit radius={2.55} speed={-0.4} color="#3b82f6" size={0.07} />
        <Orbit radius={1.95} speed={0.9} color="#22d3ee" size={0.06} />
      </Suspense>
      <OrbitControls
        enableZoom={false}
        enablePan={false}
        autoRotate
        autoRotateSpeed={0.6}
        rotateSpeed={0.4}
      />
    </Canvas>
  );
}
