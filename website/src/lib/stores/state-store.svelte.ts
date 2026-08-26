export type UserState = "IDLE" | "QUEUING_CARO" | "PLAYING_CARO";

export type CaroPiece = "" | "X" | "O";

export type CaroStatus = "PLAYING" | "X_WON" | "O_WON" | "DRAW";

export type CaroEndReason = "" | "FIVE_IN_ROW" | "FULL_BOARD" | "OUT_OF_TIME";

export type CaroMove = {
  piece: CaroPiece;
  x: number;
  y: number;
  playedAt: string;
};

// board[y][x] - backend index y trước (domain/caro.go)
export type CaroBoard = CaroPiece[][];

// Dùng chung cho match đang chơi lẫn match đã kết thúc.
// Field null = chỉ biết sau khi trận kết thúc.
export type CaroMatch = {
  id: string;
  isRated: boolean;
  xPlayerId: string;
  xPlayerRatingBefore: number;
  xPlayerRatingAfter: number | null;
  oPlayerId: string;
  oPlayerRatingBefore: number;
  oPlayerRatingAfter: number | null;
  board: CaroBoard;
  moves: CaroMove[];
  turnOf: CaroPiece;
  status: CaroStatus;
  endReason: CaroEndReason;
  startedAt: string;
  endedAt: string | null;
};

export const stateStore = $state<{
  state: UserState;
  match: CaroMatch | null;
  connected: boolean;
  // đã nhận ít nhất 1 message từ server -> state bên dưới mới đáng tin
  hydrated: boolean;
}>({
  state: "IDLE",
  match: null,
  connected: false,
  hydrated: false,
});
