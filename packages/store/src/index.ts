/**
 * @mmup/store — Zustand 전역 상태 관리
 * 
 * 사용자 세션, 테마, 알림 등 클라이언트 전역 상태
 */
import { create } from 'zustand';
import { persist } from 'zustand/middleware';

// ===== 인증 상태 =====
export interface AuthState {
  user: {
    id: string;
    email: string;
    name: string;
    persona: string;
    role: string;
  } | null;
  isAuthenticated: boolean;
  setUser: (user: AuthState['user']) => void;
  clearUser: () => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      user: null,
      isAuthenticated: false,
      setUser: (user) => set({ user, isAuthenticated: !!user }),
      clearUser: () => set({ user: null, isAuthenticated: false }),
    }),
    { name: 'mmup-auth' }
  )
);

// ===== 테마 상태 =====
export interface ThemeState {
  mode: 'dark' | 'light' | 'system';
  setMode: (mode: ThemeState['mode']) => void;
}

export const useThemeStore = create<ThemeState>()(
  persist(
    (set) => ({
      mode: 'light',
      setMode: (mode) => set({ mode }),
    }),
    { name: 'mmup-theme' }
  )
);

// ===== 알림 상태 =====
export interface Notification {
  id: string;
  type: 'info' | 'success' | 'warning' | 'error';
  title: string;
  message: string;
  read: boolean;
  createdAt: string;
}

export interface NotificationState {
  notifications: Notification[];
  unreadCount: number;
  addNotification: (n: Omit<Notification, 'id' | 'read' | 'createdAt'>) => void;
  markAsRead: (id: string) => void;
  markAllAsRead: () => void;
  clearAll: () => void;
}

export const useNotificationStore = create<NotificationState>()((set, get) => ({
  notifications: [],
  unreadCount: 0,
  addNotification: (n) => {
    const notification: Notification = {
      ...n,
      id: `notif_${Date.now()}`,
      read: false,
      createdAt: new Date().toISOString(),
    };
    set((s) => ({
      notifications: [notification, ...s.notifications].slice(0, 50),
      unreadCount: s.unreadCount + 1,
    }));
  },
  markAsRead: (id) =>
    set((s) => ({
      notifications: s.notifications.map((n) => (n.id === id ? { ...n, read: true } : n)),
      unreadCount: Math.max(0, s.unreadCount - 1),
    })),
  markAllAsRead: () =>
    set((s) => ({
      notifications: s.notifications.map((n) => ({ ...n, read: true })),
      unreadCount: 0,
    })),
  clearAll: () => set({ notifications: [], unreadCount: 0 }),
}));

// ===== 현재 환자 컨텍스트 =====
export interface PatientContextState {
  selectedPatientId: string | null;
  selectPatient: (id: string | null) => void;
}

export const usePatientContextStore = create<PatientContextState>()((set) => ({
  selectedPatientId: null,
  selectPatient: (id) => set({ selectedPatientId: id }),
}));
