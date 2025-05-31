'use client';

import { Header } from '../../components/Header';
import { ProxyRequestsTable } from '../../components/ProxyRequestsTable';
import { RefreshProvider } from '../../hooks/useRefreshContext';

export default function RequestsHistoryPage() {
  return (
    <RefreshProvider>
      <div className="h-screen flex flex-col">
        <Header />
        <main className="flex-1 p-6 overflow-auto">
          <ProxyRequestsTable />
        </main>
      </div>
    </RefreshProvider>
  );
} 