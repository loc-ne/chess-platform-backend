export type PieceType = 'pawn' | 'knight' | 'bishop' | 'rook' | 'queen' | 'king';
export type PieceColor = 'white' | 'black';

export interface Position {
  row: number;  // 0-7
  col: number;  // 0-7
}

export interface Piece {
  type: PieceType;
  color: PieceColor;
}



