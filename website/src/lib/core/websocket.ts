import { PUBLIC_WS_URL } from "$env/static/public";
import { stateStore } from "$lib/stores/state-store.svelte";
import { goto } from "$app/navigation";
import { toast } from "svelte-sonner";

let ws: WebSocket | null = null;

export function connectWs(token: string) {
  if(ws?.readyState === WebSocket.OPEN) {
    sendPingMessage();
    return;
  }

  ws = new WebSocket(PUBLIC_WS_URL + `?token=${token}`);

  ws.onopen = () => {
    sendPingMessage();
  };
  
  ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    stateStore.state = message.state;
    toast.info(stateStore.state);

    if(message.state === "PLAYING_CARO") {
      goto("/caro/match");
    }
  };
}

export function sendPingMessage() {
  if(!ws|| ws?.readyState !== WebSocket.OPEN) return;
  ws.send(JSON.stringify({
    code: "PING"
  }));
}

export function sendCaroJoinQueueMessage() {
  if(!ws|| ws?.readyState !== WebSocket.OPEN) return;
  ws.send(JSON.stringify({
    code: "CARO:JOIN_QUEUE"
  }));
}