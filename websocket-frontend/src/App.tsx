import React from "react";
import "./App.css";
import MainContent from "./components/mainContent";
import { WebSocketProvider } from "./contexts/websocket";

function App() {
  return (
    <WebSocketProvider>
      <MainContent />
    </WebSocketProvider>
  );
}

export default App;
