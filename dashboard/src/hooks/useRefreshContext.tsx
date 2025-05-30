'use client';

import React, { createContext, useContext, useCallback } from 'react';

interface RefreshContextType {
  triggerRefresh: () => void;
  onRefresh: (callback: () => void) => () => void;
}

const RefreshContext = createContext<RefreshContextType | undefined>(undefined);

export function RefreshProvider({ children }: { children: React.ReactNode }) {
  const callbacks = React.useRef<Set<() => void>>(new Set());

  const triggerRefresh = useCallback(() => {
    callbacks.current.forEach(callback => callback());
  }, []);

  const onRefresh = useCallback((callback: () => void) => {
    callbacks.current.add(callback);
    
    // Return cleanup function
    return () => {
      callbacks.current.delete(callback);
    };
  }, []);

  return (
    <RefreshContext.Provider value={{ triggerRefresh, onRefresh }}>
      {children}
    </RefreshContext.Provider>
  );
}

export function useRefresh() {
  const context = useContext(RefreshContext);
  if (context === undefined) {
    throw new Error('useRefresh must be used within a RefreshProvider');
  }
  return context;
} 