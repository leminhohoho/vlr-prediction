import pandas as pd
import sqlite3

from src.utils import load_players_stats

_cache = {}


def cache_get(query: str):
    if query in _cache:
        print(f"get {query} from cache")
        return _cache[query]

    return None


def cache_save(dat, query):
    _cache[query] = dat

    return dat


def compute_map_stats(conn: sqlite3.Connection):
    def f(row: pd.Series):
        # NOTE: Since there is only one map, player stats can be filtered based on match_id
        print(f"Processing map {row['match_id']}-{row['map_id']}")
        team_1_map_players_stats_query = f"{row['match_id']}-{row['map_id']}-{row['team_1_id']}"
        team_1_map_players_stats_df = cache_get(team_1_map_players_stats_query)
        if team_1_map_players_stats_df is None:
            team_1_map_players_stats_df = cache_save(
                load_players_stats(conn, row["team_1_id"], match_id=row["match_id"]), team_1_map_players_stats_query
            )
        team_1_map_players_stats_df = team_1_map_players_stats_df[team_1_map_players_stats_df["side"] == "def"]
        team_2_map_players_stats_query = f"{row['match_id']}-{row['map_id']}-{row['team_2_id']}"
        team_2_map_players_stats_df = cache_get(team_2_map_players_stats_query)
        if team_2_map_players_stats_df is None:
            team_2_map_players_stats_df = cache_save(
                load_players_stats(conn, row["team_2_id"], match_id=row["match_id"]), team_2_map_players_stats_query
            )
        team_2_map_players_stats_df = team_2_map_players_stats_df[team_2_map_players_stats_df["side"] == "def"]

        team_1_comps = sorted(team_1_map_players_stats_df["role"].to_list())
        team_2_comps = sorted(team_2_map_players_stats_df["role"].to_list())

        return pd.Series(
            {
                "team_1_def_score": row["team_1_def_score"],
                "team_1_atk_score": row["team_1_atk_score"],
                "team_1_ot_score": row["team_1_ot_score"],
                "team_2_def_score": row["team_2_def_score"],
                "team_2_atk_score": row["team_2_atk_score"],
                "team_2_ot_score": row["team_2_ot_score"],
                "duration": row["duration"],
                "team_1_agents": sorted(team_1_map_players_stats_df["agent_id"].unique().tolist()),
                "team_2_agents": sorted(team_2_map_players_stats_df["agent_id"].unique().tolist()),
                "team_1_duelists": team_1_comps.count("duelist"),
                "team_1_controllers": team_1_comps.count("controller"),
                "team_1_sentinels": team_1_comps.count("sentinel"),
                "team_1_initiators": team_1_comps.count("initiator"),
                "team_2_duelists": team_2_comps.count("duelist"),
                "team_2_controllers": team_2_comps.count("controller"),
                "team_2_sentinels": team_2_comps.count("sentinel"),
                "team_2_initiators": team_2_comps.count("initiator"),
            }
        )

    return f


def get_maps_stats(conn: sqlite3.Connection, match_maps_df: pd.DataFrame):
    return match_maps_df.apply(compute_map_stats(conn), axis=1)
