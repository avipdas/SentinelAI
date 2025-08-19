import typer
import requests
import json

app = typer.Typer()
API_URL = "http://localhost:8001"

@app.command()
def analyze():
    """Run anomaly detection on recent logs."""
    resp = requests.get(f"{API_URL}/analyze")
    if resp.status_code != 200:
        typer.echo("❌ Error running analysis.")
        raise typer.Exit()

    data = resp.json()
    typer.echo(f"✅ Analysis complete. {len(data)} entries analyzed.\n")
    for entry in data:
        status = "🚨 ANOMALY" if entry["anomaly"] == -1 else "✅ Normal"
        typer.echo(f"[{entry['timestamp']}] {status}: {entry['message'][:60]}")

@app.command()
def anomalies(show_explanations: bool = typer.Option(False, "--explain", help="Show LLM explanations")):
    """Fetch and display past detected anomalies."""
    resp = requests.get(f"{API_URL}/api/anomalies")
    if resp.status_code != 200:
        typer.echo("❌ Error fetching anomalies.")
        raise typer.Exit()

    data = resp.json()
    for entry in data:
        if entry["anomaly"]:
            typer.echo(f"\n🧠 [{entry['timestamp']}] {entry['message']}")
            if show_explanations:
                typer.echo(f"   ➤ {entry['explanation']}")

@app.command()
def export(output: str = typer.Argument("anomalies.json")):
    """Export anomalies to a JSON file."""
    resp = requests.get(f"{API_URL}/api/anomalies")
    if resp.status_code != 200:
        typer.echo("❌ Error exporting data.")
        raise typer.Exit()

    with open(output, "w") as f:
        json.dump(resp.json(), f, indent=2)
    typer.echo(f"✅ Exported to {output}")

if __name__ == "__main__":
    app()
