import React from 'react';
import Square from './Square';
import Piece from './Piece';
import { ChessEngine } from '../game/Board';
import { Position, Piece as PieceInterface } from '../game/types';
import { ClientGameState, GameStateManager } from '../game/GameState';

const STARTING_FEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1";

// Define MoveLabel type
export type MoveLabel = 'brilliant' | 'best' | 'excellent' | 'good' | 'book' | 'inaccuracy' | 'mistake' | 'blunder';

interface BoardProps {
  playerColor: 'white' | 'black';
  initialFen?: string;
  onMove?: (from: Position, to: Position) => void;
  gameState?: ClientGameState | null;
  // New prop for replay highlighting
  highlightedMove?: {
    from: Position;
    to: Position;
    label: MoveLabel;
  };
}

const Board: React.FC<BoardProps> = ({
  playerColor = 'white',
  initialFen = STARTING_FEN,
  onMove,
  gameState: externalGameState,
  highlightedMove, // New prop
}) => {
  //   1. Local game state 
  const [gameState, setGameState] = React.useState<ClientGameState>(() =>
    externalGameState || GameStateManager.convertFenToClientState(initialFen)
  );

  const [selectedSquare, setSelectedSquare] = React.useState<Position | null>(null);
  const [validMoves, setValidMoves] = React.useState<Position[]>([]);

  //   2. Premove support 
  const [premoveQueue, setPremoveQueue] = React.useState<[Position, Position][]>([]);

  //   3. Drag and drop state
  const [dragState, setDragState] = React.useState<{
    isDragging: boolean;
    piece: PieceInterface | null;
    from: Position;
    startCoords: { x: number; y: number } | null;
  } | null>(null);

  const [ghostPieceStyle, setGhostPieceStyle] = React.useState<React.CSSProperties>({
    display: 'none'
  });

  React.useEffect(() => {
    if (externalGameState) {
      setGameState(externalGameState);
    }
  }, [externalGameState]);

  const boardRef = React.useRef<HTMLDivElement>(null);

  //   4. Memoized board derivation
  const board = React.useMemo(() => {
    return ChessEngine.bitboardToBoard(gameState.bitboards);
  }, [gameState.bitboards]);

  //   5. Move generation cache
  const validMovesCache = React.useRef(new Map<string, Position[]>());

  const getValidMoves = React.useCallback((position: Position): Position[] => {
    const key = `${position.row}-${position.col}-${gameState.currentFen}`;
    if (validMovesCache.current.has(key)) {
      return validMovesCache.current.get(key)!;
    }

    const moves = ChessEngine.generateMovesForPiece(
      gameState.bitboards,
      gameState,
      position
    );

    if (validMovesCache.current.size > 100) {
      validMovesCache.current.clear();
    }

    validMovesCache.current.set(key, moves);
    return moves;
  }, [gameState]);

  const handleMovePiece = React.useCallback((from: Position, to: Position) => {
    onMove?.(from, to);

    setSelectedSquare(null);
    setValidMoves([]);
    validMovesCache.current.clear();
  }, [onMove]);

  const handleSquareClick = (row: number, col: number) => {
    const clickedPosition = { row, col };

    // Deselect if clicking same square
    if (selectedSquare && selectedSquare.row === row && selectedSquare.col === col) {
      setSelectedSquare(null);
      setValidMoves([]);
      return;
    }

    // Make move if clicking valid move square
    if (selectedSquare && validMoves.some(move => move.row === row && move.col === col)) {
      if (gameState.activeColor === playerColor) {
        handleMovePiece(selectedSquare, clickedPosition);
      } else {
        setPremoveQueue(prev => [...prev, [selectedSquare, clickedPosition]]);
        setSelectedSquare(null);
        setValidMoves([]);
      }
      return;
    }

    // Select piece
    const piece = ChessEngine.getPieceAt(gameState.bitboards, clickedPosition);
    if (piece && piece.color === playerColor) {
      setSelectedSquare(clickedPosition);
      const moves = getValidMoves(clickedPosition);
      setValidMoves(moves);
    } else {
      setSelectedSquare(null);
      setValidMoves([]);
    }
  };

  // ========== HIGHLIGHT HELPER FUNCTIONS ==========
  
  // Check if square is part of highlighted move (from square)
  const isFromSquare = React.useCallback((row: number, col: number): boolean => {
    return highlightedMove?.from.row === row && highlightedMove?.from.col === col;
  }, [highlightedMove]);

  // Check if square is part of highlighted move (to square)
  const isToSquare = React.useCallback((row: number, col: number): boolean => {
    return highlightedMove?.to.row === row && highlightedMove?.to.col === col;
  }, [highlightedMove]);

  //   9. Helper functions
  const isLightSquare = React.useCallback((row: number, col: number) => {
    return (row + col) % 2 != 0;
  }, []);

  const isSquareSelected = React.useCallback((row: number, col: number) => {
    return selectedSquare !== null && selectedSquare.row === row && selectedSquare.col === col;
  }, [selectedSquare]);

  const isSquareValidMove = React.useCallback((row: number, col: number) => {
    return validMoves.some(move => move.row === row && move.col === col);
  }, [validMoves]);

  //   10. Drag and drop helpers
  const getSquareFromCoordinates = React.useCallback((x: number, y: number): Position | null => {
    const boardElement = boardRef.current;
    if (!boardElement) return null;

    const rect = boardElement.getBoundingClientRect();
    const relativeX = x - rect.left;
    const relativeY = y - rect.top;
    const squareSize = rect.width / 8;

    // Calculate grid position (0-7)
    const gridCol = Math.floor(relativeX / squareSize);
    const gridRow = Math.floor(relativeY / squareSize);

    // Validate grid bounds
    if (gridCol < 0 || gridCol >= 8 || gridRow < 0 || gridRow >= 8) {
      return null;
    }

    //  Convert grid position to actual board position based on player perspective
    let actualRow: number;
    let actualCol: number;

    if (playerColor === 'white') {
      // White perspective: top-left is (0,0)
      actualRow = 7 - gridRow;
      actualCol = gridCol;
    } else {
      // Black perspective: board is flipped
      actualRow = gridRow;  // Flip vertically
      actualCol = 7 - gridCol;  // Flip horizontally
    }

    return { row: actualRow, col: actualCol };
  }, [playerColor]);

  const handlePieceMouseDown = (
    e: React.MouseEvent,
    piece: PieceInterface,
    piecePosition: Position
  ) => {
    // if (!canInteract) return;

    const canDragPiece = piece.color === playerColor;
    if (!canDragPiece) return;

    e.preventDefault();

    setDragState({
      isDragging: true,
      piece,
      from: piecePosition,
      startCoords: { x: e.clientX, y: e.clientY }
    });

    const moves = getValidMoves(piecePosition);
    setValidMoves(moves);

    const isPlayerTurn = gameState.activeColor === playerColor;

    setGhostPieceStyle({
      display: 'block',
      position: 'fixed',
      left: e.clientX - 40,
      top: e.clientY - 40,
      zIndex: 9999,
      pointerEvents: 'none',
      opacity: 0.9,
      transform: 'scale(1.1)',
      filter: isPlayerTurn ? 'drop-shadow(0 4px 8px rgba(0,0,0,0.3))' : 'sepia(0.3)'
    });
  };

  React.useEffect(() => {
    if (!dragState?.isDragging) return;

    const handleMouseMove = (e: MouseEvent) => {
      setGhostPieceStyle(prev => ({
        ...prev,
        left: e.clientX - 40,
        top: e.clientY - 40
      }));
    };

    const handleMouseUp = (e: MouseEvent) => {
      const dropSquare = getSquareFromCoordinates(e.clientX, e.clientY);
      console.log(dropSquare);
      if (dropSquare && dragState.from) {
        const isValid = validMoves.some(move =>
          move.row === dropSquare.row && move.col === dropSquare.col
        );

        if (isValid) {
          if (gameState.activeColor === playerColor) {

            handleMovePiece(dragState.from, dropSquare);
          } else {

            setPremoveQueue(prev => [...prev, [dragState.from, dropSquare]]);
            setSelectedSquare(null);
            setValidMoves([]);
          }
        }
      }

      // Cleanup
      setDragState(null);
      setGhostPieceStyle({ display: 'none' });

      if (!dropSquare || !validMoves.some(move =>
        move.row === dropSquare.row && move.col === dropSquare.col
      )) {
        setValidMoves([]);
      }
    };

    document.addEventListener('mousemove', handleMouseMove);
    document.addEventListener('mouseup', handleMouseUp);

    return () => {
      document.removeEventListener('mousemove', handleMouseMove);
      document.removeEventListener('mouseup', handleMouseUp);
    };
  }, [
    playerColor,
    gameState.activeColor,
    dragState,
    validMoves,
    handleMovePiece,
    getSquareFromCoordinates,
  ]);

  //   12. Piece visibility helper
  const isPieceVisible = React.useCallback((piecePosition: Position) => {
    if (!dragState?.isDragging || !dragState.from) return true;
    return !(dragState.from.row === piecePosition.row &&
      dragState.from.col === piecePosition.col);
  }, [dragState]);

  //   13. Ghost piece renderer
  const renderGhostPiece = () => {
    if (!dragState?.isDragging || !dragState.piece) return null;

    return (
      <div style={ghostPieceStyle}>
        <Piece
          type={dragState.piece.type}
          color={dragState.piece.color}
        />
      </div>
    );
  };

  //   14. Execute premoves - SIMPLIFIED
  React.useEffect(() => {
    if (playerColor === gameState.activeColor && premoveQueue.length > 0) {

      const [from, to] = premoveQueue[0];
      const piece = ChessEngine.getPieceAt(gameState.bitboards, from);

      if (!piece || piece.color !== playerColor) {
        setPremoveQueue([]);
        return;
      }

      const currentValidMoves = ChessEngine.generateMovesForPiece(
        gameState.bitboards,
        gameState,
        from
      );

      const isMoveValid = currentValidMoves.some(move =>
        move.row === to.row && move.col === to.col
      );

      if (isMoveValid) {
        setPremoveQueue(prev => prev.slice(1));

        setTimeout(() => {
          handleMovePiece(from, to);
        }, 150);
      } else {
        setPremoveQueue([]);
        setSelectedSquare(null);
        setValidMoves([]);
      }
    }
  }, [
    playerColor,
    gameState.activeColor,
    premoveQueue,
    gameState.bitboards,
    gameState,
    handleMovePiece
  ]);

  //   17. Board rendering with perspective
  const renderBoard = () => {
    const rows = playerColor === 'black' ? [0, 1, 2, 3, 4, 5, 6, 7] : [7, 6, 5, 4, 3, 2, 1, 0];
    const cols = playerColor === 'black' ? [7, 6, 5, 4, 3, 2, 1, 0] : [0, 1, 2, 3, 4, 5, 6, 7];

    return (
      <div ref={boardRef} className="inline-block bg-amber-900 border-4 border-gray-800">
        {rows.map(row => (
          <div key={row} className="flex">
            {cols.map(col => {
              const piece = board[row][col];

              return (
                <Square
                  key={`${row}-${col}`}
                  isLight={isLightSquare(row, col)}
                  onClick={() => handleSquareClick(row, col)}
                  isSelected={isSquareSelected(row, col)}
                  isValidMove={isSquareValidMove(row, col)}
                  isFromSquare={isFromSquare(row, col)}
                  isToSquare={isToSquare(row, col)}
                  moveLabel={highlightedMove?.label}
                >
                  {piece && isPieceVisible({ row, col }) && (
                    <div
                      onMouseDown={(e) => handlePieceMouseDown(e, piece, { row, col })}
                      style={{
                        cursor: piece.color === playerColor ? 'grab' : 'default'
                      }}
                    >
                      <Piece type={piece.type} color={piece.color} />
                    </div>
                  )}
                </Square>
              );
            })}
          </div>
        ))}
        {renderGhostPiece()}
      </div>
    );
  };

  return (
    <div className="relative">
      {renderBoard()}
    </div>
  );
};

export default React.memo(Board);