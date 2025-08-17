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
    AND "date" > date(:date, '-90 days')
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
  AND m."date" > date(:date, '-90 days')
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
    AND m."date" > date(:date, '-90 days')
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


def load_team_clutches_stats(conn, team_id, date):
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
            mm.match_id = ph.match_id
            AND mm.map_id = ph.map_id
    ) AS rounds,
    SUM(
        CASE
            WHEN ph.highlight_type = '1v1' THEN 1
            ELSE 0
        END
    ) / 1 AS p_1v1s,
    SUM(
        CASE
            WHEN ph.highlight_type = '1v2' THEN 1
            ELSE 0
        END
    ) / 2 AS p_1v2s,
    SUM(
        CASE
            WHEN ph.highlight_type = '1v3' THEN 1
            ELSE 0
        END
    ) / 3 AS p_1v3s,
    SUM(
        CASE
            WHEN ph.highlight_type = '1v4' THEN 1
            ELSE 0
        END
    ) / 4 AS p_1v4s,
    SUM(
        CASE
            WHEN ph.highlight_type = '1v5' THEN 1
            ELSE 0
        END
    ) / 5 AS p_1v5s
FROM
    player_highlights AS ph
    JOIN matches AS m ON ph.match_id = m.id
WHERE
    ph.team_id = :team_id
    AND ph.highlight_type IN ('1v1', '1v2', '1v3', '1v4', '1v5')
    AND m."date" > date(:date, '-90 days')
    AND m."date" < :date
GROUP BY
    ph.match_id,
    ph.map_id,
    ph.team_id
ORDER BY
    m.date DESC
        """,
        conn,
        params={"team_id": team_id, "date": date},
    )
