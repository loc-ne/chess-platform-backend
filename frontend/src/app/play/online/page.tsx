'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import TimeSelector from '@/components/TimeSelector';
import DisplayBoard from '@/components/DisplayBoard';
import { useAuth } from '@/contexts/AuthContext';

const PlayOnlinePage: React.FC = () => {
  const { user } = useAuth();
  const router = useRouter();

  const [selectedTime, setSelectedTime] = useState('5 + 0');
  const [selectedCategory, setSelectedCategory] = useState('Blitz');
  const [isSearching, setIsSearching] = useState(false);
  const [searchTime, setSearchTime] = useState(0);
  const currentFen = 'rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1';
  const wsRef = React.useRef<WebSocket | null>(null);

  // Timer for search duration
  useEffect(() => {
    let interval: NodeJS.Timeout;
    if (isSearching) {
      interval = setInterval(() => {
        setSearchTime(prev => prev + 1);
      }, 1000);
    } else {
      setSearchTime(0);
    }
    return () => clearInterval(interval);
  }, [isSearching]);

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  const handleTimeSelect = (time: string, category: string) => {
    setSelectedTime(time);
    setSelectedCategory(category);
  };

  function toTimeControlDto(timeStr: string, category: string) {
    const [minStr, incStr] = timeStr.replace(/\s+/g, '').split('+');
    const initialTime = parseInt(minStr, 10) * 60;
    const increment = incStr ? parseInt(incStr, 10) : 0;
    const type = ['bullet', 'blitz', 'rapid', 'classical'].includes(category.toLowerCase())
      ? category.toLowerCase()
      : 'blitz';
    return { type, initialTime, increment };
  }

  const handleStartSearch = async () => {
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_AUTH_SERVICE_URL}/api/v1/auth/me`, {
        credentials: 'include'
      });

      const data = await response.json();
      console.log(data);

      const timeControl = toTimeControlDto(selectedTime, selectedCategory);
      const body = {
        timeControl,
        player: {
          userId: data.data.id,
          username: data.data.username
        }
      };

      const matchResponse = await fetch(`${process.env.NEXT_PUBLIC_MATCHMAKING_SERVICE_URL}/api/v1/matchmaking/join`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(body)
      });

      const result = await matchResponse.json();
      console.log(result);

      if (matchResponse.ok) {
        setIsSearching(true);
      }
    } catch (error) {
      console.error('Find game error:', error);
    }
  };

  const handleCancelSearch = async () => {
    try {
      const response = await fetch(`${process.env.NEXT_PUBLIC_AUTH_SERVICE_URL}/api/v1/auth/me`, {
        credentials: 'include'
      });

      const data = await response.json();

      const timeControl = toTimeControlDto(selectedTime, selectedCategory);
      const body = {
        timeControl,
        player: {
          userId: data.data.id,
          username: data.data.username
        }
      };

      await fetch(`${process.env.NEXT_PUBLIC_MATCHMAKING_SERVICE_URL}/api/v1/matchmaking/leave`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(body)
      });

      setIsSearching(false);
      setSearchTime(0);
    } catch (error) {
      console.error('Cancel search error:', error);
      setIsSearching(false);
      setSearchTime(0);
    }
  };

  // WebSocket connection
  useEffect(() => {
    if (!user?.id || wsRef.current) return;

    const wsUrl = process.env.NEXT_PUBLIC_WEBSOCKET_URL || 'ws://localhost:3005/ws';
    const ws = new WebSocket(wsUrl);
    wsRef.current = ws;

    ws.onopen = () => {
      console.log('Connected to WebSocket');
      ws.send(JSON.stringify({
        type: 'joinMatchmaking',
        userId: user.id,
        username: user.username
      }));
      console.log('Joined matchmaking queue');
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      console.log('Received:', data);

      if (data.type === 'matchFound') {
        setIsSearching(false);
        router.push(`/game/${data.roomId}`);
      }
    };

    ws.onclose = () => {
      console.log('WebSocket connection closed');
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    return () => {
      ws.close();
      wsRef.current = null;
    };
  }, [user, router]);

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">

          <div className="lg:col-span-2">
            <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
              <div className="flex justify-center">
                <DisplayBoard
                  playerColor='white'
                  initialFen={currentFen}
                />
              </div>
            </div>
          </div>

          <div className="lg:col-span-1">
            <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
              {!isSearching ? (
                <div>
                  <div className="mb-6">
                    <TimeSelector
                      initialTime={selectedTime}
                      initialCategory={selectedCategory}
                      onTimeSelect={handleTimeSelect}
                      disabled={isSearching}
                      size="medium"
                    />
                  </div>
                  <button
                    className="w-full bg-gradient-to-r from-blue-500 to-blue-600 hover:from-blue-600 hover:to-blue-700 text-white font-bold text-lg rounded-xl py-4 transition-all duration-300 shadow-lg hover:shadow-xl transform hover:scale-105 flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none disabled:shadow-none"
                    onClick={handleStartSearch}
                  >
                    Bắt đầu ván cờ
                  </button>
                </div>
              ) : (
                <div className="text-center py-8">
                  {/* Animated KANGYOO Text */}
                  <div className="mb-8">
                    <div className="flex items-center justify-center mb-4">
                      <span className="text-4xl font-bold text-gray-800 animate-pulse">
                        KANGYOO
                      </span>
                    </div>

                    {/* Subtitle */}
                    <p className="text-lg text-blue-600 animate-fade-in-up" style={{ animationDelay: '0.5s' }}>
                      Đang tìm đối thủ
                    </p>
                  </div>

                  {/* Search Animation Container */}
                  <div className="bg-blue-50 border border-blue-200 rounded-2xl p-6 mb-6">

                    {/* Search Stats */}
                    <div className="space-y-3">
                      <div className="flex items-center justify-center gap-2">
                        <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
                        <span className="text-gray-700 text-sm font-medium">
                          Tìm kiếm {selectedCategory} {selectedTime}
                        </span>
                      </div>

                      <div className="text-gray-500 text-sm">
                        Thời gian: {formatTime(searchTime)}
                      </div>


                    </div>
                  </div>

                  {/* Cancel Button */}
                  <button
                    className="w-full bg-gradient-to-r from-red-500 to-red-600 hover:from-red-600 hover:to-red-700 text-white font-bold py-3 rounded-xl transition-all duration-300 shadow-lg hover:shadow-xl transform hover:scale-105"
                    onClick={handleCancelSearch}
                  >
                    Hủy tìm kiếm
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>

      </div>
    </div>


  );
};

export default PlayOnlinePage;