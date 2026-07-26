'use client';

import React from 'react';
import Image from 'next/image';

interface PlayerStats {
  username: string;
  avatar?: string;
  accuracy: number;
  brilliant: number;
  great: number;
  best: number;
  excellent: number;
  good: number;
  inaccuracy: number;
  mistake: number;
  blunder: number;
  missedWin: number;
}

interface GameAnalysisProps {
  player1: PlayerStats;
  player2: PlayerStats;
  winner?: 'player1' | 'player2' | 'draw';
  className?: string;
  onShowReplay?: () => void;
}

const GameAnalysis: React.FC<GameAnalysisProps> = ({
  player1,
  player2,
  winner,
  className = '',
  onShowReplay
}) => {
  // Statistics configuration
  const stats = [
    { 
      label: 'Thiên tài', 
      icon: '/assets/brilliant_1024x.png', 
      key: 'brilliant' 
    },
    { 
      label: 'Tuyệt', 
      icon: '/assets/great_find_1024x.png', 
      key: 'great' 
    },
    { 
      label: 'Tốt nhất', 
      icon: '/assets/best_1024x.png', 
      key: 'best' 
    },
    { 
      label: 'Tuyệt vời', 
      icon: '/assets/excellent_1024x.png', 
      key: 'excellent' 
    },
    { 
      label: 'Tốt', 
      icon: '/assets/correct_1024x.png', 
      key: 'good' 
    },
    { 
      label: 'Chủ đề sách', 
      icon: '/assets/book_1024x.png', 
      key: 'inaccuracy' 
    },
    { 
      label: 'Không chính xác', 
      icon: '/assets/inaccuracy_1024x.png', 
      key: 'mistake' 
    },
    { 
      label: 'Sai lầm', 
      icon: '/assets/mistake_1024x.png', 
      key: 'mistake' 
    },
    { 
      label: 'Bỏ lỡ', 
      icon: '/assets/incorrect_1024x.png', 
      key: 'miss' 
    },
    { 
      label: 'Nước sai lầm nghiêm trọng', 
      icon: '/assets/blunder_1024x.png', 
      key: 'blunder' 
    }
  ];

  // Helper function to get winner indicator
  const getWinnerIndicator = () => {
    switch (winner) {
      case 'player1': return '1-0';
      case 'player2': return '0-1';
      default: return '½-½';
    }
  };

  const getAccuracyBgColor = (isPlayer2: boolean = false) => {
    if (isPlayer2 && winner === 'player2') {
      return 'bg-green-600';
    }
    if (!isPlayer2 && winner === 'player1') {
      return 'bg-green-600';
    }
    return 'bg-gray-700';
  };

  return (
    <div className={`bg-white border border-gray-300 rounded-lg p-4 mt-4 ${className}`}>
      <div className="mb-6">
        {/* Player Names & Result */}
        <div className="flex items-center justify-between mb-4">
          <div className="text-center flex-1">
            <h3 className="text-sm font-semibold text-gray-800">
              {player1.username}
            </h3>
          </div>
          
          <div className="text-center px-4">
            <div className="text-xl font-bold text-gray-900">
              {getWinnerIndicator()}
            </div>
          </div>
          
          <div className="text-center flex-1">
            <h3 className="text-sm font-semibold text-gray-800">
              {player2.username}
            </h3>
          </div>
        </div>

        {/* Accuracy Comparison */}
        <div className="space-y-2">
          <h4 className="text-sm text-gray-600 font-medium">
            Độ chính xác
          </h4>
          
          <div className="flex gap-4">
            <div className={`${getAccuracyBgColor()} px-4 py-2 rounded-lg text-center text-white flex-1`}>
              <div className="font-bold text-lg">
                {player1.accuracy}%
              </div>
            </div>
            
            <div className={`${getAccuracyBgColor(true)} px-4 py-2 rounded-lg text-center text-white flex-1`}>
              <div className="font-bold text-lg">
                {player2.accuracy}%
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="space-y-2">
        <h4 className="text-sm text-gray-600 font-medium mb-3">
          Thống kê nước đi
        </h4>
        
        {stats.map((stat, index) => (
          <div 
            key={index} 
            className="flex items-center justify-between py-1 hover:bg-gray-50 rounded px-2 transition-colors"
          >
            {/* Label */}
            <div className="flex items-center gap-2 flex-1">
              <span className="text-xs text-gray-700 font-medium">
                {stat.label}
              </span>
            </div>
            
            {/* Stats Row */}
            <div className="flex items-center gap-3">
              {/* Player 1 Count */}
              <span className="w-8 text-center text-sm font-semibold text-gray-800">
                {player1[stat.key as keyof PlayerStats] || 0}
              </span>
              
              {/* Icon */}
              <div className="w-6 h-6 flex items-center justify-center">
                <Image
                  src={stat.icon}
                  alt={stat.label}
                  width={20}
                  height={20}
                  className="object-contain"
                />
              </div>
              
              {/* Player 2 Count */}
              <span className="w-8 text-center text-sm font-semibold text-gray-800">
                {player2[stat.key as keyof PlayerStats] || 0}
              </span>
            </div>
          </div>
        ))}
      </div>

      <div className="mt-6">
        <button 
          onClick={onShowReplay}
          className="w-full bg-green-600 hover:bg-green-700 active:bg-green-800 text-white font-semibold py-3 px-4 rounded-lg transition-all duration-200 transform hover:scale-[1.02] active:scale-[0.98] shadow-sm hover:shadow-md"
        >
          Xem lại
        </button>
      </div>
    </div>
  );
};

export default GameAnalysis;