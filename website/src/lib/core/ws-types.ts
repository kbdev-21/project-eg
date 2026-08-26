import type { CaroMatch, UserState } from "$lib/stores/state-store.svelte";

export const CARO_BOARD_SIZE = 15;

// domain.CARO_MAX_MOVE_TIME
export const CARO_MAX_MOVE_TIME_MS = 20_000;

export type ClientMessage =
  | { code: "PING" }
  | { code: "CARO:JOIN_QUEUE" }
  | { code: "CARO:LEAVE_QUEUE" }
  | { code: "CARO:PLAY_MOVE"; data: { x: number; y: number } };

export type ServerMessageCode =
  | "OK"
  | "ERROR"
  | "CARO:MATCH_FOUND"
  | "CARO:NEW_BOARD_STATE"
  | "CARO:MATCH_ENDED"
  | "CARO:MATCH_ENDED_OUT_OF_TIME";

// data.match kèm theo khi state === "PLAYING_CARO" (BuildServerMessage),
// và kèm match đã kết thúc ở 2 code MATCH_ENDED*.
export type ServerMessage = {
  code: ServerMessageCode;
  state: UserState;
  data: { match?: CaroMatch } | null;
};
