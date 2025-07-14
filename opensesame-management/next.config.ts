import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      {
        source: "/management/:path*",
        destination: "http://localhost:11072/management/:path*",
      },
    ];
  },
};

export default nextConfig;
