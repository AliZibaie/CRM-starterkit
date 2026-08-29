"use client";

import { useEffect, useRef } from "react";
import { useAuthStore } from "@/stores/auth-store";
import type { TicketMessage } from "@/lib/support-api";

const WS_BASE =
  process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080/api/v1";

export function useTicketWebSocket(
  ticketId: string | null,
  onMessage: (message: TicketMessage) => void
) {
  const accessToken = useAuthStore((s) => s.accessToken);
  const handlerRef = useRef(onMessage);
  handlerRef.current = onMessage;

  useEffect(() => {
    if (!ticketId || !accessToken) return;

    let ws: WebSocket | null = null;
    let retries = 0;
    let timer: ReturnType<typeof setTimeout> | null = null;
    let closed = false;

    const connect = () => {
      ws = new WebSocket(`${WS_BASE}/support/ws/${ticketId}?token=${accessToken}`);
      ws.onopen = () => {
        retries = 0;
      };
      ws.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data as string) as {
            type?: string;
            message?: TicketMessage;
          };
          if (data.type === "message" && data.message) {
            handlerRef.current(data.message);
          }
        } catch {
          /* ignore malformed payloads */
        }
      };
      ws.onclose = () => {
        if (closed || retries >= 5) return;
        retries += 1;
        timer = setTimeout(connect, Math.min(retries * 2000, 10000));
      };
    };

    connect();

    return () => {
      closed = true;
      if (timer) clearTimeout(timer);
      ws?.close();
    };
  }, [ticketId, accessToken]);
}
