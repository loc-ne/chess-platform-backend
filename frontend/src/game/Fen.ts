import { ClientGameState, CastlingRights} from './GameState';
import { BitboardGame, ChessEngine } from './Board';
import { Position, Piece } from './types';

export interface FENComponents {
  pieces: string;
  activeColor: string;
  castlingRights: string;
  enPassantSquare: string;
  halfmoveClock: string;
  fullmoveNumber: string;
}

export const STARTING_FEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1";

export const FenUtils = {
  parseFEN(fen: string): FENComponents {
      const parts = fen.trim().split(' ');
      if (parts.length !== 6) {
        throw new Error(`Invalid FEN: expected 6 parts, got ${parts.length}`);
      }
  
      return {
        pieces: parts[0],
        activeColor: parts[1],
        castlingRights: parts[2],
        enPassantSquare: parts[3],
        halfmoveClock: parts[4],
        fullmoveNumber: parts[5]
      };
    },
  
    validateFEN(fen: string): boolean {
      try {
        const components = this.parseFEN(fen);
  
        if (!['w', 'b'].includes(components.activeColor)) return false;
        if (!/^[KQkq\-]*$/.test(components.castlingRights)) return false;
        if (!/^([a-h][36]|\-)$/.test(components.enPassantSquare)) return false;
  
        const halfmove = parseInt(components.halfmoveClock);
        const fullmove = parseInt(components.fullmoveNumber);
        if (isNaN(halfmove) || isNaN(fullmove) || halfmove < 0 || fullmove < 1) return false;
  
        return true;
      } catch {
        return false;
      }
    },

    fenPiecesToBitboards(fenPieces: string): BitboardGame {
    const bitboards: BitboardGame = {
      whitePawns: BigInt(0), whiteRooks: BigInt(0), whiteKnights: BigInt(0),
      whiteBishops: BigInt(0), whiteQueens: BigInt(0), whiteKing: BigInt(0),
      blackPawns: BigInt(0), blackRooks: BigInt(0), blackKnights: BigInt(0),
      blackBishops: BigInt(0), blackQueens: BigInt(0), blackKing: BigInt(0)
    };

    const ranks = fenPieces.split('/');
    for (let row = 0; row < 8; row++) {
      const rank = ranks[7 - row];
      let col = 0;

      for (const char of rank) {
        if (char >= '1' && char <= '8') {
          col += parseInt(char);
        } else {
          const bitPosition = row * 8 + col;
          const bit = BigInt(1) << BigInt(bitPosition);

          switch (char) {
            case 'P': bitboards.whitePawns |= bit; break;
            case 'R': bitboards.whiteRooks |= bit; break;
            case 'N': bitboards.whiteKnights |= bit; break;
            case 'B': bitboards.whiteBishops |= bit; break;
            case 'Q': bitboards.whiteQueens |= bit; break;
            case 'K': bitboards.whiteKing |= bit; break;
            case 'p': bitboards.blackPawns |= bit; break;
            case 'r': bitboards.blackRooks |= bit; break;
            case 'n': bitboards.blackKnights |= bit; break;
            case 'b': bitboards.blackBishops |= bit; break;
            case 'q': bitboards.blackQueens |= bit; break;
            case 'k': bitboards.blackKing |= bit; break;
          }
          col++;
        }
      }
    }
    return bitboards;
  },

  parseCastlingRights(castlingString: string): CastlingRights {
    return {
      whiteKingSide: castlingString.includes('K'),
      whiteQueenSide: castlingString.includes('Q'),
      blackKingSide: castlingString.includes('k'),
      blackQueenSide: castlingString.includes('q')
    };
  },

  parseEnPassantSquare(enPassantString: string): Position | null {
    if (enPassantString === '-') return null;
    if (enPassantString.length !== 2) return null;

    const file = enPassantString.charCodeAt(0) - 'a'.charCodeAt(0);
    const rank = 8 - parseInt(enPassantString[1]);

    if (file >= 0 && file <= 7 && rank >= 0 && rank <= 7) {
      return { row: rank, col: file };
    }
    return null;
  },

  convertClientStateToFen(
    clientState: ClientGameState,
    halfmoveClock: number = 0,
    fullmoveNumber: number = 1
  ): string {
    const piecePlacement = this.bitboardsToPiecePlacement(clientState.bitboards);
    const activeColor = clientState.activeColor === 'white' ? 'w' : 'b';
    const castlingRights = this.castlingRightsToString(clientState.castlingRights);
    const enPassantSquare = clientState.enPassantSquare
      ? this.positionToAlgebraic(clientState.enPassantSquare)
      : '-';

    return [
      piecePlacement,
      activeColor,
      castlingRights,
      enPassantSquare,
      halfmoveClock.toString(),
      fullmoveNumber.toString()
    ].join(' ');
  },

  bitboardsToPiecePlacement(bitboards: BitboardGame): string {
    const board = ChessEngine.bitboardToBoard(bitboards);
    const ranks: string[] = [];

    for (let row = 0; row < 8; row++) {
      let rankString = '';
      let emptyCount = 0;

      for (let col = 0; col < 8; col++) {
        const piece = board[7 - row][col];

        if (piece === null) {
          emptyCount++;
        } else {
          if (emptyCount > 0) {
            rankString += emptyCount.toString();
            emptyCount = 0;
          }
          rankString += this.pieceToFenChar(piece);
        }
      }

      if (emptyCount > 0) {
        rankString += emptyCount.toString();
      }

      ranks.push(rankString);
    }

    return ranks.join('/');
  },

convertBitboardToFen(
  bitboards: BitboardGame,
  activeColor: 'w' | 'b' = 'w',
  castlingRights: string = 'KQkq',
  enPassant: string = '-',
  halfmoveClock: number = 0,
  fullmoveNumber: number = 1
): string {
  const piecePlacement = FenUtils.bitboardsToPiecePlacement(bitboards);
  return [
    piecePlacement,
    activeColor,
    castlingRights,
    enPassant,
    halfmoveClock.toString(),
    fullmoveNumber.toString()
  ].join(' ');
},

  pieceToFenChar(piece: Piece): string {
    const charMap: Record<string, string> = {
      'white-pawn': 'P', 'white-rook': 'R', 'white-knight': 'N',
      'white-bishop': 'B', 'white-queen': 'Q', 'white-king': 'K',
      'black-pawn': 'p', 'black-rook': 'r', 'black-knight': 'n',
      'black-bishop': 'b', 'black-queen': 'q', 'black-king': 'k'
    };

    return charMap[`${piece.color}-${piece.type}`] || '?';
  },


  castlingRightsToString(castlingRights: CastlingRights): string {
    let rights = '';
    if (castlingRights.whiteKingSide) rights += 'K';
    if (castlingRights.whiteQueenSide) rights += 'Q';
    if (castlingRights.blackKingSide) rights += 'k';
    if (castlingRights.blackQueenSide) rights += 'q';
    return rights || '-';
  },

  positionToAlgebraic(position: Position): string {
    const file = String.fromCharCode('a'.charCodeAt(0) + position.col);
    const rank = (8 - position.row).toString();
    return file + rank;
  },

  algebraicToPosition(algebraic: string): Position | null {
    if (algebraic.length !== 2) return null;

    const file = algebraic.charCodeAt(0) - 'a'.charCodeAt(0);
    const rank = parseInt(algebraic[1]);

    if (file < 0 || file > 7 || rank < 1 || rank > 8) return null;

    return {
      row: 8 - rank,  // Convert rank 1-8 to row 7-0
      col: file       // Convert file a-h to col 0-7
    };
  }
};

