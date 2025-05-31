/** @type {import('next').NextConfig} */
const nextConfig = {
  basePath: '/meet/hello-world-next-js',
  distDir: '../dist/hello-world-next-js',
  output: 'export',
  experimental: {
    serverComponentsExternalPackages: []
  }
};

export default nextConfig;
