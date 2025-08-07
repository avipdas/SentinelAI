import React, { useEffect, useState } from "react";
import axios from "axios";

interface Anomaly {
  id: number;
  message: string;
  timestamp: string;
  anomaly: boolean;
}

function Anomalies() {
  const [anomalies, setAnomalies] = useState<Anomaly[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    axios
      .get<Anomaly[]>("http://localhost:8001/api/anomalies")
      .then((res) => {
        setAnomalies(res.data);
        setLoading(false);
      })
      .catch((err) => {
        console.error("Error fetching anomalies:", err);
        setLoading(false);
      });
  }, []);

  return (
    <div style={{ padding: "20px" }}>
      <h1>Anomalies History</h1>
      {loading ? (
        <p>Loading anomalies...</p>
      ) : anomalies.length === 0 ? (
        <p>No anomalies found.</p>
      ) : (
        <table border={1} cellPadding={10}>
          <thead>
            <tr>
              <th>ID</th>
              <th>Message</th>
              <th>Timestamp</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            {anomalies.map((a) => (
              <tr key={a.id} style={{ backgroundColor: a.anomaly ? "#ffcccc" : "#ccffcc" }}>
                <td>{a.id}</td>
                <td>{a.message}</td>
                <td>{new Date(a.timestamp).toLocaleString()}</td>
                <td style={{ color: a.anomaly ? "red" : "green" }}>
                  {a.anomaly ? "Suspicious" : "Normal"}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}

export default Anomalies;
