import { RequestConfig, ApiResponse } from '../types/api';

// Backend request record from the Go API
export interface BackendRequestRecord {
  id: string;
  timestamp: string;
  method: string;
  url: string;
  request_headers: Record<string, string>;
  request_body?: string;
  response_status: number;
  response_headers: Record<string, string>;
  response_body?: string;
  proxy_start_time: string;
  upstream_start_time: string;
  upstream_end_time: string;
  proxy_end_time: string;
  proxy_overhead_ms: number;
  upstream_latency_ms: number;
  total_duration_ms: number;
  request_size: number;
  response_size: number;
  success: boolean;
  error?: string;
}

export interface BackendHistoryResponse {
  records: BackendRequestRecord[];
  total: number;
}

export interface RequestStats {
  total_requests: number;
  success_count: number;
  error_count: number;
  avg_duration_ms: number;
  avg_upstream_latency_ms: number;
  avg_proxy_overhead_ms: number;
  total_request_size: number;
  total_response_size: number;
  status_codes: Record<number, number>;
  methods: Record<string, number>;
}

class ApiService {
  private baseUrl: string;
  private proxyHost: string;
  private proxyPort: number;
  private adminPort: number;

  constructor() {
    this.baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3001';
    this.proxyHost = process.env.NEXT_PUBLIC_PROXY_HOST || 'localhost';
    this.proxyPort = parseInt(process.env.NEXT_PUBLIC_PROXY_PORT || '8080');
    this.adminPort = parseInt(process.env.NEXT_PUBLIC_ADMIN_PORT || '8081');
  }

  async makeRequest(config: RequestConfig): Promise<ApiResponse> {
    const startTime = Date.now();
    
    try {
      // Convert headers object to fetch headers
      const headers = new Headers();
      Object.entries(config.headers).forEach(([key, value]) => {
        if (key && value) {
          headers.append(key, value);
        }
      });

      // Add content type for POST/PUT/PATCH requests with body
      if (config.body && !headers.has('Content-Type')) {
        try {
          JSON.parse(config.body);
          headers.set('Content-Type', 'application/json');
        } catch {
          headers.set('Content-Type', 'text/plain');
        }
      }

      const fetchOptions: RequestInit = {
        method: config.method,
        headers,
        body: config.body || undefined,
      };

      // For now, we'll make direct requests through the browser
      // In the future, this could be routed through a backend API that uses the proxy
      const response = await fetch(config.url, fetchOptions);
      
      const responseBody = await response.text();
      const endTime = Date.now();

      // Convert response headers to object
      const responseHeaders: Record<string, string> = {};
      response.headers.forEach((value, key) => {
        responseHeaders[key] = value;
      });

      return {
        statusCode: response.status,
        headers: responseHeaders,
        body: responseBody,
        timestamp: startTime,
        duration: endTime - startTime,
      };
    } catch (error: unknown) {
      const endTime = Date.now();
      throw {
        statusCode: 0,
        headers: {},
        body: error instanceof Error ? error.message : 'Unknown error',
        timestamp: startTime,
        duration: endTime - startTime,
      };
    }
  }

  // Get request history from the backend
  async getRequestHistory(): Promise<BackendRequestRecord[]> {
    try {
      const response = await fetch(`http://${this.proxyHost}:${this.adminPort}/requests`);
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }
      const data: BackendHistoryResponse = await response.json();
      return data.records;
    } catch (error) {
      console.error('Failed to fetch request history:', error);
      return [];
    }
  }

  // Get request statistics from the backend
  async getRequestStats(): Promise<RequestStats | null> {
    try {
      const response = await fetch(`http://${this.proxyHost}:${this.adminPort}/requests/stats`);
      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }
      return await response.json();
    } catch (error) {
      console.error('Failed to fetch request stats:', error);
      return null;
    }
  }

  // Clear request history on the backend
  async clearRequestHistory(): Promise<boolean> {
    try {
      const response = await fetch(`http://${this.proxyHost}:${this.adminPort}/requests/clear`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
      });
      return response.ok;
    } catch (error) {
      console.error('Failed to clear request history:', error);
      return false;
    }
  }

  // Health check for proxy (now on admin port)
  async checkProxyHealth(): Promise<boolean> {
    try {
      const response = await fetch(`http://${this.proxyHost}:${this.adminPort}/healthz`);
      return response.ok;
    } catch {
      return false;
    }
  }
}

export const apiService = new ApiService(); 