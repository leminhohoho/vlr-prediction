import os
import sqlite3
import pandas as pd
from src.features import get_duel_stats

module_dir = os.path.dirname(os.path.abspath(__file__))
conn = sqlite3.connect("/database/vlr.db")


def test_duel_stats():
    def test(match_id, map_id, player_id, opp_id, stats):
        duel_stats = get_duel_stats(conn, match_id, map_id, player_id, opp_id).iloc[0].to_list()

        for i in range(len(stats)):
            if duel_stats[i] != stats[i]:
                print(f"Want {stats}, get {duel_stats}")
                assert False

    test(530364, 10, 24895, 25255, [5, 3, 1, 1, 1, 0])
    test(530364, 10, 24895, 3977, [6, 4, 0, 1, 1, 0])

    assert True
