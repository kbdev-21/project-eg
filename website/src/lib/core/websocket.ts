import { PUBLIC_WS_URL } from "$env/static/public";
import { stateStore } from "$lib/stores/state-store.svelte";
import { toast } from "svelte-sonner";

let ws: WebSocket | null = null;

export function connectWs(token: string) {
  if(ws?.readyState === WebSocket.OPEN) return;
  ws = new WebSocket(PUBLIC_WS_URL + `?token=${token}`);

  ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    stateStore.state = message.state;
    toast.info(stateStore.state);
  } 
}

export function sendPingMessage() {
  if(!ws|| ws?.readyState !== WebSocket.OPEN) return;
  ws.send(JSON.stringify({
    code: "PING"
  }));
}