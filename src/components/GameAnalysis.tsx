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
}

const GameAnalysis: React.FC<GameAnalysisProps> = ({
  player1,
  player2,
  winner,
  className = ''
}) => {
  const stats = [
    { 
      label: 'Thiên tài', 
      icon: '/assets/brilliant_1024x.png', 
      color: 'bg-cyan-500/20 border-cyan-500/30', 
      key: 'brilliant' 
    },
    { 
      label: 'Tuyệt', 
      icon: '/assets/excellent_1024x.png', 
      color: 'bg-blue-500/20 border-blue-500/30', 
      key: 'great' 
    },
    { 
      label: 'Tốt nhất', 
      icon: '/assets/best_1024x.png', 
      color: 'bg-green-500/20 border-green-500/30', 
      key: 'best' 
    },
    { 
      label: 'Tuyệt vời', 
      icon: '/assets/excellent_1024x.png', 
      color: 'bg-green-400/20 border-green-400/30', 
      key: 'excellent' 
    },
    { 
      label: 'Tốt', 
      icon: '/assets/correct_1024x.png', 
      color: 'bg-green-600/20 border-green-600/30', 
      key: 'good' 
    },
    { 
      label: 'Chủ đề sách', 
      icon: '/assets/book_1024x.png', 
      color: 'bg-orange-500/20 border-orange-500/30', 
      key: 'inaccuracy' 
    },
    { 
      label: 'Không chính xác', 
      icon: '/assets/inaccuracy_1024x.png', 
      color: 'bg-yellow-500/20 border-yellow-500/30', 
      key: 'mistake' 
    },
    { 
      label: 'Sai lầm', 
      icon: '/assets/mistake_1024x.png', 
      color: 'bg-orange-600/20 border-orange-600/30', 
      key: 'mistake' 
    },
    { 
      label: 'Bỏ lỡ', 
      icon: '/assets/blunder_1024x.png', 
      color: 'bg-red-500/20 border-red-500/30', 
      key: 'blunder' 
    },
    { 
      label: 'Nước sai làm ngỡ ngàn', 
      icon: '/assets/incorrect_1024x.png', 
      color: 'bg-red-600/20 border-red-600/30', 
      key: 'missedWin' 
    }
  ];

  return (
    <div className="bg-white border border-gray-300 rounded-lg p-3 mt-4">
      {/* Game Result Chart */}
      <div className="mb-6">
        <div className="flex items-center justify-between mb-3">
          <div className="text-center">
            <div className="text-sm font-medium">{player1.username}</div>
          </div>
          
          <div className="text-center text-xl font-bold">
            {winner === 'player1' ? '1-0' : winner === 'player2' ? '0-1' : '½-½'}
          </div>
          
          <div className="text-center">
     
              
            <div className="text-sm font-medium">{player2.username}</div>
          </div>
        </div>

        {/* Accuracy */}
        <div className="mb-4">
          <div className="text-sm text-black-400 mb-2">Độ chính xác</div>
          <div className="flex justify-between">
            <div className="bg-gray-700 px-3 py-2 rounded text-center">
              <div className="font-bold">{player1.accuracy}</div>
            </div>
            <div className={`px-3 py-2 rounded text-center ${
              winner === 'player2' ? 'bg-green-600' : 'bg-gray-700'
            }`}>
              <div className="font-bold">{player2.accuracy}</div>
            </div>
          </div>
        </div>
      </div>

      {/* Statistics */}
      <div className="space-y-3">
        {stats.map((stat, index) => (
          <div key={index} className="flex items-center justify-between">
            <div className="flex items-center gap-2 flex-1">
              <span className="text-xs">{stat.label}</span>
            </div>
            
            <div className="flex items-center gap-4">
              <span className="w-8 text-center text-sm font-medium">
                {player1[stat.key as keyof PlayerStats] || 0}
              </span>
              
              <div className={`w-8 h-8  flex items-center justify-center`}>
                <Image
                  src={stat.icon}
                  alt={stat.label}
                  width={20}
                  height={20}
                  className="object-contain"
                />
              </div>
              
              <span className="w-8 text-center text-sm font-medium">
                {player2[stat.key as keyof PlayerStats] || 0}
              </span>
            </div>
          </div>
        ))}
      </div>

      {/* Action Button */}
      <button className="w-full mt-6 bg-green-600 hover:bg-green-700 text-white font-bold py-3 rounded-lg transition-colors">
        Xem lại
      </button>
    </div>
  );
};

export default GameAnalysis;