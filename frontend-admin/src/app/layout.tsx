import type { Metadata } from 'next';
import { Geist, Geist_Mono } from 'next/font/google';
import './globals.css';

const geistSans = Geist({ variable: '--font-geist-sans', subsets: ['latin'] });
const geistMono = Geist_Mono({ variable: '--font-geist-mono', subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'Lamina Admin Panel',
  description: 'Invite and manage users',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className="bg-brand-beige text-brand-grey-dark">
      <body className={`${geistSans.variable} ${geistMono.variable} font-sans`}>
        {children}
      </body>
    </html>
  );
}
