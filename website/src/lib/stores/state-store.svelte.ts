export type UserState = "IDLE" | "QUEUING_CARO" | "PLAYING_CARO";

export type CaroPiece = "" | "X" | "O";

export type CaroMove = {
  piece: CaroPiece;
  x: number;
  y: number;
  playedAt: string;
};

export type CaroMatch = {
  id: string;
  isRated: boolean;
  xPlayerId: string;
  xPlayerRating: number;
  oPlayerId: string;
  oPlayerRating: number;
  board: CaroPiece[][];
  moves: CaroMove[];
  turnOf: CaroPiece;
  winner: CaroPiece;
  isEnded: boolean;
  startedAt: string;
};

export const stateStore = $state<{
  state: UserState;
  currentMatch: CaroMatch | null;
}>({
  state: "IDLE",
  currentMatch: null,
});