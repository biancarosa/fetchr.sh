'use client';

import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { Badge } from './ui/badge';
import { Button } from './ui/button';
import { ScrollArea } from './ui/scroll-area';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from './ui/tooltip';
import { History, Trash2, ExternalLink, Server, RefreshCw } from 'lucide-react';
import { useRequestStore } from '../hooks/useRequestStore';
import { apiService, BackendRequestRecord } from '../services/api';
import { formatDistanceToNow } from 'date-fns';
import { HttpMethod } from '../types/api';

export function Sidebar() {
  const { 
    setMethod,
    setUrl,
    setHeaders,
    setBody,
  } = useRequestStore();

  const [backendHistory, setBackendHistory] = useState<BackendRequestRecord[]>([]);
  const [isLoadingBackend, setIsLoadingBackend] = useState(false);
  const [lastRefresh, setLastRefresh] = useState<Date | null>(null);

  const fetchBackendHistory = async () => {
    setIsLoadingBackend(true);
    try {
      const records = await apiService.getRequestHistory();
      setBackendHistory(records);
      setLastRefresh(new Date());
    } catch (error) {
      console.error('Failed to fetch backend history:', error);
    } finally {
      setIsLoadingBackend(false);
    }
  };

  const clearBackendHistory = async () => {
    try {
      const success = await apiService.clearRequestHistory();
      if (success) {
        setBackendHistory([]);
        setLastRefresh(new Date());
      }
    } catch (error) {
      console.error('Failed to clear backend history:', error);
    }
  };

  useEffect(() => {
    fetchBackendHistory();
  }, []);

  const loadBackendHistoryItem = (item: BackendRequestRecord) => {
    setMethod(item.method as HttpMethod);
    setUrl(item.url);
    setBody(item.request_body || '');
    
    // Convert headers object back to HeaderPair array
    const headerPairs = Object.entries(item.request_headers || {}).map(([key, value]) => ({
      key,
      value: value as string,
      enabled: true,
    }));
    
    // Add an empty header pair at the end
    headerPairs.push({ key: '', value: '', enabled: true });
    
    setHeaders(headerPairs);
  };

  const getMethodColor = (method: string) => {
    switch (method) {
      case 'GET': return 'bg-green-100 text-green-800';
      case 'POST': return 'bg-blue-100 text-blue-800';
      case 'PUT': return 'bg-orange-100 text-orange-800';
      case 'DELETE': return 'bg-red-100 text-red-800';
      case 'PATCH': return 'bg-purple-100 text-purple-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  };

  const getStatusColor = (statusCode?: number) => {
    if (!statusCode) return 'bg-gray-100 text-gray-800';
    if (statusCode >= 200 && statusCode < 300) return 'bg-green-100 text-green-800';
    if (statusCode >= 300 && statusCode < 400) return 'bg-blue-100 text-blue-800';
    if (statusCode >= 400 && statusCode < 500) return 'bg-yellow-100 text-yellow-800';
    if (statusCode >= 500) return 'bg-red-100 text-red-800';
    return 'bg-gray-100 text-gray-800';
  };

  const truncateUrl = (url: string, maxLength: number = 40) => {
    if (url.length <= maxLength) return url;
    return '...' + url.slice(-(maxLength - 3));
  };

  return (
    <div className="h-full bg-background">
      <Card className="h-full rounded-none border-0">
        <CardHeader className="pb-4">
          <CardTitle className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <History className="h-5 w-5" />
              <span>Proxy Requests</span>
            </div>
            <div className="flex items-center gap-1">
              <Badge variant="outline" className="text-xs">
                {backendHistory.length}
              </Badge>
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button 
                      variant="ghost" 
                      size="sm" 
                      onClick={fetchBackendHistory}
                      disabled={isLoadingBackend}
                    >
                      <RefreshCw className={`h-4 w-4 ${isLoadingBackend ? 'animate-spin' : ''}`} />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>Refresh proxy history</TooltipContent>
                </Tooltip>
              </TooltipProvider>
            </div>
          </CardTitle>
          {lastRefresh && (
            <p className="text-xs text-muted-foreground">
              Last updated: {formatDistanceToNow(lastRefresh, { addSuffix: true })}
            </p>
          )}
        </CardHeader>
        <CardContent className="p-0">
          <div className="px-2 pb-2">
            {backendHistory.length > 0 && (
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button 
                      variant="ghost" 
                      size="sm" 
                      onClick={clearBackendHistory}
                      className="w-full text-xs"
                    >
                      <Trash2 className="h-4 w-4 mr-1" />
                      Clear history
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>Clear all proxy request history</TooltipContent>
                </Tooltip>
              </TooltipProvider>
            )}
          </div>
          
          <ScrollArea className="h-[calc(100vh-10rem)]">
            {backendHistory.length === 0 ? (
              <div className="p-4 text-center text-muted-foreground">
                <Server className="h-8 w-8 mx-auto mb-2 opacity-50" />
                <p className="text-sm">No proxy requests yet</p>
                <p className="text-xs mt-1">Requests through the proxy will appear here</p>
              </div>
            ) : (
              <div className="space-y-2 p-2">
                {backendHistory.map((item) => (
                  <div
                    key={item.id}
                    className="border rounded-lg p-3 hover:bg-muted/50 cursor-pointer transition-colors"
                    onClick={() => loadBackendHistoryItem(item)}
                  >
                    <div className="flex items-center justify-between mb-2">
                      <Badge 
                        className={`text-xs ${getMethodColor(item.method)}`}
                        variant="secondary"
                      >
                        {item.method}
                      </Badge>
                      <div className="flex gap-1">
                        {item.success ? (
                          <Badge 
                            className={`text-xs ${getStatusColor(item.response_status)}`}
                            variant="secondary"
                          >
                            {item.response_status}
                          </Badge>
                        ) : (
                          <Badge className="text-xs bg-red-100 text-red-800" variant="secondary">
                            Error
                          </Badge>
                        )}
                      </div>
                    </div>
                    
                    <div className="space-y-1">
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <p className="text-sm font-medium truncate">
                              {truncateUrl(item.url)}
                            </p>
                          </TooltipTrigger>
                          <TooltipContent>
                            <p className="max-w-md break-all">{item.url}</p>
                          </TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                      
                      <div className="flex items-center justify-between text-xs text-muted-foreground">
                        <span>
                          {formatDistanceToNow(new Date(item.timestamp), { addSuffix: true })}
                        </span>
                        <div className="flex gap-2">
                          <TooltipProvider>
                            <Tooltip>
                              <TooltipTrigger asChild>
                                <span>{item.total_duration_ms}ms</span>
                              </TooltipTrigger>
                              <TooltipContent>
                                <div className="text-xs">
                                  <div>Total: {item.total_duration_ms}ms</div>
                                  <div>Upstream: {item.upstream_latency_ms}ms</div>
                                  <div>Proxy: {item.proxy_overhead_ms}ms</div>
                                </div>
                              </TooltipContent>
                            </Tooltip>
                          </TooltipProvider>
                        </div>
                      </div>
                      
                      {!item.success && item.error && (
                        <p className="text-xs text-red-600 truncate">
                          {item.error}
                        </p>
                      )}
                    </div>
                    
                    <div className="flex items-center justify-end mt-2">
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <Button
                              variant="ghost"
                              size="sm"
                              className="h-6 w-6 p-0"
                              onClick={(e) => {
                                e.stopPropagation();
                                window.open(item.url, '_blank');
                              }}
                            >
                              <ExternalLink className="h-3 w-3" />
                            </Button>
                          </TooltipTrigger>
                          <TooltipContent>Open in new tab</TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </ScrollArea>
        </CardContent>
      </Card>
    </div>
  );
} 