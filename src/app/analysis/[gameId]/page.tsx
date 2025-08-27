'use client';

import React, { useState, useEffect, useRef } from 'react';
import { useParams } from 'next/navigation';
import Board from '@/components/Board';
import PlayerInfo from '@/components/PlayerInfo';
import EvaluationBar from '@/components/EvaluationBar';
import GameAnalysis from '@/components/GameAnalysis';
import ReplayControls from '@/components/ReplayControls';

import { ClientGameState, GameStateManager } from '@/game/GameState';
import { BitboardGame, ChessEngine } from '@/game/Board';
import { Position } from '@/game/types';
// Import MoveLabel type from Board component
import type { MoveLabel } from '@/components/Board';

const STARTING_FEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1";

interface MoveData {
  move_number: number;
  player: 'white' | 'black';
  move: string;
  from: string;
  to: string;
  label: 'brilliant' | 'best' | 'excellent' | 'good' | 'book' | 'inaccuracy' | 'mistake' | 'blunder';
  evaluation: number;
}

interface AnalysisData {
  opening_name: string;
  moves: MoveData[];
}

interface ApiResponse {
  status: string;
  message: string;
  data: AnalysisData;
}

interface MoveHistoryItem {
  moveNumber: number;
  white: string;
  black?: string;
  whiteLabel?: string;
  blackLabel?: string;
}

const GamePage: React.FC = () => {
  const [gameState, setGameState] = useState<ClientGameState>();
  const [evaluation, setEvaluation] = useState(0.0);
  const [analysisData, setAnalysisData] = useState<AnalysisData | null>(null);
  const [currentMoveIndex, setCurrentMoveIndex] = useState(0);
  const [pgn, setPgn] = useState('');
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [showPgnForm, setShowPgnForm] = useState(true);
  const [showReplayPanel, setShowReplayPanel] = useState(false);
  const [isReplaying, setIsReplaying] = useState(false);

  // Ref for auto-scroll functionality
  const moveHistoryRef = useRef<HTMLDivElement>(null);
  const activeMoveRef = useRef<HTMLButtonElement>(null);

  // Move labels mapping cho UI
  const moveLabelColors = {
    brilliant: 'bg-cyan-500 text-white',
    best: 'bg-green-500 text-white',
    excellent: 'bg-green-400 text-white',
    good: 'bg-green-600 text-white',
    book: 'bg-orange-500 text-white',
    inaccuracy: 'bg-yellow-500 text-black',
    mistake: 'bg-orange-600 text-white',
    blunder: 'bg-red-500 text-white'
  };

  const myColor: 'white' | 'black' = 'white';

  // Sample data for players
  const player1 = {
    username: 'White',
    accuracy: 74.2,
    brilliant: 0,
    great: 1,
    best: 5,
    excellent: 7,
    good: 1,
    inaccuracy: 6,
    mistake: 1,
    blunder: 2,
    missedWin: 0
  };

  const player2 = {
    username: 'Black',
    accuracy: 91.2,
    brilliant: 0,
    great: 0,
    best: 8,
    excellent: 8,
    good: 0,
    inaccuracy: 6,
    mistake: 1,
    blunder: 1,
    missedWin: 0
  };

  // Auto-scroll to active move
  useEffect(() => {
    if (showReplayPanel && activeMoveRef.current && moveHistoryRef.current) {
      const activeElement = activeMoveRef.current;
      const container = moveHistoryRef.current;

      // Calculate scroll position to center the active move
      const containerRect = container.getBoundingClientRect();
      const elementRect = activeElement.getBoundingClientRect();

      const scrollTop = container.scrollTop;
      const elementTop = elementRect.top - containerRect.top + scrollTop;
      const elementCenter = elementTop - (container.clientHeight / 2) + (elementRect.height / 2);

      container.scrollTo({
        top: Math.max(0, elementCenter),
        behavior: 'smooth'
      });
    }
  }, [currentMoveIndex, showReplayPanel]);

  // Initialize game state
  useEffect(() => {
    const initialState = GameStateManager.convertFenToClientState(STARTING_FEN);
    setGameState(initialState);
  }, []);

  // Submit PGN for analysis
  const handlePgnSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!pgn.trim()) return;

    setIsAnalyzing(true);
    try {
      const response = await fetch('http://localhost:8080/api/v1/analysis/pgn', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ pgn: pgn.trim() })
      });

      if (!response.ok) {
        throw new Error('Failed to analyze PGN');
      }

      const result: ApiResponse = await response.json();
      setAnalysisData(result.data);
      console.log(result.data);
      setCurrentMoveIndex(0);
      setShowPgnForm(false);

      // Reset to initial position
      const initialState = GameStateManager.convertFenToClientState(STARTING_FEN);
      setGameState(initialState);
      setEvaluation(result.data.moves[0]?.evaluation || 0);

    } catch (error) {
      console.error('Error analyzing PGN:', error);
      alert('Error analyzing PGN. Please check your input.');
    } finally {
      setIsAnalyzing(false);
    }
  };

  // Convert algebraic notation to Position
  const algebraicToPosition = (algebraic: string): Position => {
    const col = algebraic.charCodeAt(0) - 97; // 'a' = 0, 'b' = 1, ..., 'h' = 7
    const row = parseInt(algebraic[1]) - 1;    // '1' = 0, '2' = 1, ..., '8' = 7
    return { row, col };
  };

  // Apply moves up to current index
  const applyMovesToPosition = (moveIndex: number) => {
    if (!analysisData || !gameState) return;

    let currentState = GameStateManager.convertFenToClientState(STARTING_FEN);

    for (let i = 0; i < moveIndex; i++) {
      const moveData = analysisData.moves[i];
      if (moveData) {
        const from = algebraicToPosition(moveData.from);
        const to = algebraicToPosition(moveData.to);

        // Update game state using existing function
        currentState = ChessEngine.updateGameState(currentState, from, to);
      }
    }

    setGameState(currentState);

    // Update evaluation
    const currentMove = analysisData.moves[moveIndex - 1];
    if (currentMove) {
      setEvaluation(currentMove.evaluation);
    }
  };

  // ========== GET CURRENT HIGHLIGHTED MOVE ==========
  const getCurrentHighlightedMove = () => {
    if (!analysisData || currentMoveIndex === 0) return undefined;

    const currentMove = analysisData.moves[currentMoveIndex - 1];
    if (!currentMove) return undefined;

    return {
      from: algebraicToPosition(currentMove.from),
      to: algebraicToPosition(currentMove.to),
      label: currentMove.label as MoveLabel
    };
  };

  // Convert moves data to move history format
  const getMoveHistory = (): MoveHistoryItem[] => {
    if (!analysisData) return [];

    const history: MoveHistoryItem[] = [];

    for (let i = 0; i < analysisData.moves.length; i += 2) {
      const whiteMove = analysisData.moves[i];
      const blackMove = analysisData.moves[i + 1];

      if (whiteMove) {
        const moveItem: MoveHistoryItem = {
          moveNumber: whiteMove.move_number,
          white: whiteMove.move,
          whiteLabel: whiteMove.label,
        };

        if (blackMove) {
          moveItem.black = blackMove.move;
          moveItem.blackLabel = blackMove.label;
        }

        history.push(moveItem);
      }
    }

    return history;
  };

  // Handle move click in history
  const handleMoveClick = (moveIndex: number) => {
    setCurrentMoveIndex(moveIndex);
    applyMovesToPosition(moveIndex);
  };

  // Replay control handlers
  const handleReplayStart = () => {
    setCurrentMoveIndex(0);
    const initialState = GameStateManager.convertFenToClientState(STARTING_FEN);
    setGameState(initialState);
    setEvaluation(analysisData?.moves[0]?.evaluation || 0);
  };

  const handleReplayPrevious = () => {
    if (currentMoveIndex > 0) {
      const newIndex = currentMoveIndex - 1;
      setCurrentMoveIndex(newIndex);
      applyMovesToPosition(newIndex);
    }
  };

  const handleReplayNext = () => {
    if (analysisData && currentMoveIndex < analysisData.moves.length) {
      const newIndex = currentMoveIndex + 1;
      setCurrentMoveIndex(newIndex);
      applyMovesToPosition(newIndex);
    }
  };

  const handleReplayEnd = () => {
    if (analysisData) {
      const endIndex = analysisData.moves.length;
      setCurrentMoveIndex(endIndex);
      applyMovesToPosition(endIndex);
    }
  };


const timeoutIdsRef = useRef<number[]>([]);

const handleAutoReplay = () => {
  if (isReplaying) {
    timeoutIdsRef.current.forEach((id) => clearTimeout(id));
    timeoutIdsRef.current = [];
    setIsReplaying(false);
  } else {
    setIsReplaying(true);
    if (analysisData) {
      const newTimeouts: number[] = [];

      for (let i = currentMoveIndex + 1; i < analysisData.moves.length + 1; i++) {
        const timeoutId = window.setTimeout(() => {
          setCurrentMoveIndex(i);
          applyMovesToPosition(i);
          if (i === analysisData.moves.length - 1) {
            setIsReplaying(false);
            timeoutIdsRef.current = [];
          }
        }, (i - currentMoveIndex) * 1000);

        newTimeouts.push(timeoutId);
      }

      timeoutIdsRef.current = newTimeouts;
    }
  }
};



  const handleMove = (from: Position, to: Position) => {
    // For analysis mode, moves are controlled by replay
    // You can implement manual move functionality here if needed
  };

  // Handle show replay panel
  const handleShowReplay = () => {
    setShowReplayPanel(true);
  };

  // Handle hide replay panel
  const handleHideReplay = () => {
    setShowReplayPanel(false);
  };

  // Get current move label for display
  const getCurrentMoveLabel = (): string | null => {
    if (!analysisData || currentMoveIndex === 0) return null;
    const currentMove = analysisData.moves[currentMoveIndex - 1];
    return currentMove?.label || null;
  };

  if (showPgnForm) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="bg-white rounded-xl shadow-lg border border-gray-200 p-8 max-w-2xl w-full mx-4">
          <h1 className="text-2xl font-bold text-gray-800 mb-6 text-center">
            Phân tích PGN
          </h1>

          <form onSubmit={handlePgnSubmit} className="space-y-4">
            <div>
              <label htmlFor="pgn" className="block text-sm font-medium text-gray-700 mb-2">
                Nhập PGN của trận đấu:
              </label>
              <textarea
                id="pgn"
                value={pgn}
                onChange={(e) => setPgn(e.target.value)}
                placeholder="1. d4 Nf6 2. c4 g6 3. Nc3 Bg7..."
                className="w-full h-32 p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 resize-none"
                disabled={isAnalyzing}
              />
            </div>

            <button
              type="submit"
              disabled={isAnalyzing || !pgn.trim()}
              className="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-semibold py-3 px-4 rounded-lg transition-colors"
            >
              {isAnalyzing ? (
                <span className="flex items-center justify-center">
                  <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Đang phân tích...
                </span>
              ) : (
                'Phân tích PGN'
              )}
            </button>
          </form>
        </div>
      </div>
    );
  }

  const moveHistory = getMoveHistory();

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-6">

          {/* Evaluation Bar */}
          <div className="lg:col-span-1">
            <EvaluationBar
              evaluation={evaluation}
              playerColor={myColor}
              className="h-96"
            />
          </div>

          {/* Main Board */}
          <div className="lg:col-span-8">
            <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
              <div className="flex justify-center relative">
                <Board
                  playerColor={myColor}
                  gameState={gameState}
                  onMove={handleMove}
                  highlightedMove={getCurrentHighlightedMove()}
                />
              </div>
            </div>
          </div>

          {/* Analysis Panel */}
          <div className="lg:col-span-3">
            {!showReplayPanel ? (
              // Show Game Analysis
              <GameAnalysis
                player1={player1}
                player2={player2}
                winner="draw"
                onShowReplay={handleShowReplay}
              />
            ) : (
              // Show Move History with Replay Controls
              analysisData && (
                <>
                  {/* Header with opening name - Above Move History */}
                  <div className="mb-4 bg-white rounded-xl shadow-sm border border-gray-200 p-4">
                    <div className="space-y-3">
                      <div>
                        <h1 className="text-lg font-semibold text-gray-800">
                          {analysisData.opening_name}
                        </h1>
                      </div>
                      <button
                        onClick={() => setShowPgnForm(true)}
                        className="w-full px-3 py-2 bg-gray-600 hover:bg-gray-700 text-white rounded-lg transition-colors text-sm"
                      >
                        Phân tích PGN
                      </button>
                    </div>
                  </div>

                  {/* Move History Panel */}
                  <div className="bg-white border border-gray-300 rounded-lg h-96 flex flex-col">
                    {/* Header with Back Button */}
                    <div className="bg-gray-100 px-4 py-2 border-b border-gray-300 rounded-t-lg flex items-center justify-between">
                      <h3 className="font-semibold text-gray-800">Move History</h3>
                      <button
                        onClick={handleHideReplay}
                        className="text-gray-600 hover:text-gray-800 text-sm"
                      >
                        ← Tổng quan
                      </button>
                    </div>

                    {/* Move History */}
                    <div
                      ref={moveHistoryRef}
                      className="flex-1 p-4 overflow-y-auto"
                    >
                      <div className="space-y-1">
                        {(!moveHistory || moveHistory.length === 0) ? (
                          <div className="text-gray-500 text-sm text-center py-4">
                            No moves yet
                          </div>
                        ) : (
                          moveHistory.map((move, index) => {
                            const whiteMoveIndex = index * 2 + 1;
                            const blackMoveIndex = index * 2 + 2;
                            const isWhiteActive = currentMoveIndex === whiteMoveIndex;
                            const isBlackActive = currentMoveIndex === blackMoveIndex;

                            return (
                              <div
                                key={index}
                                className="grid grid-cols-3 text-sm p-1 rounded gap-1"
                              >
                                {/* Move Number */}
                                <div className="text-gray-500 font-mono text-center">
                                  {move.moveNumber}.
                                </div>

                                {/* White Move */}
                                <button
                                  ref={isWhiteActive ? activeMoveRef : null}
                                  onClick={() => handleMoveClick(whiteMoveIndex)}
                                  className={`font-mono text-gray-800 rounded p-1 hover:bg-gray-100 text-left ${isWhiteActive
                                    ? 'bg-blue-400 font-bold shadow border border-yellow-600 text-white'
                                    : ''
                                    }`}
                                >
                                  <div className="flex items-center justify-between">
                                    <span>{move.white}</span>
                                  </div>
                                </button>

                                {/* Black Move */}
                                {move.black ? (
                                  <button
                                    ref={isBlackActive ? activeMoveRef : null}
                                    onClick={() => handleMoveClick(blackMoveIndex)}
                                    className={`font-mono text-gray-800 rounded p-1 hover:bg-gray-100 text-left ${isBlackActive
                                      ? 'bg-blue-400 font-bold shadow border border-yellow-600 text-white'
                                      : ''
                                      }`}
                                  >
                                    <div className="flex items-center justify-between">
                                      <span>{move.black}</span>
                                    </div>
                                  </button>
                                ) : (
                                  <div></div>
                                )}
                              </div>
                            );
                          })
                        )}
                      </div>
                    </div>

                    {/* Replay Controls */}
                    <div className="flex justify-center w-full">
                      <div className="bg-white border border-gray-300 shadow px-4 py-2 mt-2 w-full flex justify-center">
                        <ReplayControls
                          onStart={handleReplayStart}
                          onPrevious={handleReplayPrevious}
                          onNext={handleReplayNext}
                          onEnd={handleReplayEnd}
                          onAutoReplay={handleAutoReplay}
                          isReplaying={isReplaying}
                          currentMove={currentMoveIndex}
                          totalMoves={analysisData.moves.length}
                        />
                      </div>
                    </div>
                  </div>
                </>
              )
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default GamePage;