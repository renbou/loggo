import logging
import random
import signal
import sys
import time

import structlog


def main():
    # structlog configuration for structured logging,
    # mostly taken from the official docs: https://www.structlog.org/en/stable/standard-library.html
    logging.basicConfig(
        format="%(message)s",
        stream=sys.stdout,
        level=logging.INFO,
    )

    structlog.configure(
        processors=[
            structlog.stdlib.add_logger_name,
            structlog.stdlib.add_log_level,
            structlog.processors.TimeStamper(fmt="iso"),
            structlog.processors.StackInfoRenderer(),
            structlog.processors.EventRenamer("message"),
            structlog.processors.JSONRenderer(),
        ],
        wrapper_class=structlog.stdlib.BoundLogger,
        logger_factory=structlog.stdlib.LoggerFactory(),
    )
    logger = structlog.stdlib.get_logger("main")

    # watch for signals to support graceful shutdown
    alive = True

    def shutdown(*_):
        nonlocal alive
        alive = False

    for sig in (signal.SIGTERM, signal.SIGINT):
        signal.signal(sig, shutdown)

    logger.warn("started python example app", awesome=True)

    # log random messages, both structured and unstructured,
    # to showcase loggo's filtration capabilities
    while alive:
        time.sleep(random.random())
        if random.randint(0, 1) == 0:
            logger.info(
                "here's a structured log message, filterable by field",
                dice_roll=random.randint(1, 6),
            )
        else:
            logging.info(
                "and this is a non-structured log message, filterable by substring"
            )

    logger.warn("shutting down python example app", farewell="See you!")
