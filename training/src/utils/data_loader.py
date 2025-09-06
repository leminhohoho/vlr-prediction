import pandas as pd
import sqlite3


def load_duel_stats(conn: sqlite3.Connection, match_id: int, map_id: int, player_id: int, opp_id: int):
    return pd.read_sql(
        """
SELECT
    match_id, map_id,
    (CASE WHEN team_1_player_id = :player_id THEN team_1_player_id ELSE team_2_player_id END) AS player_id, 
    (CASE WHEN team_1_player_id = :opp_id THEN team_1_player_id ELSE team_2_player_id END) opp_id, 
    (CASE WHEN team_1_player_id = :player_id THEN  team_1_player_kills_vs_team_2_player ELSE  team_2_player_kills_vs_team_1_player END) AS player_kills_vs_opp, 
    (CASE WHEN team_1_player_id = :opp_id THEN  team_1_player_kills_vs_team_2_player ELSE  team_2_player_kills_vs_team_1_player END) AS opp_kills_vs_players, 
    (CASE WHEN team_1_player_id = :player_id THEN  team_1_player_first_kills_vs_team_2_player ELSE  team_2_player_first_kills_vs_team_1_player END) AS player_first_kills_vs_opp, 
    (CASE WHEN team_1_player_id = :opp_id THEN  team_1_player_first_kills_vs_team_2_player ELSE  team_2_player_first_kills_vs_team_1_player END) AS opp_first_kills_vs_players,
    (CASE WHEN team_1_player_id = :player_id THEN  team_1_player_op_kills_vs_team_2_player ELSE  team_2_player_op_kills_vs_team_1_player END) AS player_op_kills_vs_opp, 
    (CASE WHEN team_1_player_id = :opp_id THEN  team_1_player_op_kills_vs_team_2_player ELSE  team_2_player_op_kills_vs_team_1_player END) AS opp_op_kills_vs_players
FROM players_duel_stats 
WHERE 
    match_id  = :match_id AND 
    map_id = :map_id AND 
    :player_id IN(team_1_player_id, team_2_player_id) AND
    :opp_id IN(team_1_player_id, team_2_player_id)

        """,
        conn,
        params={"match_id": match_id, "map_id": map_id, "player_id": player_id, "opp_id": opp_id},
    )
