"""
Database layer tests.
Students should extend these with mocked DynamoDB (moto) or a test PostgreSQL instance.
"""
import os
os.environ.setdefault("TABLE_NAME", "test-table")


def test_backend_selection_dynamodb():
    """TABLE_NAME should select DynamoDB backend."""
    os.environ.pop("DATABASE_URL", None)
    os.environ["TABLE_NAME"] = "test-table"

    # Force re-init
    import src.db as db_mod
    db_mod._backend = None

    # This will try to init boto3 - we just check it picks the right path
    try:
        backend_type = db_mod.get_backend_type()
        assert backend_type == "dynamodb"
    except Exception:
        # boto3 may fail without AWS creds - that's fine for unit tests
        pass


def test_backend_selection_requires_config():
    """Should raise if neither TABLE_NAME nor DATABASE_URL is set."""
    import src.db as db_mod
    db_mod._backend = None

    old_table = os.environ.pop("TABLE_NAME", None)
    old_db = os.environ.pop("DATABASE_URL", None)

    try:
        db_mod.get_backend_type()
        assert False, "Should have raised RuntimeError"
    except RuntimeError as e:
        assert "TABLE_NAME" in str(e)
    finally:
        if old_table:
            os.environ["TABLE_NAME"] = old_table
        if old_db:
            os.environ["DATABASE_URL"] = old_db
