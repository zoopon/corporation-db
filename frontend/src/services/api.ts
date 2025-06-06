import createClient from 'openapi-fetch';
import type { paths } from '../types/api';

// API base URL - nginx経由で/apiプレフィックスを使用
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

// nginx経由の場合は/apiプレフィックスを追加
const getApiBaseUrl = () => {
  if (API_BASE_URL === 'http://localhost') {
    return `${API_BASE_URL}/api`;
  }
  return API_BASE_URL;
};

// OpenAPI クライアントを作成
export const api = createClient<paths>({
  baseUrl: getApiBaseUrl(),
  headers: {
    'Content-Type': 'application/json',
  },
});

// APIエラーハンドリング用の型
export interface ApiError {
  message: string;
  status?: number;
}

// APIレスポンスのラッパー型
export type ApiResponse<T> = {
  data: T | null;
  error: ApiError | null;
  loading: boolean;
};
