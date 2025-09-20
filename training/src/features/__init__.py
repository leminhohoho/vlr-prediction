from src.utils.data_loader import load_matches, load_players_stats
from .duel_stats import get_duel_stats
from .players_stats import get_players_stats
from .rounds_stats import get_rounds_stats
from .match_maps_stats import get_match_maps_stats
from .maps_stats import get_maps_stats
from src.utils import load_maps
import pandas as pd
import sqlite3


def append_maps(conn: sqlite3.Connection, match, maps: pd.DataFrame, players_stats_df: pd.DataFrame):
    team_match_players_stats = players_stats_df[players_stats_df["match_id"] == match["id"]]
    opp_team_match_players_stats = load_players_stats(conn, match["team_against_id"], match_id=match["id"])
    match_maps = maps[maps["match_id"] == match["id"]]

    match_maps_stats = get_match_maps_stats(match_maps, team_match_players_stats, opp_team_match_players_stats).to_dict(orient="records")

    for match_map in match_maps_stats:
        m = maps[maps["map_id"] == match_map["map_id"]]
        print(len(m))
        match_map["map_form"] = get_maps_stats(conn, m).to_dict(orient="records")

        match_map.pop("map_id")

    return match_maps_stats


def compute_stats(conn: sqlite3.Connection, team_1_id: int, team_2_id: int, stage: str, tier: int, date: str):
    team_1_matches = load_matches(conn, team_1_id, date=date).to_dict(orient="records")
    team_2_matches = load_matches(conn, team_2_id, date=date).to_dict(orient="records")
    team_1_players_stats_df = load_players_stats(conn, team_1_id, date=date)
    team_2_players_stats_df = load_players_stats(conn, team_2_id, date=date)
    maps = load_maps(conn, date)

    for match in team_1_matches:
        match["maps"] = append_maps(conn, match, maps, team_1_players_stats_df)
        match.pop("id")
        match.pop("url")
        match.pop("date")
        match.pop("tournament_id")
        match.pop("team_id")
        match.pop("team_against_id")

    for match in team_2_matches:
        match["maps"] = append_maps(conn, match, maps, team_2_players_stats_df)
        match.pop("id")
        match.pop("url")
        match.pop("date")
        match.pop("tournament_id")
        match.pop("team_id")
        match.pop("team_against_id")

    return {
        "team_1_matches": team_1_matches,
        "team_2_matches": team_2_matches,
        "stage": stage,
        "tier": tier,
    }
