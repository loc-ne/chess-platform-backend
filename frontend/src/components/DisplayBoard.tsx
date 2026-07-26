import React from 'react';
import Square from './Square';
import Piece from './Piece';
import { ChessEngine } from '../game/Board';
import { GameStateManager } from '../game/GameState';

const STARTING_FEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1";

interface DisplayBoardProps {
  playerColor: 'white' | 'black';
  initialFen?: string;
  size?: 'small' | 'medium' | 'large';
  className?: string;
}

const DisplayBoard: React.FC<DisplayBoardProps> = ({
  playerColor = 'white',
  initialFen = STARTING_FEN,
  size = 'medium',
  className = ''
}) => {
  
  const board = React.useMemo(() => {
    const gameState = GameStateManager.convertFenToClientState(initialFen);
    return ChessEngine.bitboardToBoard(gameState.bitboards);
  }, [initialFen]);

  const isLightSquare = React.useCallback((row: number, col: number) => {
    return (row + col) % 2 !== 0;
  }, []);


  const renderBoard = () => {
    const rows = playerColor === 'black' ? [0, 1, 2, 3, 4, 5, 6, 7] : [7, 6, 5, 4, 3, 2, 1, 0];
    const cols = playerColor === 'black' ? [7, 6, 5, 4, 3, 2, 1, 0] : [0, 1, 2, 3, 4, 5, 6, 7];

    return (
      <div className={`inline-block bg-amber-900 border-4 border-gray-800  ${className}`}>
        {rows.map(row => (
          <div key={row} className="flex h-1/8">
            {cols.map(col => {
              const piece = board[row][col];

              return (
                <Square
                  key={`${row}-${col}`}
                  isLight={isLightSquare(row, col)}
                  onClick={() => { }}
                  isSelected={false}
                  isValidMove={false}
                >
                  {piece && (
                    <div className="pointer-events-none">
                      <Piece type={piece.type} color={piece.color} />
                    </div>
                  )}
                </Square>
              );
            })}
          </div>
        ))}
      </div>
    );
  };

  return (
    <div className="relative">
      {renderBoard()}
    </div>
  );
};

export default React.memo(DisplayBoard);