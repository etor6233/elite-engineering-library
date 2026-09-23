import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  output: "standalone",
  outputFileTracingIncludes: {"/api/document-preview/assets/*": ["./third_party/pdfjs/*.mjs"]},
  generateBuildId: async () => {
    const id = process.env.ELITE_SOURCE_SHA256;
    if (!id) return null;
    if (!/^[0-9a-f]{64}$/.test(id)) throw new Error("ELITE_SOURCE_SHA256 must bind the source inventory");
    return id;
  },
  poweredByHeader: false,
  typedRoutes: true,
  // This reference does not transform images at runtime. Re-admit an image
  // pipeline and its dependencies before enabling the built-in optimizer.
  images: { unoptimized: true },
  async headers() {
    return [
      {
        source: "/(.*)",
        headers: [
          { key: "X-Content-Type-Options", value: "nosniff" },
          { key: "Referrer-Policy", value: "strict-origin-when-cross-origin" },
          { key: "Permissions-Policy", value: "camera=(), microphone=(), geolocation=()" },
          { key: "X-Frame-Options", value: "DENY" },
          { key: "Cross-Origin-Opener-Policy", value: "same-origin" }
        ]
      }
    ];
  }
};

export default nextConfig;
