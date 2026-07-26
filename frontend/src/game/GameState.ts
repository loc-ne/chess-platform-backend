import { Position, Piece } from './types';
import { BitboardGame, ChessEngine } from './Board';
import { FenUtils, STARTING_FEN } from './Fen'

export interface CastlingRights {
  whiteKingSide: boolean;
  whiteQueenSide: boolean;
  blackKingSide: boolean;
  blackQueenSide: boolean;
}

export interface MoveRecord {
  from: Position;
  to: Position;
  piece: Piece;
  capturedPiece?: Piece;
  moveNumber: number;
  fen: string;
  timestamp: number;
  timeLeft?: {
    white: number;
    black: number;
  };
}

export interface ClientGameState {
  currentFen: string;
  bitboards: BitboardGame;
  activeColor: 'white' | 'black';
  castlingRights: CastlingRights;
  enPassantSquare: Position | null;
}


export interface TimeControl {
  type: 'rapid' | 'blitz' | 'bullet';
  initialTime: number;
  increment: number;
}

export const GameStateManager = {
  convertFenToClientState(fen: string): ClientGameState {
    const components = FenUtils.parseFEN(fen);

    return {
      currentFen: fen,
      bitboards: FenUtils.fenPiecesToBitboards(components.pieces),
      activeColor: components.activeColor === 'w' ? 'white' : 'black',
      castlingRights: FenUtils.parseCastlingRights(components.castlingRights),
      enPassantSquare: FenUtils.parseEnPassantSquare(components.enPassantSquare)
    };
  },

  getPieceTypeFromSan(san: string): string {
    if (/^[a-h][1-8]$/.test(san) || /^[a-h]x[a-h][1-8]$/.test(san)) return 'pawn';
    if (san.startsWith('N')) return 'knight';
    if (san.startsWith('B')) return 'bishop';
    if (san.startsWith('R')) return 'rook';
    if (san.startsWith('Q')) return 'queen';
    if (san.startsWith('K')) return 'king';
    return 'pawn';
  },

  algebraicToPosition(move: string): Position | null {
    if (!move || move.length < 2) return null;
    const file = move[0].toLowerCase().charCodeAt(0) - 97;
    const rank = parseInt(move[1], 10);
    if (file < 0 || file > 7 || rank < 1 || rank > 8) return null;
    return { row: rank - 1, col: file };
  },


  getTargetPositionFromSan(san: string): Position | null {
    const match = san.match(/^(?:(O-O(?:-O)?)|(?:[KQRBN]?[a-h]?[1-8]?x?([a-h][1-8])(?:=[QRBN])?[+#]?))$/i);

    if (!match) return null;

    if (match[1]) {
      return null;
    }

    const target = match[2];
    if (!target) return null;

    return this.algebraicToPosition(target);
  },

  getBitboardStatesFromMoves(moves: string[]): ClientGameState[] {
    const game = ChessEngine.createBitboardGame();
    const states: ClientGameState[] = [];

    let currentState: ClientGameState = {
      currentFen: STARTING_FEN,
      bitboards: game,
      activeColor: 'white',
      castlingRights: {
        whiteKingSide: true,
        whiteQueenSide: true,
        blackKingSide: true,
        blackQueenSide: true,
      },
      enPassantSquare: null,
    };
    states.push({
  ...currentState,
  bitboards: ChessEngine.cloneBitboards(currentState.bitboards)
});

    for (let i = 0; i < moves.length; i++) {
      const moveStr = moves[i];
      const pieceType = this.getPieceTypeFromSan(moveStr);
      const targetPos = this.getTargetPositionFromSan(moveStr);

      if (!targetPos) continue;

      let found = false;
      for (let row = 0; row < 8 && !found; row++) {
        for (let col = 0; col < 8 && !found; col++) {
          const fromPos: Position = { row, col };
          const piece = ChessEngine.getPieceAt(game, fromPos);
          if (!piece || piece.color !== currentState.activeColor || piece.type !== pieceType) continue;

          const movesArr = ChessEngine.generateMovesForPiece(game, currentState, fromPos);
          for (const toPos of movesArr) {
            if (toPos.row === targetPos.row && toPos.col === targetPos.col) {
              ChessEngine.makeMove(game, fromPos, toPos);
              currentState.activeColor = currentState.activeColor === 'white' ? 'black' : 'white';
              found = true;
              break;
            }
          }
        }
      }

      const newFen = FenUtils.convertBitboardToFen(
        game,
        currentState.activeColor === 'white' ? 'w' : 'b',
        FenUtils.castlingRightsToString(currentState.castlingRights),
        currentState.enPassantSquare
          ? FenUtils.positionToAlgebraic(currentState.enPassantSquare)
          : '-',
        0,
        i + 1
      );
      currentState = this.convertFenToClientState(newFen);
      states.push({
  ...currentState,
  bitboards: ChessEngine.cloneBitboards(currentState.bitboards)
});
    }
    return states;
  }
};



