import sqlite3
import os
import pandas as pd


def get_rounds_stats(
    conn: sqlite3.Connection, match_id: int, map_id: int, team_id: int
):
    rounds_stats = load_map_rounds_stats(conn, match_id, map_id, team_id)
