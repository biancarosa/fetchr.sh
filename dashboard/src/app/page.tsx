'use client';

import { Header } from '../components/Header';
import { Sidebar } from '../components/Sidebar';
import { RequestBuilder } from '../components/RequestBuilder';
import { RequestStats } from '../components/RequestStats';

export default function Home() {
  return (
    <div className="h-screen flex flex-col">
      <Header />
      <div className="flex flex-1 overflow-hidden">
        <div className="flex flex-col w-80 border-r bg-background">
          <div className="flex-1 min-h-0">
            <Sidebar />
          </div>
          <div className="border-t">
            <RequestStats className="rounded-none border-0" />
          </div>
        </div>
        <main className="flex-1 p-6 overflow-auto">
          <RequestBuilder />
        </main>
      </div>
    </div>
  );
}
