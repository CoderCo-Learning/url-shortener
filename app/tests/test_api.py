import os
os.environ.setdefault("TABLE_NAME", "test-table")

from unittest.mock import patch, MagicMock
from fastapi.testclient import TestClient
from src.main import app

client = TestClient(app)


def test_health():
    r = client.get("/healthz")
    assert r.status_code == 200
    assert r.json()["status"] == "ok"
    assert "db" in r.json()


@patch("src.db.put_mapping")
def test_shorten(mock_put):
    r = client.post("/shorten", json={"url": "https://example.com"})
    assert r.status_code == 200
    data = r.json()
    assert "short" in data
    assert data["url"] == "https://example.com"
    mock_put.assert_called_once()


def test_shorten_missing_url():
    r = client.post("/shorten", json={})
    assert r.status_code == 400


@patch("src.db.get_mapping", return_value={"id": "abc", "url": "https://example.com"})
def test_resolve(mock_get):
    r = client.get("/abc", follow_redirects=False)
    assert r.status_code == 307 or r.status_code == 302


@patch("src.db.get_mapping", return_value=None)
def test_resolve_not_found(mock_get):
    r = client.get("/nonexistent")
    assert r.status_code == 404
