import React from 'react';
import Image from 'next/image';

export type MoveLabel = 'best' | 'blunder' | 'book' | 'brilliant' | 'correct' | 'excellent' | 'inaccuracy' | 'incorrect' | 'mistake' | 'good';

interface SquareProps {
  isLight: boolean;
  onClick?: () => void;
  className?: string;
  children?: React.ReactNode; 
  isSelected: boolean;
  isValidMove: boolean;
  // New props for move highlighting
  isFromSquare?: boolean;
  isToSquare?: boolean;
  moveLabel?: MoveLabel;
}

const Square = ({ 
  isLight, 
  onClick, 
  className = '', 
  children, 
  isSelected, 
  isValidMove,
  isFromSquare = false,
  isToSquare = false,
  moveLabel
}: SquareProps) => {
  
  // Label colors mapping
  const labelColors = {
    best: 'bg-lime-400/60',           // xanh lá cây đậm
    blunder: 'bg-red-500/60',          // màu đỏ
    book: 'bg-amber-700/60',           // màu nâu
    brilliant: 'bg-cyan-500/60',       // màu xanh nước biển
    correct: 'bg-lime-400/60',        // màu xanh lá cây nhạt
    excellent: 'bg-lime-400/60',      // màu xanh lá cây nhạt
    inaccuracy: 'bg-yellow-400/60',    // màu vàng
    incorrect: 'bg-red-500/60',        // màu đỏ
    mistake: 'bg-orange-500/60',       // màu cam
    good: 'bg-green-600/60'            // màu xanh lá cây
  };

  // Label assets mapping
  const labelIcons = {
    best: '/assets/best_1024x.png',
    blunder: '/assets/blunder_1024x.png',
    book: '/assets/book_1024x.png',
    brilliant: '/assets/brilliant_1024x.png',
    correct: '/assets/correct_1024x.png',
    excellent: '/assets/excellent_1024x.png',
    inaccuracy: '/assets/inaccuracy_1024x.png',
    incorrect: '/assets/incorrect_1024x.png',
    mistake: '/assets/mistake_1024x.png',
    good: '/assets/correct_1024x.png' // Tạm dùng correct icon cho good
  };

  const baseClasses = `
    w-20 h-20 
    flex items-center justify-center 
    cursor-pointer 
    transition-all duration-200 
    hover:brightness-110
    relative
    ${isLight ? 'bg-[#f0d9b5]' : 'bg-[#b58863]'}
    ${isSelected ? 'ring-4 ring-yellow-400 ring-inset' : ''}
    ${isValidMove ? 'ring-4 ring-green-400 ring-inset' : ''}
  `;
  
  const combinedClasses = className ? `${baseClasses} ${className}` : baseClasses;

  // Get highlight overlay for move squares
  const getMoveHighlight = () => {
    if ((isFromSquare || isToSquare) && moveLabel) {
      return (
        <div className={`absolute inset-0 ${labelColors[moveLabel]} rounded-sm`} />
      );
    }
    return null;
  };

  // Get label icon for to square - Updated with larger size and positioned outside
  const getLabelIcon = () => {
    if (isToSquare && moveLabel) {
      return (
        <div className="absolute -top-2 -right-2 w-8 h-8 z-20">
          <Image
            src={labelIcons[moveLabel]}
            alt={moveLabel}
            width={32}
            height={32}
            className="w-full h-full object-contain drop-shadow-lg"
          />
        </div>
      );
    }
    return null;
  };

  return (
    <div 
      className={combinedClasses}
      onClick={onClick}
    >
      {/* Move highlight overlay */}
      {getMoveHighlight()}
      
      {/* Piece */}
      <div className="relative z-10">
        {children}
      </div>
      
      {/* Valid move indicator */}
      {isValidMove && !children && (
        <div className="w-6 h-6 bg-green-600 rounded-full opacity-70 z-10" />
      )}
      
      {/* Label icon */}
      {getLabelIcon()}
    </div>
  );
};

export default Square;