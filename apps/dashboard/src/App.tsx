import React, { useEffect, useState } from "react";
import axios from "axios";

interface LogEntry {
  id: number;
  message: string;
  timestamp: string;
  anomaly: number; // -1 = Suspicious, 1 = Normal
}

function App() {
  const [logs, setLogs] = useState<LogEntry[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    axios
      .get<LogEntry[]>("http://localhost:8001/analyze")
      .then((res) => {
        console.log("Logs fetched:", res.data); // ✅ Debug line
        setLogs(res.data);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Error fetching logs:", err);
        setLoading(false);
      });
  }, []);

  return (
    <div style={{ padding: "20px" }}>
      <h1>SentinelAI Dashboard</h1>
      {loading ? (
        <p>Loading logs...</p>
      ) : logs.length === 0 ? (
        <p>No logs found.</p>
      ) : (
        <table border={1} cellPadding={10}>
          <thead>
            <tr>
              <th>ID</th>
              <th>Message</th>
              <th>Timestamp</th>
              <th>Anomaly</th>
            </tr>
          </thead>
          <tbody>
            {logs.map((log) => (
              <tr
                key={log.id}
                style={{
                  backgroundColor: log.anomaly === -1 ? "#ffcccc" : "#ccffcc",
                }}
              >
                <td>{log.id}</td>
                <td>{log.message}</td>
                <td>{new Date(log.timestamp).toLocaleString()}</td>
                <td style={{ color: log.anomaly === -1 ? "red" : "green" }}>
                  {log.anomaly === -1 ? "Suspicious" : "Normal"}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}

export default App;
