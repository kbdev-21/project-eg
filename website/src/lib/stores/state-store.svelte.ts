export type UserState = "IDLE" | "QUEUING_CARO" | "PLAYING_CARO";

export type CaroPiece = "" | "X" | "O";

export type CaroMove = {
  piece: CaroPiece;
  x: number;
  y: number;
  playedAt: string;
};

// board[y][x] - backend index y trước (domain/caro.go)
export type CaroBoard = CaroPiece[][];

export type CaroMatch = {
  id: string;
  isRated: boolean;
  xPlayerId: string;
  xPlayerRating: number;
  oPlayerId: string;
  oPlayerRating: number;
  board: CaroBoard;
  moves: CaroMove[];
  turnOf: CaroPiece;
  winner: CaroPiece;
  isEnded: boolean;
  startedAt: string;
};

// db.CaroMatch được embed nên các field bị flatten ra cùng cấp.
// winnerId = null nghĩa là hoà. Không có field winner ("X"/"O") ở đây.
export type CaroMatchResult = {
  id: string;
  xPlayerId: string;
  xPlayerRatingBefore: number;
  xPlayerRatingAfter: number;
  oPlayerId: string;
  oPlayerRatingBefore: number;
  oPlayerRatingAfter: number;
  winnerId: string | null;
  finalBoard: CaroBoard;
  moves: CaroMove[];
  createdAt: string | null;
  updatedAt: string | null;
};

export type MatchEndedReason = "NORMAL" | "OUT_OF_TIME";

export const stateStore = $state<{
  state: UserState;
  currentMatch: CaroMatch | null;
  endedMatch: CaroMatchResult | null;
  endedReason: MatchEndedReason | null;
  connected: boolean;
  // đã nhận ít nhất 1 message từ server -> state bên dưới mới đáng tin
  hydrated: boolean;
}>({
  state: "IDLE",
  currentMatch: null,
  endedMatch: null,
  endedReason: null,
  connected: false,
  hydrated: false,
});
