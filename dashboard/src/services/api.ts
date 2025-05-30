import { RequestConfig, ApiResponse, RequestHistory } from '../types/api';

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

  // Future endpoint for request history
  async getRequestHistory(): Promise<RequestHistory[]> {
    // This will be implemented when backend supports it
    return [];
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