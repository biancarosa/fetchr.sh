'use client';

import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from './ui/card';
import { Badge } from './ui/badge';
import { Button } from './ui/button';
import { ScrollArea } from './ui/scroll-area';
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from './ui/tooltip';
import { History, Trash2, Clock, ExternalLink } from 'lucide-react';
import { useRequestStore } from '../hooks/useRequestStore';
import { formatDistanceToNow } from 'date-fns';
import { RequestHistory, HttpMethod } from '../types/api';

export function Sidebar() {
  const { 
    history, 
    clearHistory,
    setMethod,
    setUrl,
    setHeaders,
    setBody,
  } = useRequestStore();

  const loadHistoryItem = (item: RequestHistory) => {
    setMethod(item.request.method as HttpMethod);
    setUrl(item.request.url);
    setBody(item.request.body || '');
    
    // Convert headers object back to HeaderPair array
    const headerPairs = Object.entries(item.request.headers || {}).map(([key, value]) => ({
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
    <div className="w-80 h-full border-r bg-background">
      <Card className="h-full rounded-none border-0">
        <CardHeader className="pb-4">
          <CardTitle className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <History className="h-5 w-5" />
              <span>Request History</span>
            </div>
            {history.length > 0 && (
              <TooltipProvider>
                <Tooltip>
                  <TooltipTrigger asChild>
                    <Button variant="ghost" size="sm" onClick={clearHistory}>
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </TooltipTrigger>
                  <TooltipContent>Clear history</TooltipContent>
                </Tooltip>
              </TooltipProvider>
            )}
          </CardTitle>
        </CardHeader>
        <CardContent className="p-0">
          <ScrollArea className="h-[calc(100vh-8rem)]">
            {history.length === 0 ? (
              <div className="p-4 text-center text-muted-foreground">
                <Clock className="h-8 w-8 mx-auto mb-2 opacity-50" />
                <p className="text-sm">No requests yet</p>
                <p className="text-xs mt-1">Your request history will appear here</p>
              </div>
            ) : (
              <div className="space-y-2 p-2">
                {history.map((item) => (
                  <div
                    key={item.id}
                    className="border rounded-lg p-3 hover:bg-muted/50 cursor-pointer transition-colors"
                    onClick={() => loadHistoryItem(item)}
                  >
                    <div className="flex items-center justify-between mb-2">
                      <Badge 
                        className={`text-xs ${getMethodColor(item.request.method)}`}
                        variant="secondary"
                      >
                        {item.request.method}
                      </Badge>
                      {item.response ? (
                        <Badge 
                          className={`text-xs ${getStatusColor(item.response.statusCode)}`}
                          variant="secondary"
                        >
                          {item.response.statusCode}
                        </Badge>
                      ) : item.error ? (
                        <Badge className="text-xs bg-red-100 text-red-800" variant="secondary">
                          Error
                        </Badge>
                      ) : null}
                    </div>
                    
                    <div className="space-y-1">
                      <TooltipProvider>
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <p className="text-sm font-medium truncate">
                              {truncateUrl(item.request.url)}
                            </p>
                          </TooltipTrigger>
                          <TooltipContent>
                            <p className="max-w-md break-all">{item.request.url}</p>
                          </TooltipContent>
                        </Tooltip>
                      </TooltipProvider>
                      
                      <div className="flex items-center justify-between text-xs text-muted-foreground">
                        <span>
                          {formatDistanceToNow(new Date(item.timestamp), { addSuffix: true })}
                        </span>
                        {item.response && (
                          <span>{item.response.duration}ms</span>
                        )}
                      </div>
                      
                      {item.error && (
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
                                window.open(item.request.url, '_blank');
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