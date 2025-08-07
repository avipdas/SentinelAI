import React, { useEffect, useState } from "react";

interface LogEntry {
  id: number;
  message: string;
  timestamp: string;
  anomaly: boolean;
  explanation?: string; // new field
}

function App() {
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8001/ws/logs");

    socket.onopen = () => {
      console.log("✅ WebSocket connected");
      setError(null);
    };

    socket.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (Array.isArray(data)) {
          setLogs(data);
        } else {
          setError("Unexpected data format from server");
        }
      } catch (err) {
        console.error("Failed to parse WebSocket message:", err);
      }
    };

    socket.onerror = () => {
      console.error("WebSocket error");
      setError("WebSocket connection failed.");
    };

    socket.onclose = () => {
      console.warn("WebSocket closed");
      setError("WebSocket connection closed.");
    };

    return () => {
      socket.close();
    };
  }, []);

  return (
    <div style={{ padding: "20px" }}>
      <h1>SentinelAI Dashboard (Live)</h1>
      {error && <p style={{ color: "red" }}>{error}</p>}
      {logs.length === 0 ? (
        <p>No logs found.</p>
      ) : (
        <table border={1} cellPadding={10}>
          <thead>
            <tr>
              <th>ID</th>
              <th>Message</th>
              <th>Timestamp</th>
              <th>Status</th>
              <th>Explanation</th>
            </tr>
          </thead>
          <tbody>
            {logs.map((log) => (
              <tr
                key={log.id}
                style={{ backgroundColor: log.anomaly ? "#ffcccc" : "#ccffcc" }}
              >
                <td>{log.id}</td>
                <td>{log.message}</td>
                <td>{new Date(log.timestamp).toLocaleString()}</td>
                <td style={{ color: log.anomaly ? "red" : "green" }}>
                  {log.anomaly ? "Suspicious" : "Normal"}
                </td>
                <td>{log.explanation || "N/A"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}

export default App;
