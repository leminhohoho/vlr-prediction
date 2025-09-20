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


def load_players_stats(conn: sqlite3.Connection, team_id: int, date=None, duration=180, match_id=None, map_id=None):
    query = """
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
    pos.team_id = :team_id
        """

    if date is not None:
        query += f" AND m.date > date('{date}', '-{duration} days') AND m.date < '{date}'"

    if match_id is not None:
        query += f" AND pos.match_id = {match_id}"

    if map_id is not None:
        query += f" AND pos.map_id = {map_id}"

    return pd.read_sql(query, conn, params={"team_id": team_id})


def load_rounds_stats(conn: sqlite3.Connection, team_id: int, date=None, duration=180):
    query = """
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
    (rs.team_1_id = :team_id OR rs.team_2_id = :team_id)
        """

    if date is not None:
        query += f" AND m.date > date('{date}', '-{duration} days') AND m.date < '{date}'"

    query += ' ORDER BY m."date" DESC, rs.match_id, rs.map_id, round_no ASC'

    return pd.read_sql(query, conn, params={"team_id": team_id})


def load_highlights(conn: sqlite3.Connection, team_id: int, date=None, duration=180):
    query = """
SELECT ph.*, m."date" FROM 
player_highlights AS ph
JOIN matches AS m ON m.id = ph.match_id
WHERE 
    ph.team_id = :team_id
        """

    if date is not None:
        query += f" AND m.date > date('{date}', '-{duration} days') AND m.date < '{date}'"

    return pd.read_sql(query, conn, params={"team_id": team_id})


def load_maps(conn: sqlite3.Connection, date=None, duration=180, match_id=-1):
    """
    match_id is optional, for distinguish between the curent match that the past matches is fetched for and the other matches
    """
    query = """
SELECT 
    mm.match_id,
    mm.map_id,
    mm.team_1_id,
    mm.team_2_id,
    mm.duration,
    mm.team_1_def_score,
    mm.team_1_atk_score,
    mm.team_1_ot_score,
    mm.team_2_def_score,
    mm.team_2_atk_score,
    mm.team_2_ot_score,
    (CASE WHEN mm.team_def_first = mm.team_1_id THEN 1 ELSE 2 END) AS team_def_first,
    (CASE
        WHEN mm.team_pick = mm.team_1_id THEN 1
        WHEN mm.team_pick = mm.team_2_id THEN 2
        ELSE 3
    END) AS team_pick,
    m."date"
FROM 
    match_maps AS mm
    JOIN matches AS m ON m.id = mm.match_id 
WHERE 
    mm.match_id != :match_id
    """

    if date is not None:
        query += f" AND m.date > date('{date}', '-{duration} days') AND m.date < '{date}'"

    return pd.read_sql(query, conn, params={"match_id": match_id})


def load_matches(conn: sqlite3.Connection, team_id: int, date=None, duration=180):
    query = """
SELECT 
    id, url, date, tournament_id, stage,
    (CASE WHEN :team_id = team_1_id THEN team_1_id ELSE team_2_id END) AS team_id,
    (CASE WHEN :team_id = team_1_id THEN team_2_id ELSE team_1_id END) AS team_against_id,
    (CASE WHEN :team_id = team_1_id THEN team_1_score ELSE team_2_score END) AS team_score,
    (CASE WHEN :team_id = team_1_id THEN team_2_score ELSE team_1_score END) AS team_against_score,
    (CASE WHEN :team_id = team_1_id THEN team_1_rating ELSE team_2_rating END) AS team_rating,
    (CASE WHEN :team_id = team_1_id THEN team_2_rating ELSE team_1_rating END) AS team_against_rating
FROM matches 
WHERE (team_1_id = :team_id OR team_2_id = :team_id)
    """

    if date is not None:
        query += f""" AND "date" > date('{date}', '-{duration} days') AND "date" < '{date}'"""

    return pd.read_sql(query, conn, params={"team_id": team_id})
