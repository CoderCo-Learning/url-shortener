"""
SQS event publisher for click analytics.

Set SQS_QUEUE_URL to enable. Without it, events are logged but not published.
"""

import os
import json
import time
import logging

logger = logging.getLogger(__name__)

_sqs_client = None


def _get_sqs():
    global _sqs_client
    if _sqs_client is not None:
        return _sqs_client

    queue_url = os.environ.get("SQS_QUEUE_URL")
    if not queue_url:
        return None

    import boto3
    _sqs_client = boto3.client("sqs")
    return _sqs_client


def publish_click_event(short_code: str, ip: str, user_agent: str, referer: str):
    event = {
        "short_code": short_code,
        "ip": ip,
        "user_agent": user_agent,
        "referer": referer,
        "timestamp": time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()),
    }

    sqs = _get_sqs()
    queue_url = os.environ.get("SQS_QUEUE_URL")

    if sqs and queue_url:
        try:
            sqs.send_message(
                QueueUrl=queue_url,
                MessageBody=json.dumps(event),
            )
        except Exception as e:
            logger.warning(f"Failed to publish click event: {e}")
    else:
        logger.info(f"Click event (no SQS): {json.dumps(event)}")
