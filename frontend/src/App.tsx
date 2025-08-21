import React, { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import LoginPage from "./pages/LoginPage";
import ChatPage from "./pages/ChatPage";

const App: React.FC = () => {
  const [token, setToken] = useState<string | null>(localStorage.getItem("token"));

  // Keep state in sync with localStorage changes
  useEffect(() => {
    if (token) {
      localStorage.setItem("token", token);
    } else {
      localStorage.removeItem("token");
    }
  }, [token]);

  return (
    <Router>
      <Routes>
        {/* if no token, go to login */}
        <Route
          path="/"
          element={token ? <Navigate to="/chat" replace /> : <Navigate to="/login" replace />}
        />
        <Route
          path="/login"
          element={!token ? <LoginPage setToken={setToken} /> : <Navigate to="/" replace />}
        />
        <Route
          path="/chat"
          element={token ? <ChatPage token={token} setToken={setToken} /> : <Navigate to="/login" replace />}
        />
      </Routes>
    </Router>
  );
};

export default App;