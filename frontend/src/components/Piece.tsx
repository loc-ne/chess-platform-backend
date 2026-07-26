import React from 'react';
import Image from 'next/image';

interface PieceProps {
  type: 'knight' | 'rook' | 'bishop' | 'queen' | 'king' | 'pawn';
  color: 'white' | 'black';
}

const Piece = ({ type, color }: PieceProps) => {
  const getImagePath = () => {
    const colorPrefix = color === 'white' ? 'w' : 'b';
    const typeMap = {
      king: 'k',
      queen: 'q', 
      rook: 'r',
      bishop: 'b',
      knight: 'n',
      pawn: 'p'
    };
    
    return `/assets/${colorPrefix}${typeMap[type]}.png`; 
  };

  const getAltText = () => {
    const colorName = color === 'white' ? 'White' : 'Black';
    const typeName = type.charAt(0).toUpperCase() + type.slice(1);
    return `${colorName} ${typeName}`;
  };


  return (
    < Image 
      src={getImagePath()}
      alt={getAltText()}
      className="w-20 h-20 select-none"
      width={80}
      height={80}
      draggable={false}
    />
  );
};

export default Piece;