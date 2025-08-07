from fastapi import FastAPI, WebSocket
from fastapi.middleware.cors import CORSMiddleware
import asyncio
import psycopg2
import pandas as pd
from sklearn.ensemble import IsolationForest
import uvicorn
import os
from openai import OpenAI

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

DB_CONFIG = {
    "dbname": os.getenv("DB_NAME", "sentinelai"),
    "user": os.getenv("DB_USER", "sentinel"),
    "password": os.getenv("DB_PASS", "sentinel123"),
    "host": os.getenv("DB_HOST", "db"),
    "port": os.getenv("DB_PORT", "5432"),
}


client = OpenAI(api_key=os.getenv("OPENAI_API_KEY"))

def generate_explanation(message: str, anomaly: bool) -> str:
    """Use LLM to explain anomalies."""
    if not anomaly:
        return "No anomaly detected."

    prompt = f"Explain why this log entry might indicate a security threat:\n\nLog: {message}\n"

    try:
        response = client.chat.completions.create(
            model="gpt-4o-mini",   # smaller, cheaper, faster; change if needed
            messages=[
                {"role": "system", "content": "You are a cybersecurity analyst."},
                {"role": "user", "content": prompt}
            ],
            max_tokens=80,
            temperature=0.5
        )
        return response.choices[0].message.content.strip()
    except Exception as e:
        return f"Explanation unavailable: {str(e)}"

@app.get("/analyze")
def analyze():
    try:
        conn = psycopg2.connect(**DB_CONFIG)
        df = pd.read_sql(
            "SELECT id, message, timestamp FROM logs ORDER BY timestamp DESC LIMIT 100;", 
            conn
        )
    except Exception:
        return []

    if df.empty:
        conn.close()
        return []

    df['msg_len'] = df['message'].apply(len)
    model = IsolationForest(contamination=0.1, random_state=42)
    df['anomaly'] = model.fit_predict(df[['msg_len']])

    cur = conn.cursor()
    for _, row in df.iterrows():
        is_suspicious = row['anomaly'] == -1
        try:
            explanation = generate_explanation(row['message'], is_suspicious)
        except Exception as e:
            explanation = f"Explanation unavailable: {str(e)}"

        cur.execute(
            """
            INSERT INTO anomalies (message, timestamp, anomaly, explanation)
            VALUES (%s, %s, %s, %s)
            ON CONFLICT DO NOTHING
            """,
            (row['message'], row['timestamp'], is_suspicious, explanation)
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
        cur.execute("SELECT id, message, timestamp, anomaly, explanation FROM anomalies ORDER BY timestamp DESC LIMIT 50;")
        rows = cur.fetchall()
        cur.close()
        conn.close()
    except Exception:
        return []

    return [
        {"id": r[0], "message": r[1], "timestamp": r[2], "anomaly": r[3], "explanation": r[4]}
        for r in rows
    ]

@app.websocket("/ws/logs")
async def websocket_logs(websocket: WebSocket):
    await websocket.accept()
    while True:
        try:
            conn = psycopg2.connect(**DB_CONFIG)
            cur = conn.cursor()
            cur.execute("SELECT id, message, timestamp, anomaly, explanation FROM anomalies ORDER BY timestamp DESC LIMIT 20;")
            rows = cur.fetchall()
            cur.close()
            conn.close()

            logs = [
                {"id": r[0], "message": r[1], "timestamp": str(r[2]), "anomaly": r[3], "explanation": r[4]}
                for r in rows
            ]
            await websocket.send_json(logs)
        except Exception as e:
            await websocket.send_json({"error": str(e)})
        await asyncio.sleep(3)  # check every 3 seconds

if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=8001, reload=True)
