import pandas as pd
import sqlite3


def load_matches(conn: sqlite3.Connection):
    return pd.read_sql(
        """
SELECT 
    matches.*, tournaments.tier_1 
FROM matches 
JOIN tournaments ON tournaments.id = matches.tournament_id 
ORDER BY date DESC
""",
        conn,
    )


def load_team_played_maps_recently(conn, team_id, date):
    return pd.read_sql(
        """
SELECT
    mms.*, m.date, m.stage, m.team_1_rating, m.team_2_rating,
    (mms.team_1_def_score + mms.team_1_atk_score + mms.team_1_ot_score) AS team_1_score,
    (mms.team_2_def_score + mms.team_2_atk_score + mms.team_2_ot_score) AS team_2_score
FROM match_maps AS mms JOIN matches AS m ON mms.match_id = m.id
WHERE
    (m.team_1_id = :team_id OR m.team_2_id = :team_id)
    AND "date" > date(:date, '-180 days')
    AND "date" < :date
""",
        conn,
        params={"team_id": team_id, "date": date},
    )


def load_team_maps_stats_recently(conn: sqlite3.Connection, team_id, date):
    return pd.read_sql(
        """
WITH mms AS (
  SELECT
    (team_1_def_score + team_1_atk_score + team_1_ot_score) AS team_1_score,
    (team_2_def_score + team_2_atk_score + team_2_ot_score) AS team_2_score,
    *
  FROM match_maps
)
SELECT
  mm.map_id,
  SUM(CASE
    WHEN mm.team_1_id = :team_id AND mm.team_1_score > mm.team_2_score THEN 1
    WHEN mm.team_2_id = :team_id AND mm.team_1_score < mm.team_2_score THEN 1
    ELSE 0
  END) AS wins,
  SUM(CASE
    WHEN mm.team_1_id = :team_id AND mm.team_1_score < mm.team_2_score THEN 1
    WHEN mm.team_2_id = :team_id AND mm.team_1_score > mm.team_2_score THEN 1
    ELSE 0
  END) AS losses
FROM mms AS mm
  JOIN matches AS m ON m.id = mm.match_id
WHERE (
    mm.team_1_id = :team_id OR mm.team_2_id = :team_id
  )
  AND m."date" > date(:date, '-180 days')
  AND m."date" < :date
GROUP BY mm.map_id
        """,
        conn,
        params={"team_id": team_id, "date": date},
    )


def load_current_map_pool(conn: sqlite3.Connection, date):
    return pd.read_sql(
        """
SELECT map_id
FROM maps_pool
WHERE date = (
    SELECT MAX(date)
    FROM maps_pool
    WHERE date <= :date
)
ORDER BY map_id
        """,
        conn,
        params={"date": date},
    )


def load_team_fkfd(conn: sqlite3.Connection, team_id, date):
    return pd.read_sql(
        """
SELECT
    (
        SELECT
            (
                mm.team_1_def_score + mm.team_1_atk_score + mm.team_1_ot_score + mm.team_2_def_score + mm.team_2_atk_score + mm.team_2_ot_score
            )
        FROM
            match_maps AS mm
        WHERE
            mm.match_id = pos.match_id
            AND mm.map_id = pos.map_id
    ) AS rounds,
    SUM(pos.first_kills) AS fks,
    SUM(pos.first_deaths) AS fds
FROM
    player_overview_stats AS pos
    JOIN matches AS m ON m.id = pos.match_id
WHERE
    pos.team_id = :team_id
    AND m."date" > date(:date, '-180 days')
    AND m."date" < :date
GROUP BY
    pos.match_id,
    pos.map_id
ORDER BY
    m."date" DESC
        """,
        conn,
        params={"team_id": team_id, "date": date},
    )


def load_team_highlights(conn, team_id, date):
    return pd.read_sql(
        """
SELECT
    ph.*
FROM
    player_highlights AS ph
    JOIN matches AS m ON ph.match_id = m.id
WHERE
    ph.team_id = :team_id
    AND m."date" > date(:date, '-180 days')
    AND m."date" < :date
""",
        conn,
        params={"team_id": team_id, "date": date},
    )


def load_team_played_rounds_recently(conn, team_id, date):
    return pd.read_sql(
        """
SELECT
    rs.*
FROM
    round_stats AS rs
    JOIN matches AS m ON m.id = rs.match_id
WHERE
    (
        rs.team_1_id = :team_id
        OR rs.team_2_id = :team_id
    )
    AND m."date" > date(:date, '-180 days')
    AND m."date" < :date
        """,
        conn,
        params={"team_id": team_id, "date": date},
    )


def load_map_players_overview_stats(
    conn: sqlite3.Connection, match_id: int, map_id: int, team_id: int
):
    return pd.read_sql(
        """
SELECT
    *,
    (
        SELECT
            agent_type
        FROM
            agents
        WHERE
            id = agent_id
        LIMIT
            1
    ) AS role
FROM
    player_overview_stats
WHERE
    match_id = :match_id
    AND map_id = :map_id
    AND team_id = :team_id
        """,
        conn,
        params={"match_id": match_id, "map_id": map_id, "team_id": team_id},
    )


def load_players_highlights_log(
    conn: sqlite3.Connection, match_id: int, map_id: int, team_id: int
):
    return pd.read_sql(
        """
SELECT
    *
FROM
    player_highlights
WHERE
    match_id = :match_id
    AND map_id = :map_id
    AND team_id = :team_id;
        """,
        conn,
        params={"match_id": match_id, "map_id": map_id, "team_id": team_id},
    )


def load_map_rounds_stats(
    conn: sqlite3.Connection, match_id: int, map_id: int, team_id: int
):
    return pd.read_sql(
        """
SELECT
    match_id,
    map_id,
    round_no,
    (
        CASE
            WHEN team_1_id = :team_id THEN team_1_id
            WHEN team_2_id = :team_id THEN team_2_id
        END
    ) AS team_id,
    (
        CASE
            WHEN team_1_id = :team_id THEN team_2_id
            WHEN team_2_id = :team_id THEN team_1_id
        END
    ) AS team_against_id,
    team_def,
    (
        CASE
            WHEN team_1_id = :team_id THEN team_1_buy_type
            WHEN team_2_id = :team_id THEN team_2_buy_type
        END
    ) AS team_buy_type,
    (
        CASE
            WHEN team_1_id = :team_id THEN team_2_buy_type
            WHEN team_2_id = :team_id THEN team_1_buy_type
        END
    ) AS team_against_buy_type,
    (
        CASE
            WHEN team_1_id = :team_id THEN team_1_bank
            WHEN team_2_id = :team_id THEN team_2_bank
        END
    ) AS team_bank,
    (
        CASE
            WHEN team_1_id = :team_id THEN team_2_bank
            WHEN team_2_id = :team_id THEN team_1_bank
        END
    ) AS team_against_bank,
    team_won,
    won_method
FROM
    round_stats
WHERE
    match_id = :match_id
    AND map_id = :map_id
    AND (
        team_1_id = :team_id
        OR team_2_id = :team_id
    );
        """,
        conn,
        params={"match_id": match_id, "map_id": map_id, "team_id": team_id},
    )
