'use client';

import React, { useState, useRef, useEffect } from 'react';
import Link from 'next/link';
import {
  ChevronDown,
  Settings,
  Volume2,
  Globe,
  Palette,
  LogOut,
  User,
  Home,
  Puzzle,
  BookOpen,
  Users
} from 'lucide-react';

interface NavbarProps {
  user?: {
    id: number;
    username: string;
    email: string;
  } | null;
  onLogout?: () => void;
}

const Navbar: React.FC<NavbarProps> = ({ user, onLogout }) => {
  const [isAccountDropdownOpen, setIsAccountDropdownOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsAccountDropdownOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const handleLogout = () => {
    onLogout?.();
    setIsAccountDropdownOpen(false);
  };

  return (
    <nav className="w-full border-b border-gray-300 bg-white shadow-sm relative z-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">

          {/* Logo */}
          <Link href="/" className="flex-shrink-0">
            <div className="text-gray-900 text-xl font-bold hover:text-blue-600 transition-colors">
              KANGYOO.com
            </div>
          </Link>

          {/* Navigation Menu */}
          <div className="hidden md:block">
            <div className="flex items-center space-x-8">

              {/* Chơi */}
              <Link
                href="/play"
                className="text-gray-700 hover:text-blue-600 transition-colors flex items-center space-x-1"
              >
                <Home className="h-4 w-4" />
                <span>Chơi</span>
              </Link>

              {/* Câu đố */}
              <Link
                href="/puzzles"
                className="text-gray-700 hover:text-blue-600 transition-colors flex items-center space-x-1"
              >
                <Puzzle className="h-4 w-4" />
                <span>Câu đố</span>
              </Link>

              {/* Học */}
              <Link
                href="/learn"
                className="text-gray-700 hover:text-blue-600 transition-colors flex items-center space-x-1"
              >
                <BookOpen className="h-4 w-4" />
                <span>Học</span>
              </Link>

              {/* Cộng đồng */}
              <Link
                href="/community"
                className="text-gray-700 hover:text-blue-600 transition-colors flex items-center space-x-1"
              >
                <Users className="h-4 w-4" />
                <span>Cộng đồng</span>
              </Link>

              {/* Cài đặt */}
              <Link
                href="/settings"
                className="text-gray-700 hover:text-blue-600 transition-colors flex items-center space-x-1"
              >
                <Settings className="h-4 w-4" />
                <span>Cài đặt</span>
              </Link>

            </div>
          </div>

          {/* Account Section */}
          <div className="relative" ref={dropdownRef}>
            {user ? (
              // Logged in user
              <div className="relative">
                <button
                  className="flex items-center space-x-2 text-gray-700 hover:text-blue-600 transition-colors focus:outline-none"
                  onClick={() => setIsAccountDropdownOpen(!isAccountDropdownOpen)}
                >
                  <User className="h-5 w-5" />
                  <span className="hidden sm:block">{user.username}</span>
                  <ChevronDown className={`h-4 w-4 transition-transform ${isAccountDropdownOpen ? 'rotate-180' : ''}`} />
                </button>

                {/* Dropdown Menu */}
                {isAccountDropdownOpen && (
                  <div className="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-gray-200 z-50">
                    <div className="py-1">

                      {/* User Info */}
                      <div className="px-4 py-2 border-b border-gray-200">
                        <div className="text-gray-900 font-medium">{user.username}</div>

                      </div>

                      {/* Ngôn ngữ */}
                      <button className="w-full px-4 py-2 text-left text-gray-700 hover:bg-gray-50 flex items-center space-x-2">
                        <Globe className="h-4 w-4" />
                        <span>Ngôn ngữ</span>
                      </button>

                      {/* Âm thanh */}
                      <button className="w-full px-4 py-2 text-left text-gray-700 hover:bg-gray-50 flex items-center space-x-2">
                        <Volume2 className="h-4 w-4" />
                        <span>Âm thanh</span>
                      </button>

                      {/* Nền */}
                      <button className="w-full px-4 py-2 text-left text-gray-700 hover:bg-gray-50 flex items-center space-x-2">
                        <Palette className="h-4 w-4" />
                        <span>Nền</span>
                      </button>

                      {/* Divider */}
                      <div className="border-t border-gray-200 my-1"></div>

                      {/* Đăng xuất */}
                      <button
                        onClick={handleLogout}
                        className="w-full px-4 py-2 text-left text-red-600 hover:bg-gray-50 hover:text-red-700 flex items-center space-x-2"
                      >
                        <LogOut className="h-4 w-4" />
                        <span>Đăng xuất</span>
                      </button>

                    </div>
                  </div>
                )}
              </div>
            ) : (
              // Not logged in
              <div className="flex items-center space-x-4">
                <Link
                  href="/login"
                  className="text-gray-700 hover:text-blue-600 transition-colors"
                >
                  Đăng nhập
                </Link>
                <Link
                  href="/register"
                  className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg transition-colors"
                >
                  Đăng ký
                </Link>
              </div>
            )}
          </div>

          {/* Mobile menu button */}
          <div className="md:hidden">
            <button className="text-gray-700 hover:text-blue-600">
              <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            </button>
          </div>

        </div>
      </div>
    </nav>
  );
};

export default Navbar;