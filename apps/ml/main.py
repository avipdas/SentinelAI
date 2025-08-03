from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
import psycopg2
import pandas as pd
from sklearn.ensemble import IsolationForest
import uvicorn
import os

app = FastAPI()

# ✅ Enable CORS for React frontend
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000"],  # React dev server
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

DB_CONFIG = {
    "dbname": os.getenv("DB_NAME", "sentinelai"),
    "user": os.getenv("DB_USER", "sentinel"),
    "password": os.getenv("DB_PASS", "sentinel123"),
    "host": os.getenv("DB_HOST", "db"),   # ✅ use docker service name
    "port": os.getenv("DB_PORT", "5432")  # ✅ match docker-compose
}

@app.get("/analyze")
def analyze():
    try:
        conn = psycopg2.connect(**DB_CONFIG)
        df = pd.read_sql(
            "SELECT id, message, timestamp FROM logs ORDER BY timestamp DESC LIMIT 100;", 
            conn
        )
    except Exception as e:
        # Always return an array so React doesn't break
        return []

    if df.empty:
        conn.close()
        return []

    df['msg_len'] = df['message'].apply(len)

    model = IsolationForest(contamination=0.1, random_state=42)
    df['anomaly'] = model.fit_predict(df[['msg_len']])

    cur = conn.cursor()
    for _, row in df.iterrows():
        cur.execute(
            """
            INSERT INTO anomalies (message, timestamp, anomaly)
            VALUES (%s, %s, %s)
            ON CONFLICT DO NOTHING
            """,
            (row['message'], row['timestamp'], row['anomaly'] == -1)
        )

    conn.commit()
    cur.close()
    conn.close()

    return df[['id', 'message', 'timestamp', 'anomaly']].to_dict(orient="records")

@app.get("/api/anomalies")
def get_anomalies():
    try:
        conn = psycopg2.connect(**DB_CONFIG)
        cur = conn.cursor()
        cur.execute("SELECT id, message, timestamp, anomaly FROM anomalies ORDER BY timestamp DESC LIMIT 50;")
        rows = cur.fetchall()
        cur.close()
        conn.close()
    except Exception as e:
        return []

    return [
        {"id": r[0], "message": r[1], "timestamp": r[2], "anomaly": r[3]}
        for r in rows
    ]

if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=8001, reload=True)
