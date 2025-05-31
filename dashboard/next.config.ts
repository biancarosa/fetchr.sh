import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // Enable static export for embedding
  output: 'export',
  trailingSlash: true,
  distDir: 'out',
  
  // Disable all caching to ensure fresh data
  experimental: {
    staleTimes: {
      dynamic: 0,
      static: 0,
    },
  },
  
  // Add headers to prevent caching
  async headers() {
    return [
      {
        source: '/:path*',
        headers: [
          {
            key: 'Cache-Control',
            value: 'no-cache, no-store, must-revalidate',
          },
          {
            key: 'Pragma',
            value: 'no-cache',
          },
          {
            key: 'Expires',
            value: '0',
          },
        ],
      },
    ];
  },
  
  // Configure static export
  images: {
    unoptimized: true,
  },
};

export default nextConfig;
