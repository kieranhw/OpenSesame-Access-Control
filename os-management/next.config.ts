import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      {
        source: "/admin/:path*",
        destination: "http://localhost:11072/admin/:path*",
      },
    ];
  },
};

export default nextConfig;
