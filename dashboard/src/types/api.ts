export interface RequestConfig {
  method: string;
  url: string;
  headers: Record<string, string>;
  body?: string;
  timeout?: number;
}

export interface ApiResponse {
  statusCode: number;
  headers: Record<string, string>;
  body: string;
  timestamp: number;
  duration: number;
}

export interface RequestHistory {
  id: string;
  request: RequestConfig;
  response?: ApiResponse;
  error?: string;
  timestamp: number;
}

export interface ProxyConfig {
  host: string;
  port: number;
  enabled: boolean;
}

export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH' | 'HEAD' | 'OPTIONS';

export interface HeaderPair {
  key: string;
  value: string;
  enabled: boolean;
}

export interface ApiError {
  body?: string;
  message?: string;
  statusCode?: number;
} 