import createClient from 'openapi-fetch';
import type { paths } from '../types/api';

// API base URL - 開発環境ではlocalhostのGoサーバーを指す
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

// OpenAPI クライアントを作成
export const api = createClient<paths>({
  baseUrl: API_BASE_URL,
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
