import { apiClient } from './client';
import type { LoginRequest, RegisterRequest, AuthResponse, User } from '../types/auth';

export const authApi = {
  async login(data: LoginRequest): Promise<AuthResponse> {
    return apiClient.request<AuthResponse>('/api/auth/login', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  async register(data: RegisterRequest): Promise<AuthResponse> {
    return apiClient.request<AuthResponse>('/api/auth/register', {
      method: 'POST',
      body: JSON.stringify(data),
    });
  },

  async getCurrentUser(token: string): Promise<User> {
    return apiClient.request<User>('/api/auth/me', {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  },

  async logout(token: string): Promise<void> {
    return apiClient.request<void>('/api/auth/logout', {
      method: 'POST',
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });
  },
};
