import { PUBLIC_WS_URL } from "$env/static/public";
import { stateStore } from "$lib/stores/state-store.svelte";
import { goto } from "$app/navigation";
import { page } from "$app/state";
import { toast } from "svelte-sonner";
import type { ClientMessage, ServerMessage } from "./ws-types";

const MATCH_ROUTE = "/caro/match";
const RECONNECT_BASE_DELAY = 1000;
const RECONNECT_MAX_DELAY = 10_000;

let ws: WebSocket | null = null;
let currentToken: string | null = null;
let closedByUs = false;
let retryCount = 0;

export function connectWs(token: string) {
  currentToken = token;
  closedByUs = false;
  // open() tự guard ws != null (gồm cả CONNECTING) nên không tạo socket trùng
  open();
}

export function disconnectWs() {
  closedByUs = true;
  ws?.close();
  ws = null;
  stateStore.connected = false;
}

function open() {
  if (ws || closedByUs || !currentToken) return;

  ws = new WebSocket(PUBLIC_WS_URL + `?token=${currentToken}`);

  ws.onopen = () => {
    retryCount = 0;
    stateStore.connected = true;
    // server không gửi gì lúc connect -> tự PING để lấy state
    sendPingMessage();
  };

  ws.onmessage = (event) => {
    let message: ServerMessage;
    try {
      message = JSON.parse(event.data);
    } catch {
      return;
    }
    handleServerMessage(message);
  };

  // không cần onerror: onclose luôn chạy sau nó và lo phần reconnect
  ws.onclose = () => {
    ws = null;
    stateStore.connected = false;
    if (closedByUs) return;
    scheduleReconnect();
  };
}

function scheduleReconnect() {
  // không cần giữ/cancel timer: open() đã guard closedByUs và ws != null
  const delay = Math.min(RECONNECT_BASE_DELAY * 2 ** retryCount, RECONNECT_MAX_DELAY);
  retryCount += 1;
  setTimeout(open, delay);
}

function handleServerMessage(message: ServerMessage) {
  stateStore.state = message.state;
  stateStore.hydrated = true;

  const match = message.data?.match ?? null;

  if (match) {
    stateStore.match = match;
  } else if (
    message.state !== "PLAYING_CARO" &&
    stateStore.match?.status === "PLAYING"
  ) {
    // chỉ dọn match mình tưởng còn sống; match đã kết thúc giữ lại để hiện kết quả
    stateStore.match = null;
  }

  if (message.code === "ERROR") {
    // server không kèm lý do -> coi như state client sai, resync lại
    toast.error("Thao tác không hợp lệ");
    sendPingMessage();
  }

  // Điều hướng dựa trên state, KHÔNG dựa trên code:
  // CARO:MATCH_FOUND còn bị tái dụng làm ack cho CARO:LEAVE_QUEUE.
  if (message.state === "PLAYING_CARO" && page.url.pathname !== MATCH_ROUTE) {
    goto(MATCH_ROUTE);
  }
}

function send(message: ClientMessage): boolean {
  if (!ws || ws.readyState !== WebSocket.OPEN) return false;
  ws.send(JSON.stringify(message));
  return true;
}

export function sendPingMessage() {
  return send({ code: "PING" });
}

export function sendCaroJoinQueueMessage() {
  return send({ code: "CARO:JOIN_QUEUE" });
}

export function sendCaroLeaveQueueMessage() {
  return send({ code: "CARO:LEAVE_QUEUE" });
}

export function sendCaroPlayMoveMessage(x: number, y: number) {
  return send({ code: "CARO:PLAY_MOVE", data: { x, y } });
}
