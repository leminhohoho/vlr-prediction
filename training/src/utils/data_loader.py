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


def load_players_stats(conn: sqlite3.Connection, team_id: int, date: str, duration=180):
    return pd.read_sql(
        """
SELECT 
    pos.*, 
      (
        SELECT agent_type FROM agents
        WHERE id = pos.agent_id
        LIMIT 1
    ) AS role
FROM player_overview_stats AS pos
JOIN matches AS m ON m.id = pos.match_id
WHERE 
    pos.team_id = :team_id AND
    m.date > date(:date, :duration) AND
    m.date < :date
        """,
        conn,
        params={"team_id": team_id, "date": date, "duration": f"-{duration} days"},
    )


def load_rounds_stats(conn: sqlite3.Connection, team_id: int, date: str, duration=180):
    return pd.read_sql(
        """
SELECT
    rs.match_id,
    rs.map_id,
    rs.round_no,
    (
        CASE
            WHEN rs.team_1_id = :team_id THEN rs.team_1_id
            WHEN rs.team_2_id = :team_id THEN rs.team_2_id
        END
    ) AS team_id,
    (
        CASE
            WHEN rs.team_1_id = :team_id THEN rs.team_2_id
            WHEN rs.team_2_id = :team_id THEN rs.team_1_id
        END
    ) AS team_against_id,
    (
        CASE
            WHEN rs.team_1_id = :team_id THEN rs.team_1_buy_type
            WHEN rs.team_2_id = :team_id THEN rs.team_2_buy_type
        END
    ) AS team_buy_type,
    (
        CASE
            WHEN rs.team_1_id = :team_id THEN rs.team_2_buy_type
            WHEN rs.team_2_id = :team_id THEN rs.team_1_buy_type
        END
    ) AS team_against_buy_type,
    (
        CASE
            WHEN rs.team_1_id = :team_id THEN rs.team_1_bank
            WHEN rs.team_2_id = :team_id THEN rs.team_2_bank
        END
    ) AS team_bank,
    (
        CASE
            WHEN rs.team_1_id = :team_id THEN rs.team_2_bank
            WHEN rs.team_2_id = :team_id THEN rs.team_1_bank
        END
    ) AS team_against_bank,
    (rs.team_def = :team_id) AS def,
    (rs.team_won = :team_id) AS won,
    rs.won_method,
    m."date"
FROM
    round_stats AS rs
    JOIN matches AS m ON m.id = rs.match_id 
WHERE
    (rs.team_1_id = :team_id OR rs.team_2_id = :team_id) AND
    m.date > date(:date, :duration) AND
    m.date < :date
ORDER BY m."date" DESC, rs.match_id, rs.map_id, round_no ASC
        """,
        conn,
        params={"team_id": team_id, "date": date, "duration": f"-{duration} days"},
    )


def load_highlights(conn: sqlite3.Connection, team_id: int, date: str, duration=180):
    return pd.read_sql(
        """
SELECT ph.*, m."date" FROM 
player_highlights AS ph
JOIN matches AS m ON m.id = ph.match_id
WHERE 
    ph.team_id = :team_id AND
    m.date > date(:date, :duration) AND
    m.date < :date
        """,
        conn,
        params={"team_id": team_id, "date": date, "duration": f"-{duration} days"},
    )
