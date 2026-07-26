'use client';

import React from 'react';
import { useParams } from 'next/navigation';

import Board from '@/components/Board';
import Clock from '@/components/Clock';
import PlayerInfo from '@/components/PlayerInfo';
import GameResultDialog from '@/components/GameResult';
import ReplayControls from '@/components/ReplayControls';

import { useAuth } from '@/contexts/AuthContext';
import { ClientGameState, GameStateManager } from '@/game/GameState';
import { BitboardGame } from '@/game/Board';
import { Position } from '@/game/types';

interface PlayerData {
  userId: number;
  username: string;
  color: 'white' | 'black';
  rating?: number;
  isOnline: boolean;
}

interface WSStateUpdateMessage {
  type: string;
  roomId: string;
  gameState?: ClientGameState;
  player1?: PlayerData;
  player2?: PlayerData;
  whiteTimeLeft?: number;
  blackTimeLeft?: number;
  error?: string;
  result?: string;
  reason?: string;
  winner?: string;
  offerId?: string;
  offerFrom?: number;
  moveHistory?: MoveHistoryItem[];
}

interface MoveHistoryItem {
  moveNumber: number;
  white: string;
  black?: string;
}

interface MoveHistoryProps {
  moves: MoveHistoryItem[];
  indexMoveHistory: number;
  setIndexMoveHistory: React.Dispatch<React.SetStateAction<number>>;
  bitboardStates: ClientGameState[];
  setGameState: (state: ClientGameState) => void;
}

// Convert functions
function convertBitboardsToBigInt(bitboards: BitboardGame): BitboardGame {
  const keys = [
    'WhitePawns', 'WhiteRooks', 'WhiteKnights', 'WhiteBishops', 'WhiteQueens', 'WhiteKing',
    'BlackPawns', 'BlackRooks', 'BlackKnights', 'BlackBishops', 'BlackQueens', 'BlackKing'
  ];
  const result: BitboardGame = {
    whitePawns: BigInt(0),
    whiteRooks: BigInt(0),
    whiteKnights: BigInt(0),
    whiteBishops: BigInt(0),
    whiteQueens: BigInt(0),
    whiteKing: BigInt(0),
    blackPawns: BigInt(0),
    blackRooks: BigInt(0),
    blackKnights: BigInt(0),
    blackBishops: BigInt(0),
    blackQueens: BigInt(0),
    blackKing: BigInt(0),
  };
  for (const key of keys) {
    const camelKey = key.charAt(0).toLowerCase() + key.slice(1);
    result[camelKey as keyof BitboardGame] =
      bitboards && typeof bitboards === 'object' && bitboards[key as keyof BitboardGame] !== undefined
        ? BigInt(bitboards[key as keyof BitboardGame])
        : BigInt(0);
  }
  return result;
}

function convertWSGameStateToClient(wsGameState: ClientGameState): ClientGameState {
  return {
    currentFen: wsGameState.currentFen,
    bitboards: convertBitboardsToBigInt(wsGameState.bitboards),
    activeColor: wsGameState.activeColor,
    castlingRights: wsGameState.castlingRights,
    enPassantSquare: wsGameState.enPassantSquare ? {
      row: wsGameState.enPassantSquare.row,
      col: wsGameState.enPassantSquare.col
    } : null,
  };
}

const GameInfo: React.FC<{ gameType: string; timeControl: string }> = ({ gameType, timeControl }) => (
  <div className="bg-white-50 border border-blue-200 rounded-lg p-4 mb-4">
    <div className="text-lg font-bold text-blue-800 mb-2">{gameType} Game</div>
    <div className="text-sm text-blue-600 space-y-1">
      <div>Time Control: {timeControl}</div>
    </div>
  </div>
);

const MoveHistory: React.FC<MoveHistoryProps> = ({
  moves,
  indexMoveHistory,
  setIndexMoveHistory,
  bitboardStates,
  setGameState
}) => {
  const handleReplayStart = () => {
    setIndexMoveHistory(0);
    if (bitboardStates.length > 0) {
      setGameState(bitboardStates[0]);
    }
  };

  const handleReplayPrevious = () => {
    setIndexMoveHistory(prev => {
      const newIndex = Math.max(prev - 1, 0);
      if (bitboardStates.length > newIndex) {
        setGameState(bitboardStates[newIndex]);
      }
      return newIndex;
    });
  };

  const handleReplayNext = () => {
    setIndexMoveHistory(prev => {
      const newIndex = Math.min(prev + 1, bitboardStates.length - 1);
      if (bitboardStates.length > newIndex) {
        setGameState(bitboardStates[newIndex]);
      }
      return newIndex;
    });
  };

  const handleReplayEnd = () => {
    const lastIndex = bitboardStates.length - 1;
    setIndexMoveHistory(lastIndex);
    if (bitboardStates.length > lastIndex) {
      setGameState(bitboardStates[lastIndex]);
    }
  };

  return (
    <div className="bg-white border border-gray-300 rounded-lg h-96 flex flex-col">
      <div className="bg-gray-100 px-4 py-2 border-b border-gray-300 rounded-t-lg">
        <h3 className="font-semibold text-gray-800">Move History</h3>
      </div>
      <div className="flex-1 p-4 overflow-y-auto">
        <div className="space-y-1">
          {(!moves || moves.length === 0) ? (
            <div className="text-gray-500 text-sm text-center py-4">
              No moves yet
            </div>
          ) : (
            moves.map((move, index) => {
              const whiteMoveIndex = index * 2 + 1;
              const blackMoveIndex = index * 2 + 2;
              return (
                <div
                  key={index}
                  className="grid grid-cols-3 text-sm p-1 rounded"
                >
                  <div className="text-gray-500 font-mono">{move.moveNumber}.</div>
                  <div
                    className={`font-mono text-gray-800 rounded mr-2 ${indexMoveHistory !== 0 && indexMoveHistory === whiteMoveIndex
                      ? 'bg-blue-400 font-bold shadow border border-yellow-600'
                      : ''
                      }`}
                  >
                    {move.white}
                  </div>
                  <div
                    className={`font-mono text-gray-800 rounded ${indexMoveHistory !== 0 && indexMoveHistory === blackMoveIndex
                      ? 'bg-blue-400 font-bold shadow border border-yellow-600'
                      : ''
                      }`}
                  >
                    {move.black || ''}
                  </div>
                </div>
              );
            })
          )}
        </div>
      </div>
      <div className="flex justify-center w-full">
        <div className="bg-white border border-gray-300 shadow px-4 py-2 mt-2 w-full flex justify-center">
          <ReplayControls
            onStart={handleReplayStart}
            onPrevious={handleReplayPrevious}
            onNext={handleReplayNext}
            onEnd={handleReplayEnd}
          />
        </div>
      </div>
    </div>
  );
};

const ConfirmTooltip: React.FC<{
  isOpen: boolean;
  title: string;
  message: string;
  onConfirm: () => void;
  onCancel: () => void;
  confirmText?: string;
  cancelText?: string;
  confirmButtonClass?: string;
  position?: 'top' | 'bottom' | 'left' | 'right';
}> = ({
  isOpen,
  title,
  message,
  onConfirm,
  onCancel,
  confirmText = "Xác nhận",
  cancelText = "Hủy",
  confirmButtonClass = "bg-red-500 hover:bg-red-600",
  position = 'right'
}) => {
    if (!isOpen) return null;

    const positionClasses = {
      top: 'bottom-full left-1/2 transform -translate-x-1/2 mb-2',
      bottom: 'top-full left-1/2 transform -translate-x-1/2 mt-2',
      left: 'right-full top-1/2 transform -translate-y-1/2 mr-2',
      right: 'left-full top-1/2 transform -translate-y-1/2 ml-2'
    };

    const arrowClasses = {
      top: 'top-full left-1/2 transform -translate-x-1/2 border-l-transparent border-r-transparent border-b-transparent border-t-white',
      bottom: 'bottom-full left-1/2 transform -translate-x-1/2 border-l-transparent border-r-transparent border-t-transparent border-b-white',
      left: 'left-full top-1/2 transform -translate-y-1/2 border-t-transparent border-b-transparent border-r-transparent border-l-white',
      right: 'right-full top-1/2 transform -translate-y-1/2 border-t-transparent border-b-transparent border-l-transparent border-r-white'
    };

    return (
      <div className={`absolute ${positionClasses[position]} z-50`}>
        <div
          className={`absolute ${arrowClasses[position]} w-0 h-0 border-8`}
          style={{ zIndex: 51 }}
        />

        <div className="bg-white rounded-lg shadow-xl border border-gray-200 p-4 min-w-[280px] max-w-[320px]">
          <h4 className="text-sm font-semibold text-gray-900 mb-2">{title}</h4>
          <p className="text-xs text-gray-600 mb-4">{message}</p>
          <div className="flex gap-2 justify-end">
            <button
              onClick={onCancel}
              className="px-3 py-1.5 bg-gray-200 hover:bg-gray-300 text-gray-700 text-xs font-medium rounded transition-colors"
            >
              {cancelText}
            </button>
            <button
              onClick={onConfirm}
              className={`px-3 py-1.5 ${confirmButtonClass} text-white text-xs font-medium rounded transition-colors`}
            >
              {confirmText}
            </button>
          </div>
        </div>
      </div>
    );
  };

const DrawOfferDialog: React.FC<{
  isOpen: boolean;
  offerFromUsername: string;
  onAccept: () => void;
  onDecline: () => void;
  position?: 'top' | 'bottom' | 'left' | 'right';
}> = ({ isOpen, offerFromUsername, onAccept, onDecline, position = 'right' }) => {
  if (!isOpen) return null;

  const positionClasses = {
    top: 'bottom-full left-1/2 transform -translate-x-1/2 mb-2',
    bottom: 'top-full left-1/2 transform -translate-x-1/2 mt-2',
    left: 'right-full top-1/2 transform -translate-y-1/2 mr-2',
    right: 'left-full top-1/2 transform -translate-y-1/2 ml-2'
  };

  const arrowClasses = {
    top: 'top-full left-1/2 transform -translate-x-1/2 border-l-transparent border-r-transparent border-b-transparent border-t-white',
    bottom: 'bottom-full left-1/2 transform -translate-x-1/2 border-l-transparent border-r-transparent border-t-transparent border-b-white',
    left: 'left-full top-1/2 transform -translate-y-1/2 border-t-transparent border-b-transparent border-r-transparent border-l-white',
    right: 'right-full top-1/2 transform -translate-y-1/2 border-t-transparent border-b-transparent border-l-transparent border-r-white'
  };

  return (
    <div className={`absolute ${positionClasses[position]} z-50`}>
      {/* Arrow */}
      <div
        className={`absolute ${arrowClasses[position]} w-0 h-0 border-8`}
        style={{ zIndex: 51 }}
      />

      {/* Modal content */}
      <div className="bg-white rounded-lg shadow-xl border border-gray-200 p-4 min-w-[320px] max-w-[380px]">
        <h3 className="text-base font-bold text-gray-900 mb-3">Đề nghị hòa</h3>
        <p className="text-sm text-gray-600 mb-4">
          <span className="font-semibold text-blue-600">{offerFromUsername}</span> đã đề nghị hòa.
          Bạn có muốn chấp nhận không?
        </p>
        <div className="flex gap-2 justify-end">
          <button
            onClick={onDecline}
            className="px-3 py-1.5 bg-gray-200 hover:bg-gray-300 text-gray-700 text-sm font-medium rounded transition-colors"
          >
            Từ chối
          </button>
          <button
            onClick={onAccept}
            className="px-3 py-1.5 bg-blue-500 hover:bg-blue-600 text-white text-sm font-medium rounded transition-colors"
          >
            Chấp nhận
          </button>
        </div>
      </div>
    </div>
  );
};

const GamePage: React.FC = () => {
  const params = useParams();
  const gameId = params?.gameId as string;
  const { user } = useAuth();

  // States
  const [players, setPlayers] = React.useState<PlayerData[]>([]);
  const [moveHistory, setMoveHistory] = React.useState<MoveHistoryItem[]>([]);
  const [gameState, setGameState] = React.useState<ClientGameState | null>(null);
  const [gameResult, setGameResult] = React.useState<{ result: string; reason: string } | null>(null);
  const [showDialog, setShowDialog] = React.useState(false);
  const [gameStatus, setGameStatus] = React.useState<'loading' | 'playing' | 'end'>('loading');
  const [playerTimes, setPlayerTimes] = React.useState<{ white: number; black: number }>({ white: 0, black: 0 });
  const [error, setError] = React.useState<string | null>(null);
  const [bitboardStates, setBitboardStates] = React.useState<ClientGameState[]>([]);
  const [indexMoveHistory, setIndexMoveHistory] = React.useState<number>(0);
  const [wsConnected, setWsConnected] = React.useState(false);

  const [showResignConfirm, setShowResignConfirm] = React.useState(false);
  const [showDrawConfirm, setShowDrawConfirm] = React.useState(false);
  const [pendingDrawOffer, setPendingDrawOffer] = React.useState<{ offerId: string; offerFrom: number } | null>(null);
  const [showDrawOfferDialog, setShowDrawOfferDialog] = React.useState(false);

  const ws = React.useRef<WebSocket | null>(null);

  // WebSocket connection
  React.useEffect(() => {
    if (!gameId || !user?.id) return;

    const wsUrl = process.env.NEXT_PUBLIC_WEBSOCKET_URL || 'ws://localhost:3005/ws';
    ws.current = new WebSocket(wsUrl);

    ws.current.onopen = () => {
      setWsConnected(true);
      setError(null);

      const joinMessage = {
        type: 'joinRoom',
        roomId: gameId,
        userId: user.id,
        username: user.username || `User${user.id}`
      };

      if (ws.current && ws.current.readyState === WebSocket.OPEN) {
        ws.current.send(JSON.stringify(joinMessage));
      }
    };

    ws.current.onmessage = (event) => {
      try {
        const message: WSStateUpdateMessage = JSON.parse(event.data);
        switch (message.type) {
          case 'gameState':
            handleGameStateMessage(message);
            break;

          case 'gameUpdate':
            handleGameUpdateMessage(message);
            break;

          case 'error':
            setError(message.error || 'Unknown error');
            break;

          case 'gameEnd':
            handleGameEndMessage(message);
            console.log("mes: ", message);
            break;

          case 'drawOffer':
            handleDrawOfferMessage(message);
            break;

          case 'drawDeclined':
            handleDrawDeclinedMessage(message);
            break;
          default:
            console.log('Unknown message type:', message.type);
        }
      } catch (err) {
        console.error('Error parsing WebSocket message:', err);
      }
    };

    ws.current.onclose = (event) => {
      setWsConnected(false);
    };

    ws.current.onerror = (error) => {
      setError('Connection error. Please check your network.');
      setWsConnected(false);
    };

    return () => {
      if (ws.current && ws.current.readyState === WebSocket.OPEN) {
        ws.current.close();
      }
    };
  }, [gameId, user?.id]);

  // Handle message types
  const handleGameStateMessage = (message: WSStateUpdateMessage) => {
    if (message.gameState) {
      const clientGameState = convertWSGameStateToClient(message.gameState);
      setGameState(clientGameState);
      setBitboardStates([clientGameState]);
      setIndexMoveHistory(0);
    }

    if (message.player1 && message.player2) {
      setPlayers([message.player1, message.player2]);
    }

    if (message.whiteTimeLeft !== undefined && message.blackTimeLeft !== undefined) {
      setPlayerTimes({
        white: message.whiteTimeLeft,
        black: message.blackTimeLeft
      });
    }

    setGameStatus('playing');
  };

  const handleGameUpdateMessage = (message: WSStateUpdateMessage) => {
    if (message.gameState) {
      const clientGameState = convertWSGameStateToClient(message.gameState);
      setGameState(clientGameState);
      setMoveHistory(message.moveHistory ?? []);
      setBitboardStates(prev => [...prev, clientGameState]);
      setIndexMoveHistory(prev => prev + 1);
    }

    if (message.whiteTimeLeft !== undefined && message.blackTimeLeft !== undefined) {
      setPlayerTimes({
        white: message.whiteTimeLeft * 1000,
        black: message.blackTimeLeft * 1000
      });
    }
  };

  const handleGameEndMessage = (message: WSStateUpdateMessage) => {
    if (message.result && message.reason) {
      setGameResult({
        result: message.result,
        reason: message.reason
      });
      setShowDialog(true);
      setGameStatus('end');
    }
  };

  const handleDrawOfferMessage = (message: WSStateUpdateMessage) => {
    if (message.offerId && message.offerFrom) {
      setPendingDrawOffer({
        offerId: message.offerId,
        offerFrom: message.offerFrom
      });
      setShowDrawOfferDialog(true);
    }
  };

  const handleDrawDeclinedMessage = (message: WSStateUpdateMessage) => {
    console.log('Draw offer was declined');
  };

  const handleMove = (from: Position, to: Position, promotion?: string) => {
    if (!wsConnected || !ws.current) {
      console.error('Cannot send move: WebSocket not connected');
      return;
    }

    const moveMessage = {
      type: 'move',
      fromRow: from.row,
      fromCol: from.col,
      toRow: to.row,
      toCol: to.col,
      ...(promotion && { promotion })
    };
    ws.current.send(JSON.stringify(moveMessage));
  };

  const handleResign = () => {
    if (!wsConnected || !ws.current) {
      console.error(' Cannot resign: WebSocket not connected');
      return;
    }
    setShowResignConfirm(true);
  };

const confirmResign = () => {
    const resignMessage = {
      type: 'gameAction', 
      action: 'resign'    
    };

    ws.current?.send(JSON.stringify(resignMessage));
    setShowResignConfirm(false);
};

  const handleDrawOffer = () => {
    if (!wsConnected || !ws.current) {
      console.error(' Cannot offer draw: WebSocket not connected');
      return;
    }
    setShowDrawConfirm(true);
  };

const confirmDrawOffer = () => {
    const drawMessage = {
      type: 'gameAction',
      action: 'drawOffer'  
    };

    ws.current?.send(JSON.stringify(drawMessage));
    setShowDrawConfirm(false);
};

const acceptDrawOffer = () => {
    if (!pendingDrawOffer) return;

    const acceptMessage = {
      type: 'gameAction',
      action: 'drawAccept',
      offerId: pendingDrawOffer.offerId  
    };

    ws.current?.send(JSON.stringify(acceptMessage));
    setShowDrawOfferDialog(false);
    setPendingDrawOffer(null);
};

  const declineDrawOffer = () => {
    if (!pendingDrawOffer) return;

    const declineMessage = {
      type: 'gameAction',
      roomId: gameId,
      playerId: user?.id,
      action: 'drawDecline',
      offerId: pendingDrawOffer.offerId
    };

    ws.current?.send(JSON.stringify(declineMessage));
    setShowDrawOfferDialog(false);
    setPendingDrawOffer(null);
  };

  // Get player info
  const myPlayer = user && players.length > 0
    ? players.find((p: PlayerData) => p.userId === user.id)
    : null;
  const opponentPlayer = user && players.length > 0
    ? players.find((p: PlayerData) => p.userId !== user.id)
    : null;

  // Get username of player who made draw offer
  const drawOfferPlayerUsername = pendingDrawOffer
    ? players.find(p => p.userId === pendingDrawOffer.offerFrom)?.username || 'Unknown'
    : '';

  // Determine player's color
  const myColor = myPlayer?.color || 'white';

  // Helper format time for Clock
  const formatTimeForClock = (timeInMs: number) => {
    const minutes = Math.floor(timeInMs / 60);
    const seconds = Math.floor(timeInMs % 60);
    return { minutes, seconds };
  };

  // Loading state
  if (gameStatus === 'loading') {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <div className="text-gray-600">Loading game...</div>
          {!wsConnected && (
            <div className="text-red-600 mt-2">Connecting to game server...</div>
          )}
        </div>
      </div>
    );
  }

  // Error state
  if (error) {
    return (
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="text-center p-8 bg-white rounded-lg shadow-lg">
          <div className="text-red-600 text-xl mb-4"> Error</div>
          <div className="text-gray-700 mb-4">{error}</div>
          <button
            onClick={() => window.location.reload()}
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100 p-4">
      <div className="max-w-7xl mx-auto">

        <div className="grid grid-cols-12 gap-6">
          {/* Left sidebar */}
          <div className="col-span-3">
            <GameInfo
              gameType="RAPID"
              timeControl="10+0"
            />

            <div className="bg-white border border-gray-300 rounded-lg p-3 mt-4">
              {opponentPlayer && (
                <div className="mb-4">
                  <PlayerInfo
                    username={opponentPlayer.username}
                    elo={opponentPlayer.rating || 1200}
                    isConnected={opponentPlayer.isOnline}
                    isCurrentPlayer={gameState?.activeColor === opponentPlayer.color}
                  />
                </div>
              )}
              {myPlayer && (
                <div>
                  <PlayerInfo
                    username={myPlayer.username}
                    elo={myPlayer.rating || 1200}
                    isConnected={myPlayer.isOnline}
                    isCurrentPlayer={gameState?.activeColor === myPlayer.color}
                  />
                </div>
              )}
            </div>
            {gameStatus === 'playing' && (
              <div className="bg-white border border-gray-300 rounded-lg p-3 mt-4 relative">
                <div className="flex flex-col gap-2">
                  {/* Draw button container */}
                  <div className="relative">
                    <button
                      onClick={handleDrawOffer}
                      disabled={!wsConnected}
                      className="w-full px-4 py-2 bg-cyan-500 hover:bg-blue-600 text-white font-medium rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                    >
                      Đề nghị hòa
                    </button>

                    {/* Draw confirm tooltip */}
                    <ConfirmTooltip
                      isOpen={showDrawConfirm}
                      title="Đề nghị hòa"
                      message="Bạn có muốn gửi đề nghị hòa cho đối thủ?"
                      onConfirm={confirmDrawOffer}
                      onCancel={() => setShowDrawConfirm(false)}
                      confirmText="Gửi đề nghị"
                      cancelText="Hủy bỏ"
                      confirmButtonClass="bg-blue-500 hover:bg-blue-600"
                      position="right"
                    />
                  </div>

                  {/* Resign button container */}
                  <div className="relative">
                    <button
                      onClick={handleResign}
                      disabled={!wsConnected}
                      className="w-full px-4 py-2 bg-red-500 hover:bg-red-600 text-white font-medium rounded-lg transition-colors disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                    >
                      Đầu hàng
                    </button>

                    {/* Resign confirm tooltip */}
                    <ConfirmTooltip
                      isOpen={showResignConfirm}
                      title="Xác nhận đầu hàng"
                      message="Bạn có chắc chắn muốn đầu hàng? Hành động này không thể hoàn tác."
                      onConfirm={confirmResign}
                      onCancel={() => setShowResignConfirm(false)}
                      confirmText="Đầu hàng"
                      cancelText="Hủy bỏ"
                      confirmButtonClass="bg-red-500 hover:bg-red-600"
                      position="right"
                    />

                    <DrawOfferDialog
                      isOpen={showDrawOfferDialog}
                      offerFromUsername={drawOfferPlayerUsername}
                      onAccept={acceptDrawOffer}
                      onDecline={declineDrawOffer}
                    />
                  </div>
                </div>
              </div>
            )}
          </div>

          {/* Game board */}
          <div className="col-span-6">
            <div className="flex flex-col items-center justify-center min-h-[600px]">
              {gameResult && (
                <GameResultDialog
                  isOpen={showDialog}
                  gameResult={gameResult}
                  onRematch={() => window.location.reload()}
                  onNewGame={() => window.location.href = '/play/online'}
                  onClose={() => setShowDialog(false)}
                />
              )}

              {/* Draw Offer Dialog */}

              <div className="relative">
                {gameState && (
                  <Board
                    playerColor={myColor}
                    gameState={gameState}
                    onMove={handleMove}

                  />
                )}
              </div>
            </div>
          </div>

          {/* Right sidebar */}
          <div className="col-span-3 flex flex-col justify-center">
            <div className="space-y-4">
              {/* Opponent timer */}
              {opponentPlayer && (
                <div className="bg-white border border-gray-300 rounded-lg p-3">
                  <Clock
                    initialMinutes={formatTimeForClock(playerTimes[opponentPlayer.color as 'white' | 'black']).minutes}
                    initialSeconds={formatTimeForClock(playerTimes[opponentPlayer.color as 'white' | 'black']).seconds}
                    increment={0}
                    isActive={gameStatus === 'playing' && gameState?.activeColor === opponentPlayer.color}
                    onTimeUp={() => console.log(`${opponentPlayer.color} time up!`)}
                  />
                </div>
              )}

              {/* Move history */}
              <MoveHistory
                moves={moveHistory}
                indexMoveHistory={indexMoveHistory}
                setIndexMoveHistory={setIndexMoveHistory}
                bitboardStates={bitboardStates}
                setGameState={setGameState}
              />

              {/* My timer */}
              {myPlayer && (
                <div className="bg-white border border-gray-300 rounded-lg p-3">
                  <Clock
                    initialMinutes={formatTimeForClock(playerTimes[myPlayer.color as 'white' | 'black']).minutes}
                    initialSeconds={formatTimeForClock(playerTimes[myPlayer.color as 'white' | 'black']).seconds}
                    increment={0}
                    isActive={gameStatus === 'playing' && gameState?.activeColor === myPlayer.color}
                    onTimeUp={() => console.log(`${myPlayer.color} time up!`)}
                  />
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default GamePage;