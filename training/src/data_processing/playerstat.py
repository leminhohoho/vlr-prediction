import pandas as pd
import os
import sqlite3
from data_loader import (
    load_map_players_overview_stats,
    load_map_rounds_stats,
    load_players_highlights_log,
)
from utils import divide


def get_highlight_per_round(
    highlights_log: pd.DataFrame,
    rounds_stats: pd.DataFrame,
    team_id: int,
    player_id: int,
    highlight_type: str,
):
    player_highlights = highlights_log[
        (highlights_log["player_id"] == player_id)
        & (highlights_log["highlight_type"] == highlight_type)
    ]

    if highlight_type == "1v1":
        divider = 1
    elif highlight_type in ["1v2", "2k"]:
        divider = 2
    elif highlight_type in ["1v3", "3k"]:
        divider = 3
    elif highlight_type in ["1v4", "4k"]:
        divider = 4
    elif highlight_type in ["1v5", "5k"]:
        divider = 5

    if len(player_highlights) == 0:
        return {"def": 0, "atk": 0, "both": 0}

    def is_side(side):
        def get_side(row):
            round_stats = rounds_stats[
                rounds_stats["round_no"] == row["round_no"]
            ].iloc[0]
            if round_stats["team_def"] == team_id:
                return "def" == side

            return "atk" == side

        return get_side

    def_higlights = player_highlights[player_highlights.apply(is_side("def"), axis=1)]
    atk_higlights = player_highlights[player_highlights.apply(is_side("atk"), axis=1)]

    return {
        "def": len(def_higlights) / divider,
        "atk": len(atk_higlights) / divider,
        "both": (len(def_higlights) + len(atk_higlights)) / divider,
    }


# NOTE: rounds != score, and each side also include the ot rounds
def get_players_stats(
    conn: sqlite3.Connection,
    match_id: int,
    map_id: int,
    team_id: int,
    def_rounds: int,
    atk_rounds: int,
):
    players_overview_stats = load_map_players_overview_stats(
        conn, match_id, map_id, team_id
    )
    highlights_log = load_players_highlights_log(conn, match_id, map_id, team_id)
    rounds_stats = load_map_rounds_stats(conn, match_id, map_id, team_id)

    if len(players_overview_stats) != 10:
        return None

    players_stats = []

    for _, stats in players_overview_stats.groupby("player_id"):
        def_stats = stats[stats["side"] == "def"].iloc[0]
        atk_stats = stats[stats["side"] == "atk"].iloc[0]

        p1v1s = get_highlight_per_round(
            highlights_log, rounds_stats, team_id, def_stats["player_id"], "1v1"
        )
        p1v2s = get_highlight_per_round(
            highlights_log, rounds_stats, team_id, def_stats["player_id"], "1v2"
        )
        p1v3s = get_highlight_per_round(
            highlights_log, rounds_stats, team_id, def_stats["player_id"], "1v3"
        )
        p1v4s = get_highlight_per_round(
            highlights_log, rounds_stats, team_id, def_stats["player_id"], "1v4"
        )
        p1v5s = get_highlight_per_round(
            highlights_log, rounds_stats, team_id, def_stats["player_id"], "1v5"
        )
        p2ks = get_highlight_per_round(
            highlights_log, rounds_stats, team_id, def_stats["player_id"], "2k"
        )
        p3ks = get_highlight_per_round(
            highlights_log, rounds_stats, team_id, def_stats["player_id"], "3k"
        )
        p4ks = get_highlight_per_round(
            highlights_log, rounds_stats, team_id, def_stats["player_id"], "4k"
        )
        p5ks = get_highlight_per_round(
            highlights_log, rounds_stats, team_id, def_stats["player_id"], "5k"
        )

        player_stats = {
            "agent_id": int(def_stats["agent_id"]),
            "role": def_stats["role"],
            "def_rating": float(def_stats["rating"]),
            "atk_rating": float(atk_stats["rating"]),
            "rating": float(
                divide(
                    (
                        def_stats["rating"] * def_rounds
                        + atk_stats["rating"] * atk_rounds
                    ),
                    (def_rounds + atk_rounds),
                )
            ),
            "def_acs": float(def_stats["acs"]),
            "atk_acs": float(atk_stats["acs"]),
            "acs": float(
                divide(
                    (def_stats["acs"] * def_rounds + atk_stats["acs"] * atk_rounds),
                    (def_rounds + atk_rounds),
                )
            ),
            "def_kills_per_round": float(def_stats["kills"] / def_rounds),
            "atk_kills_per_round": float(atk_stats["kills"] / atk_rounds),
            "kills_per_round": float(
                divide(def_stats["kills"] + atk_stats["kills"], def_rounds + atk_rounds)
            ),
            "def_deaths_per_round": float(def_stats["deaths"] / def_rounds),
            "atk_deaths_per_round": float(atk_stats["deaths"] / atk_rounds),
            "deaths_per_round": float(
                divide(
                    def_stats["deaths"] + atk_stats["deaths"], def_rounds + atk_rounds
                )
            ),
            "def_assists_per_round": float(def_stats["assists"] / def_rounds),
            "atk_assists_per_round": float(atk_stats["assists"] / atk_rounds),
            "assists_per_round": float(
                divide(
                    def_stats["assists"] + atk_stats["assists"], def_rounds + atk_rounds
                )
            ),
            "def_kast": float(def_stats["kast"]),
            "atk_kast": float(atk_stats["kast"]),
            "kast": float(
                divide(
                    (def_stats["kast"] * def_rounds + atk_stats["kast"] * atk_rounds),
                    (def_rounds + atk_rounds),
                )
            ),
            "def_adr": float(def_stats["adr"]),
            "atk_adr": float(atk_stats["adr"]),
            "adr": float(
                divide(
                    (def_stats["adr"] * def_rounds + atk_stats["adr"] * atk_rounds),
                    (def_rounds + atk_rounds),
                )
            ),
            "def_hs": float(def_stats["hs"]),
            "atk_hs": float(atk_stats["hs"]),
            "hs": float(
                divide(
                    def_stats["hs"] * def_stats["kills"]
                    + atk_stats["hs"] * atk_stats["kills"],
                    def_stats["kills"] + atk_stats["kills"],
                )
            ),
            "def_fks_per_round": float(def_stats["first_kills"] / def_rounds),
            "atk_fks_per_round": float(atk_stats["first_kills"] / atk_rounds),
            "fks_per_round": float(
                divide(
                    def_stats["first_kills"] + atk_stats["first_kills"],
                    def_rounds + atk_rounds,
                )
            ),
            "def_fds_per_round": float(def_stats["first_deaths"] / def_rounds),
            "atk_fds_per_round": float(atk_stats["first_deaths"] / atk_rounds),
            "fds_per_round": float(
                divide(
                    def_stats["first_deaths"] + atk_stats["first_deaths"],
                    def_rounds + atk_rounds,
                )
            ),
            "def_1v1s_per_round": float(divide(p1v1s["def"], def_rounds)),
            "atk_1v1s_per_round": float(divide(p1v1s["atk"], atk_rounds)),
            "1v1s_per_round": float(divide(p1v1s["both"], def_rounds + atk_rounds)),
            "def_1v2s_per_round": float(divide(p1v2s["def"], def_rounds)),
            "atk_1v2s_per_round": float(divide(p1v2s["atk"], atk_rounds)),
            "1v2s_per_round": float(divide(p1v2s["both"], def_rounds + atk_rounds)),
            "def_1v3s_per_round": float(divide(p1v3s["def"], def_rounds)),
            "atk_1v3s_per_round": float(divide(p1v3s["atk"], atk_rounds)),
            "1v3s_per_round": float(divide(p1v3s["both"], def_rounds + atk_rounds)),
            "def_1v4s_per_round": float(divide(p1v4s["def"], def_rounds)),
            "atk_1v4s_per_round": float(divide(p1v4s["atk"], atk_rounds)),
            "1v4s_per_round": float(divide(p1v4s["both"], def_rounds + atk_rounds)),
            "def_1v5s_per_round": float(divide(p1v5s["def"], def_rounds)),
            "atk_1v5s_per_round": float(divide(p1v5s["atk"], atk_rounds)),
            "1v5s_per_round": float(divide(p1v5s["both"], def_rounds + atk_rounds)),
            "def_2ks_per_round": float(divide(p2ks["def"], def_rounds)),
            "atk_2ks_per_round": float(divide(p2ks["atk"], atk_rounds)),
            "2ks_per_round": float(divide(p2ks["both"], def_rounds + atk_rounds)),
            "def_3ks_per_round": float(divide(p3ks["def"], def_rounds)),
            "atk_3ks_per_round": float(divide(p3ks["atk"], atk_rounds)),
            "3ks_per_round": float(divide(p3ks["both"], def_rounds + atk_rounds)),
            "def_4ks_per_round": float(divide(p4ks["def"], def_rounds)),
            "atk_4ks_per_round": float(divide(p4ks["atk"], atk_rounds)),
            "4ks_per_round": float(divide(p4ks["both"], def_rounds + atk_rounds)),
            "def_5ks_per_round": float(divide(p5ks["def"], def_rounds)),
            "atk_5ks_per_round": float(divide(p5ks["atk"], atk_rounds)),
            "5ks_per_round": float(divide(p5ks["both"], def_rounds + atk_rounds)),
        }

        players_stats.append(player_stats)

    return players_stats


if __name__ == "__main__":
    import pprint

    pp = pprint.PrettyPrinter(indent=4)

    db_path = os.path.join(
        os.path.dirname(os.path.abspath(__file__)), "../../../database/vlr.db"
    )
    conn = sqlite3.connect(db_path)

    ps_stats = get_players_stats(conn, 530924, 4, 2593, 8, 12)
    pp.pprint(ps_stats)
