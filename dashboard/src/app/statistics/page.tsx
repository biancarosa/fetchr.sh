'use client';

import { Header } from '../../components/Header';
import { StatisticsOverview } from '../../components/StatisticsOverview';
import { RefreshProvider } from '../../hooks/useRefreshContext';

export default function StatisticsPage() {
  return (
    <RefreshProvider>
      <div className="h-screen flex flex-col">
        <Header />
        <main className="flex-1 p-6 overflow-auto">
          <StatisticsOverview />
        </main>
      </div>
    </RefreshProvider>
  );
} 