import sqlite3
from ..utils import load_duel_stats


def get_duel_stats(conn: sqlite3.Connection, match_id: int, map_id: int, player_id: int, opp_id: int):
    duel_stats_df = load_duel_stats(conn, match_id, map_id, player_id, opp_id)

    return duel_stats_df.drop(["match_id", "map_id", "player_id", "opp_id"], axis=1)
