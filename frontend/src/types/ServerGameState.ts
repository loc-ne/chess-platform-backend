export interface MaterialCount {
  pawns: number;
  knights: number;
  bishops: number;
  rooks: number;
  queens: number;
}

export interface CastlingRights {
  whiteKingSide: boolean;
  whiteQueenSide: boolean;
  blackKingSide: boolean;
  blackQueenSide: boolean;
}

export interface Position {
  row: number;
  col: number;
}

export interface MoveNotation {
  moveNumber: number;
  white: string;
  black?: string;
}

export interface BitboardGame {
  WhitePawns: string;   
  WhiteRooks: string;
  WhiteKnights: string;
  WhiteBishops: string;
  WhiteQueens: string;
  WhiteKing: string;

  BlackPawns: string;
  BlackRooks: string;
  BlackKnights: string;
  BlackBishops: string;
  BlackQueens: string;
  BlackKing: string;
  [key: string]: string;
}

export interface ServerGameState {
  currentFen: string;
  bitboards: BitboardGame;
  activeColor: 'white' | 'black';
  castlingRights: CastlingRights;
  enPassantSquare: Position | null;
  moveHistory: MoveNotation[];
  fullMoveNumber: number;
  halfMoveClock: number;
  positionCounts: Record<string, number>;
  materialCount: Record<string, MaterialCount>;
}