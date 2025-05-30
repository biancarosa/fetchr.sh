'use client';

import { Header } from '../components/Header';
import { Sidebar } from '../components/Sidebar';
import { RequestBuilder } from '../components/RequestBuilder';

export default function Home() {
  return (
    <div className="h-screen flex flex-col">
      <Header />
      <div className="flex flex-1 overflow-hidden">
        <Sidebar />
        <main className="flex-1 p-6 overflow-auto">
          <RequestBuilder />
        </main>
      </div>
    </div>
  );
}
