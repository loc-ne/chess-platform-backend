'use client';

import React, { useState, useEffect } from 'react';
import { useParams } from 'next/navigation';
import Board from '@/components/Board';
import PlayerInfo from '@/components/PlayerInfo';
import EvaluationBar from '@/components/EvaluationBar';
import GameAnalysis from '@/components/GameAnalysis';

import { ClientGameState, GameStateManager } from '@/game/GameState';
import { BitboardGame } from '@/game/Board';
import { Position } from '@/game/types';

const GamePage: React.FC = () => {
  const [gameState, setGameState] = useState<ClientGameState>();
  const [evaluation, setEvaluation] = useState(0.0); // Current position evaluation

  const myColor: 'white' | 'black' = 'white'; // Determine from game data

  // Sample data - replace with real data
  const player1 = {
    username: 'Vishwajeet17',
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
    username: 'Noop114',
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

  const handleMove = (from: Position, to: Position) => {
    // Handle move and update evaluation
    // setEvaluation(newEval);
  };

  interface GamePlayer {
    userId: number;
    username: string;
    elo: number;
  }

  interface GameData {
    players: {
      white: GamePlayer;
      black: GamePlayer;
    };
    moves: string[];
    result: string;
    createdAt: string;
    timeControl: string;
    gameType: string;
    winnerId: string;
    whiteTimeLeft: number;
    blackTimeLeft: number;
    reason: string;
    lastFen: string;
  }
  const params = useParams();
  const gameId = params?.gameId as string;

  const [gameData, setGameData] = useState<GameData | null>(null);
  const [currentMoveIndex, setCurrentMoveIndex] = useState(0);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchGameData = async () => {
      if (!gameId) return;

      try {
        setLoading(true);
        const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/games/${gameId}`, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
          },
        });
        
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json(); 
        setGameData(data);
        
        
        setError(null);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to load game');
        console.error('Error fetching game data:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchGameData();
  }, [gameId]);

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
              <div className="flex justify-center">
                <Board
                  playerColor={myColor}
                  gameState={gameState}
                  onMove={handleMove}
                />
              </div>
            </div>
          </div>

          {/* Analysis Panel */}
          <div className="lg:col-span-3">
            <GameAnalysis
              player1={player1}
              player2={player2}
              winner="player2"
            />
          </div>

        </div>
      </div>
    </div>
  );
};

export default GamePage;