// src/contexts/WebSocketContext.tsx
import React, { createContext, useContext, useEffect, useState } from "react";

const API_URL = process.env.REACT_APP_API_URL || "";

interface WebSocketContextType {
  socket: WebSocket | null;
  sendMessage: (message: object) => void;
  isLoading: boolean;
  url: string | null;
}

interface WebSocketMessage {
  url: string;
}

const WebSocketContext = createContext<WebSocketContextType | undefined>(
  undefined
);

export const WebSocketProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [url, setUrl] = useState<string | null>(null);

  useEffect(() => {
    const ws = new WebSocket(API_URL);
    setSocket(ws);

    ws.onopen = () => {
      console.log("Connected to WebSocket");
    };

    ws.onmessage = (event) => {
      console.log("Message from WebSocket: ", event.data);
      try {
        const data: WebSocketMessage = JSON.parse(event.data);
        setUrl(data.url);
      } catch (error) {
        console.error("Error parsing WebSocket message: ", error);
      } finally {
        setIsLoading(false);
        ws.close();
      }
    };

    ws.onerror = (error) => {
      console.error("WebSocket Error: ", error);
    };

    ws.onclose = () => {
      console.log("WebSocket connection closed");
    };

    return () => {
      ws.close();
    };
  }, []);

  const sendMessage = (message: object) => {
    if (socket && socket.readyState === WebSocket.OPEN) {
      setIsLoading(true);
      socket.send(JSON.stringify(message));
    }
  };

  return (
    <WebSocketContext.Provider value={{ socket, sendMessage, isLoading, url }}>
      {children}
    </WebSocketContext.Provider>
  );
};

export const useWebSocket = (): WebSocketContextType => {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error("useWebSocket must be used within a WebSocketProvider");
  }
  return context;
};
