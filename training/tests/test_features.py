import os
import sqlite3
from src.features import get_duel_stats, get_players_stats, get_rounds_stats, get_match_maps_stats
from src.utils import load_players_stats, load_rounds_stats, load_highlights, load_maps

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


def test_players_stats():
    players_stats = load_players_stats(conn, 878, "2025-09-07")
    rounds = load_rounds_stats(conn, 878, "2025-09-07")
    highlights = load_highlights(conn, 878, "2025-09-07")

    players_stats = players_stats[(players_stats["match_id"] == 530364) & (players_stats["map_id"] == 10)]
    rounds = rounds[(rounds["match_id"] == 530364) & (rounds["map_id"] == 10)]
    highlights = highlights[(highlights["match_id"] == 530364) & (highlights["map_id"] == 10)]

    players = get_players_stats(players_stats, rounds, highlights)
    players.to_csv("/tmp/players_test.csv")
    print()
    print(players)


def test_round_stats():
    rounds = load_rounds_stats(conn, 878, "2025-09-07")
    highlights = load_highlights(conn, 878, "2025-09-07")
    opp_highlights = load_highlights(conn, 8304, "2025-09-07")

    rounds = rounds[(rounds["match_id"] == 530364) & (rounds["map_id"] == 10)]
    highlights = highlights[(highlights["match_id"] == 530364) & (highlights["map_id"] == 10)]
    opp_highlights = opp_highlights[(opp_highlights["match_id"] == 530364) & (opp_highlights["map_id"] == 10)]

    rounds_stats = get_rounds_stats(rounds, highlights, opp_highlights)
    rounds_stats.to_csv("/tmp/rounds_test.csv")

    print()
    print(rounds_stats)


def test_match_maps_stats():
    players_stats = load_players_stats(conn, 878, "2025-09-07", match_id=530364)
    opp_players_stats = load_players_stats(conn, 8304, "2025-09-07", match_id=530364)
    maps = load_maps(conn, "2025-09-07")

    maps = maps[maps["match_id"] == 530364]

    match_maps_stats = get_match_maps_stats(maps, players_stats, opp_players_stats)
    match_maps_stats.to_csv("/tmp/match_maps_test.csv")

    print()
    print(match_maps_stats)
